package ellyn_ast

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/lvyahui8/ellyn"
	"github.com/lvyahui8/ellyn/ellyn_agent"
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/lvyahui8/ellyn/ellyn_common/goroutine"
	"github.com/lvyahui8/ellyn/ellyn_common/log"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type ProgramContext struct {
	rootPkgPath string
	extra       sync.Map
}

func (p *ProgramContext) RootPkgPath() string {
	return p.rootPkgPath
}

// Program 封装程序信息，解析遍历程序的所有包、函数、代码块
type Program struct {
	mainPkg *ellyn_agent.Package
	rootPkg *ellyn_agent.Package
	pkgMap  map[string]*ellyn_agent.Package
	modFile string

	fileMethodsMap *collections.ConcurrentMap[uint32, *treeset.Set]

	allFuncs  *collections.ConcurrentMap[uint32, *ellyn_agent.Method]
	allBlocks *collections.ConcurrentMap[uint32, *ellyn_agent.Block]

	progCtx      *ProgramContext
	initOnce     sync.Once
	fileCounter  uint32
	funcCounter  uint32
	blockCounter uint32
	executor     *goroutine.RoutinePool
	w            *sync.WaitGroup
	targetPath   string
}

func NewProgram(mainPkgDir string) *Program {
	prog := &Program{
		mainPkg: &ellyn_agent.Package{
			Dir: mainPkgDir,
		},
		pkgMap:         make(map[string]*ellyn_agent.Package),
		executor:       goroutine.NewRoutinePool(runtime.NumCPU()<<1, false),
		w:              &sync.WaitGroup{},
		fileMethodsMap: collections.NewNumberKeyConcurrentMap[uint32, *treeset.Set](8),
		allFuncs:       collections.NewNumberKeyConcurrentMap[uint32, *ellyn_agent.Method](32),
		allBlocks:      collections.NewNumberKeyConcurrentMap[uint32, *ellyn_agent.Block](32),
		targetPath:     mainPkgDir + "/target/",
	}
	return prog
}

// Visit 触发项目扫描处理动作
func (p *Program) Visit() {
	p._init()
	p.scanSourceFiles()
	p.buildApp()
	p.buildAgent()
}

// _init 初始化基础信息，为文件迭代做准备
func (p *Program) _init() {
	p.initOnce.Do(func() {
		packages := utils.Go.AllPackages(p.mainPkg.Dir)
		for pkgPath, pkgDir := range packages {
			p.pkgMap[pkgDir] = ellyn_agent.NewPackage(pkgDir, pkgPath)
		}
		p.mainPkg.Name = p.pkgMap[p.mainPkg.Dir].Name

		p.modFile = utils.Go.GetModFile(p.mainPkg.Dir)
		rootPkgPath := p.getProjectRootPkgPath(p.modFile)
		p.rootPkg = ellyn_agent.NewPackage(path.Dir(p.modFile), rootPkgPath)
		p.progCtx = &ProgramContext{rootPkgPath: rootPkgPath}
	})
}

func (p *Program) addMethod(fileId uint32, fcName string, begin, end token.Position) *ellyn_agent.Method {
	f := &ellyn_agent.Method{
		Id:       p.funcCounter,
		FileId:   fileId,
		FullName: fcName,
		Begin:    ellyn_agent.NewPos(begin.Offset, begin.Line, begin.Column),
		End:      ellyn_agent.NewPos(end.Offset, end.Line, end.Column),
	}
	p.allFuncs.Store(f.Id, f)
	fileAllFuncs, ok := p.fileMethodsMap.Load(fileId)
	if !ok {
		fileAllFuncs = treeset.NewWith(func(a, b interface{}) int {
			return a.(*ellyn_agent.Method).Begin.Offset - b.(*ellyn_agent.Method).Begin.Offset
		})
		p.fileMethodsMap.Store(fileId, fileAllFuncs)
	}
	atomic.AddUint32(&p.funcCounter, 1)

	fileAllFuncs.Add(f)
	return f
}

