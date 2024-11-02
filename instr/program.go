package instr

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/lvyahui8/ellyn"
	"github.com/lvyahui8/ellyn/sdk/agent"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/goroutine"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/mod/modfile"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type fileHandler func(pkg *agent.Package, fileAbsPath string)

// Program 封装程序信息，解析遍历程序的所有包、函数、代码块
type Program struct {
	conf    agent.Configuration
	mainPkg *agent.Package
	rootPkg *agent.Package

	path2pkgMap    map[string]*agent.Package
	dir2pkgMap     map[string]*agent.Package
	allFiles       *collections.ConcurrentMap[uint32, *agent.File]
	allMethods     *collections.ConcurrentMap[uint32, *agent.Method]
	fileMethodsMap *collections.ConcurrentMap[uint32, *treeset.Set]
	allBlocks      *collections.ConcurrentMap[uint32, *agent.Block]

	initOnce       sync.Once
	packageCounter uint32
	fileCounter    int32
	methodCounter  int32
	blockCounter   int32

	executor *goroutine.RoutinePool

	sdkImportPkgPath string
	modFilePath      string
	modFile          *modfile.File
	targetPath       string
	updatedFiles     []string
	specifySdkDir    string
}

func NewProgram(mainPkgDir string) *Program {
	return NewProgram2(mainPkgDir, agent.Configuration{})
}
func NewProgram2(mainPkgDir string, configuration agent.Configuration) *Program {
	mainPkgDir = filepath.ToSlash(mainPkgDir)
	prog := &Program{
		conf: configuration,
		mainPkg: &agent.Package{
			Dir: mainPkgDir,
		},
		path2pkgMap:    make(map[string]*agent.Package),
		dir2pkgMap:     make(map[string]*agent.Package),
		executor:       goroutine.NewRoutinePool(runtime.NumCPU()<<1, false),
		allFiles:       collections.NewNumberKeyConcurrentMap[uint32, *agent.File](4),
		fileMethodsMap: collections.NewNumberKeyConcurrentMap[uint32, *treeset.Set](8),
		allMethods:     collections.NewNumberKeyConcurrentMap[uint32, *agent.Method](32),
		allBlocks:      collections.NewNumberKeyConcurrentMap[uint32, *agent.Block](32),
		targetPath:     mainPkgDir,
		fileCounter:    -1,
		methodCounter:  -1,
		blockCounter:   -1,
	}
	prog._init()
	return prog
}

// _init 初始化基础信息，为文件迭代做准备
func (p *Program) _init() {
	p.initOnce.Do(func() {
		packages := utils.Go.AllPackages(p.mainPkg.Dir)
		for pkgPath, pkgDir := range packages {
			if strings.Contains(pkgPath, ellyn.AgentPkg) {
				continue
			}
			pkg := agent.NewPackage(pkgDir, pkgPath)
			pkg.Id = p.packageCounter
			p.packageCounter++
			p.dir2pkgMap[pkgDir] = pkg
			p.path2pkgMap[pkgPath] = pkg
		}
		p.mainPkg.Name = p.dir2pkgMap[p.mainPkg.Dir].Name
		p.modFilePath = utils.Go.GetModFile(p.mainPkg.Dir)
		p.modFile = p.parseModFile(p.modFilePath)
		rootPkgPath := p.modFile.Module.Mod.Path
		p.rootPkg = agent.NewPackage(path.Dir(p.modFilePath), rootPkgPath)
		p.sdkImportPkgPath = fmt.Sprintf("%s/ellyn_agent", p.rootPkg.Path)
	})
}

// Visit 触发项目扫描处理动作
func (p *Program) Visit() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("rollbackAll,panic err %v,stack:%s", err, string(debug.Stack()))
			p.RollbackAll()
		}
	}()
	p.scanSourceFiles(p.updateFile)
	p.buildApp()
	p.buildAgent()
	p.buildMeta()
}