func (p *Program) findMethod(fileId uint32, offset int) *ellyn_agent.Method {
	set, ok := p.fileMethodsMap.Load(fileId)
	if !ok {
		return nil
	}
	values := set.Values()
	var target *ellyn_agent.Method
	for _, v := range values {
		f := v.(*ellyn_agent.Method)
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
	asserts.True(ok)
	fileMethods.Each(func(index int, value interface{}) {
		m := value.(*ellyn_agent.Method)
		sort.Slice(m.Blocks, func(i, j int) bool {
			return m.Blocks[i].Begin.Offset-m.Blocks[j].Begin.Offset < 0
		})
		for offset, b := range m.Blocks {
			b.MethodOffset = offset
		}
	})
	// 计算匿名函数名
}

func (p *Program) addBlock(fileId uint32, begin, end token.Position) *ellyn_agent.Block {
	method := p.findMethod(fileId, begin.Offset)
	b := &ellyn_agent.Block{
		Id:       p.blockCounter,
		MethodId: method.Id,
		Begin:    ellyn_agent.NewPos(begin.Offset, begin.Line, begin.Column),
		End:      ellyn_agent.NewPos(end.Offset, end.Line, end.Column),
	}
	method.Blocks = append(method.Blocks, b)
	atomic.AddUint32(&p.blockCounter, 1)
	p.allBlocks.Store(b.Id, b)
	return b
}

func (p *Program) scanSourceFiles() {
	for pkgDir, pkg := range p.pkgMap {
		if !strings.HasPrefix(pkg.Path, p.progCtx.RootPkgPath()) {
			continue
		}
		files, err := os.ReadDir(pkgDir)
		asserts.IsNil(err)
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".go") {
				continue
			}
			if file.IsDir() {
				continue
			}

			if utils.Go.IsTestFile(file.Name()) {
				continue
			}
			if utils.Go.IsAutoGenFile(pkg.Dir + string(os.PathSeparator) + file.Name()) {
				continue
			}

			// 将文件加入遍历队列并发处理，加快文件处理速度
			p.parseFile(pkg, file)
		}
	}

	// 等待所有文件处理完成
	p.w.Wait()
	p.executor.Shutdown()
}

func (p *Program) buildAgent() {
	for _, sdkPath := range ellyn.SdkPaths {
		p.copySdk(sdkPath)
	}
}

func (p *Program) copySdk(sdkPath string) {
	files, err := ellyn.SdkFs.ReadDir(sdkPath)
	asserts.IsNil(err)
	for _, file := range files {
		if !file.IsDir() && !utils.Go.IsSourceFile(file.Name()) {
			continue
		}
		fmt.Printf("sdk file :%s\n", file.Name())
		rPath := path.Join(sdkPath, file.Name())
		if file.IsDir() {
			p.copySdk(rPath)
		} else {
			bytes, err := ellyn.SdkFs.ReadFile(rPath)
			asserts.IsNil(err)
			updated := strings.ReplaceAll(
				utils.String.Bytes2string(bytes), ellyn.SdkRawRootPkg, p.rootPkg.Path)
			utils.OS.WriteTo(path.Join(p.targetPath, rPath), utils.String.String2bytes(updated))
		}
	}
}

func (p *Program) parseFile(pkg *ellyn_agent.Package, file os.DirEntry) {
	p.w.Add(1)
	// 这里使用阻塞队列，队列不限制容量，确保文件不会被丢弃
	p.executor.Submit(func() {
		defer p.w.Done()
		fileAbsPath := pkg.Dir + string(os.PathSeparator) + file.Name()
		content, err := os.ReadFile(fileAbsPath)
		asserts.IsNil(err)
		log.Infof("dir %s,file %s", pkg.Dir, file.Name())
		fileId := atomic.AddUint32(&p.fileCounter, 1)
		visitor := &FileVisitor{
			fileId:  fileId,
			prog:    p,
			content: content,
			file:    strings.ReplaceAll(utils.File.FormatFilePath(fileAbsPath), p.mainPkg.Dir, ""),
		}
		fset := token.NewFileSet()
		visitor.fset = fset
		parsedFile, err := parser.ParseFile(fset, fileAbsPath, content, parser.ParseComments)
		ast.Walk(visitor, parsedFile)
		p.buildMethods(fileId)
		visitor.WriteTo(p.targetPath)
	})
}

func (p *Program) buildApp() {
	utils.OS.CopyFile(p.modFile, filepath.Join(p.targetPath, "go.mod"))
}

// getProjectRootPkgPath 获取项目go.mod文件所在的package name
func (p *Program) getProjectRootPkgPath(modFilePath string) string {
	content, err := os.ReadFile(modFilePath)
	asserts.IsNil(err)
	modFile, err := modfile.Parse("go.mod", content, nil)
	asserts.IsNil(err)
	return modFile.Module.Mod.Path
}