func (p *Program) addFile(pkgId uint32, file string) *agent.File {
	f := &agent.File{
		FileId:       uint32(atomic.AddInt32(&p.fileCounter, 1)),
		PackageId:    pkgId,
		RelativePath: file,
	}
	p.allFiles.Store(f.FileId, f)
	return f
}

func (p *Program) Destroy() {
	p.executor.Shutdown()
}

func (p *Program) addMethod(fileId uint32, methodName string, begin, end token.Position, funcType *ast.FuncType) *agent.Method {
	file, ok := p.allFiles.Load(fileId)
	asserts.True(ok)

	f := &agent.Method{
		Id:         uint32(atomic.AddInt32(&p.methodCounter, 1)),
		FileId:     fileId,
		Name:       methodName,
		FullName:   methodName,
		PackageId:  file.PackageId,
		Begin:      agent.NewPos(begin.Offset, begin.Line, begin.Column),
		End:        agent.NewPos(end.Offset, end.Line, end.Column),
		ArgsList:   p.filedList2VarDefList(funcType.Params),
		ReturnList: p.filedList2VarDefList(funcType.Results),
	}

	p.allMethods.Store(f.Id, f)
	fileAllFuncs, ok := p.fileMethodsMap.Load(fileId)
	if !ok {
		fileAllFuncs = treeset.NewWith(func(a, b interface{}) int {
			return a.(*agent.Method).Begin.Offset - b.(*agent.Method).Begin.Offset
		})
		p.fileMethodsMap.Store(fileId, fileAllFuncs)
	}

	fileAllFuncs.Add(f)
	return f
}

func (p *Program) filedList2VarDefList(fieldList *ast.FieldList) *agent.VarDefList {
	if fieldList == nil || fieldList.List == nil {
		return agent.NewVarDefList(nil)
	}
	var list []*agent.VarDef
	for _, field := range fieldList.List {
		var names []string
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
		list = append(list, &agent.VarDef{
			Names: names,
			Type:  types.ExprString(field.Type), // 获取type的string表示
		})
	}
	return agent.NewVarDefList(list)
}

func (p *Program) findMethod(fileId uint32, offset int) *agent.Method {
	set, ok := p.fileMethodsMap.Load(fileId)
	if !ok {
		return nil
	}
	values := set.Values()
	var target *agent.Method
	for _, v := range values {
		f := v.(*agent.Method)
		if f.Begin.Offset > offset {
			break
		}
		if f.Begin.Offset <= offset && f.End.Offset >= offset {
			target = f
		}
	}
	return target
}

// buildMethods 完成方法内容的善后工作
func (p *Program) buildMethods(fileId uint32) {
	// 计算Block Offset
	fileMethods, ok := p.fileMethodsMap.Load(fileId)
	if !ok {
		// 文件没有函数
		// fmt.Printf("file %d no methods\n", fileId)
		return
	}
	fileMethods.Each(func(index int, value interface{}) {
		m := value.(*agent.Method)
		sort.Slice(m.Blocks, func(i, j int) bool {
			return m.Blocks[i].Begin.Offset-m.Blocks[j].Begin.Offset < 0
		})
		for offset, b := range m.Blocks {
			b.MethodOffset = offset
		}
	})
	// 计算匿名函数名
}

func (p *Program) addBlock(fileId uint32, begin, end token.Position) *agent.Block {
	method := p.findMethod(fileId, begin.Offset)
	b := &agent.Block{
		Id:       uint32(atomic.AddInt32(&p.blockCounter, 1)),
		MethodId: method.Id,
		FileId:   fileId,
		Begin:    agent.NewPos(begin.Offset, begin.Line, begin.Column),
		End:      agent.NewPos(end.Offset, end.Line, end.Column),
	}
	method.Blocks = append(method.Blocks, b)
	p.allBlocks.Store(b.Id, b)
	return b
}

func (p *Program) scanSourceFiles(handler fileHandler) {
	fileGroup := &sync.WaitGroup{}
	for pkgDir, pkg := range p.dir2pkgMap {
		if !strings.HasPrefix(pkg.Path, p.rootPkg.Path) || strings.HasSuffix(pkg.Path, ellyn.AgentPkg) {
			continue
		}
		files, err := os.ReadDir(pkgDir)
		asserts.IsNil(err)
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".go") ||
				file.IsDir() ||
				utils.Go.IsTestFile(file.Name()) ||
				utils.Go.IsAutoGenFile(pkg.Dir+string(os.PathSeparator)+file.Name()) {
				continue
			}
			// 将文件加入遍历队列并发处理，加快文件处理速度
			p.handleFile(pkg, file, handler, fileGroup)
		}
	}

	// 等待所有文件处理完成
	fileGroup.Wait()
}

func (p *Program) buildAgent() {
	p.copySdk(ellyn.SdkDir)
}

func (p *Program) copySdk(sdkPath string) {
	files, err := ellyn.SdkFs.ReadDir(sdkPath)
	asserts.IsNil(err)
	for _, file := range files {
		if !file.IsDir() && !utils.Go.IsSourceFile(file.Name()) &&
			!strings.Contains(filepath.ToSlash(sdkPath), "/page") {
			continue
		}
		fmt.Printf("sdk relativePath :%s\n", file.Name())
		rPath := path.Join(sdkPath, file.Name())
		if file.IsDir() {
			p.copySdk(rPath)
		} else {
			if strings.HasSuffix(filepath.ToSlash(rPath), agent.AgentApiFile) {
				if !p.require(ellyn.ApiPackage) {
					continue
				}
			}
			bytes, err := ellyn.SdkFs.ReadFile(rPath)
			asserts.IsNil(err)

			updated := strings.ReplaceAll(
				utils.String.Bytes2string(bytes), ellyn.SdkRawRootPkg, p.rootPkg.Path+"/"+ellyn.AgentPkg)
			bytes = utils.String.String2bytes(updated)

			utils.OS.WriteTo(path.Join(p.targetPath,
				strings.Replace(rPath, ellyn.SdkDir, ellyn.AgentPkg, 1)), bytes)
		}
	}
}

func (p *Program) handleFile(pkg *agent.Package, file os.DirEntry, handler fileHandler, fileGroup *sync.WaitGroup) {
	fileGroup.Add(1)
	// 这里使用阻塞队列，队列不限制容量，确保文件不会被丢弃
	p.executor.Submit(func() {
		defer fileGroup.Done()
		fileAbsPath := filepath.Join(pkg.Dir, file.Name())
		fmt.Printf("dir %s,relativePath %s\n", pkg.Dir, file.Name())
		handler(pkg, fileAbsPath)
	})
}

func (p *Program) updateFile(pkg *agent.Package, fileAbsPath string) {
	p.backup(fileAbsPath)
	relativePath := strings.ReplaceAll(filepath.ToSlash(fileAbsPath), filepath.ToSlash(p.mainPkg.Dir), "")
	f := p.addFile(pkg.Id, relativePath)
	content, err := os.ReadFile(fileAbsPath)
	asserts.IsNil(err)
	p.copySource(relativePath, content)
	visitor := &FileVisitor{
		fileId:       f.FileId,
		prog:         p,
		content:      content,
		relativePath: relativePath,
	}
	fset := token.NewFileSet()
	visitor.fset = fset
	parsedFile, err := parser.ParseFile(fset, fileAbsPath, content, parser.ParseComments)
	asserts.IsNil(err)
	f.LineNum = fset.Position(parsedFile.End()).Line
	ast.Walk(visitor, parsedFile)
	p.buildMethods(f.FileId)
	visitor.WriteTo(fileAbsPath)
	p.updatedFiles = append(p.updatedFiles, fileAbsPath)
}

func (p *Program) copySource(relativePath string, content []byte) {
	sourcesPath := filepath.Join(p.targetPath, ellyn.AgentPkg, agent.SourcesRelativePath)
	utils.OS.WriteTo(filepath.Join(sourcesPath, relativePath)+agent.SourcesFileExt, content)
}

func (p *Program) buildApp() {
	utils.OS.CopyFile(p.modFilePath, filepath.Join(p.targetPath, "go.mod"))
}

// parseModFile 获取项目go.mod文件所在的package name
func (p *Program) parseModFile(modFilePath string) *modfile.File {
	content, err := os.ReadFile(modFilePath)
	asserts.IsNil(err)
	modFile, err := modfile.Parse("go.mod", content, nil)
	asserts.IsNil(err)
	return modFile
}

func (p *Program) require(pkgPath string) bool {
	for _, r := range p.modFile.Require {
		if strings.HasPrefix(pkgPath, r.Mod.Path) {
			return true
		}
	}
	return false
}

func (p *Program) backup(fileAbsPath string) {
	bakFile := fileAbsPath + ".bak"
	fmt.Printf("backup file:%s, target file:%s\n", fileAbsPath, bakFile)
	utils.OS.CopyFile(fileAbsPath, bakFile)
}

func (p *Program) rollback(fileAbsPath string) {
	bakFile := fileAbsPath + ".bak"
	if utils.OS.NotExists(bakFile) {
		return
	}
	fmt.Printf("rollback file:%s, target file:%s\n", bakFile, fileAbsPath)
	utils.OS.CopyFile(bakFile, fileAbsPath)
	utils.OS.Remove(bakFile)
}

func (p *Program) RollbackAll() {
	if p.updatedFiles != nil {
		for _, f := range p.updatedFiles {
			p.rollback(f)
		}
	} else {
		p.scanSourceFiles(func(pkg *agent.Package, fileAbsPath string) {
			p.rollback(fileAbsPath)
		})
	}
	utils.OS.Remove(filepath.ToSlash(filepath.Join(p.targetPath, ellyn.AgentPkg)))
}

func (p *Program) cleanBackupFiles() {
	p.scanSourceFiles(func(pkg *agent.Package, fileAbsPath string) {
		utils.OS.Remove(fileAbsPath + ".bak")
	})
}

// buildMeta 构建元数据，将元数据写入项目
func (p *Program) buildMeta() {
	target := p.targetPath
	metaPath := filepath.Join(target, ellyn.AgentPkg, agent.MetaRelativePath)
	utils.OS.MkDirs(metaPath)
	// 写入运行时配置
	confBytes, err := json.Marshal(p.conf)
	asserts.IsNil(err)
	utils.OS.WriteTo(filepath.Join(metaPath, agent.RuntimeConfFile), confBytes)
	// 写入包、文件、函数、块数据
	pkgList := utils.GetMapValues(p.dir2pkgMap)
	sort.Slice(pkgList, func(i, j int) bool {
		return pkgList[i].Id < pkgList[j].Id
	})
	utils.OS.WriteTo(filepath.Join(metaPath, agent.MetaPackages), agent.EncodeCsvRows(pkgList))
	utils.OS.WriteTo(filepath.Join(metaPath, agent.MetaFiles),
		agent.EncodeCsvRows(p.allFiles.SortedValues(func(a, b *agent.File) bool {
			return a.FileId < b.FileId
		})))
	utils.OS.WriteTo(filepath.Join(metaPath, agent.MetaMethods),
		agent.EncodeCsvRows(p.allMethods.SortedValues(func(a, b *agent.Method) bool {
			return a.Id < b.Id
		})))
	utils.OS.WriteTo(filepath.Join(metaPath, agent.MetaBlocks),
		agent.EncodeCsvRows(p.allBlocks.SortedValues(func(a, b *agent.Block) bool {
			return a.Id < b.Id
		})))
}
