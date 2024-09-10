package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"github.com/lvyahui8/ellyn/ellyn_common/goroutine"
	"github.com/lvyahui8/ellyn/ellyn_common/log"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
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
	mainPkg      Package
	rootPkg      Package
	pkgMap       map[string]Package
	modFile      string
	allFuncs     []*GoFunc
	allBlocks    []*Block
	progCtx      *ProgramContext
	initOnce     sync.Once
	funcCounter  int64
	blockCounter int64
	executor     *goroutine.RoutinePool
	w            *sync.WaitGroup
}

func NewProgram(mainPkgDir string) *Program {
	prog := &Program{
		mainPkg: Package{
			Dir: mainPkgDir,
		},
		pkgMap:   make(map[string]Package),
		executor: goroutine.NewRoutinePool(runtime.NumCPU() << 1),
		w:        &sync.WaitGroup{},
	}
	return prog
}

// Visit 触发项目扫描处理动作
func (p *Program) Visit() {
	p._init()
	p.scanFiles()
}

// _init 初始化基础信息，为文件迭代做准备
func (p *Program) _init() {
	p.initOnce.Do(func() {
		packages := utils.Go.AllPackages(p.mainPkg.Dir)
		for pkgPath, pkgDir := range packages {
			p.pkgMap[pkgDir] = NewPackage(pkgDir, pkgPath)
		}
		p.mainPkg.Name = p.pkgMap[p.mainPkg.Dir].Name

		p.modFile = utils.Go.GetModFile(p.mainPkg.Dir)
		rootPkgPath := utils.Go.GetProjectRootPkgPath(p.modFile)
		p.rootPkg = NewPackage(path.Dir(p.modFile), rootPkgPath)
		p.progCtx = &ProgramContext{rootPkgPath: rootPkgPath}
	})
}

func (p *Program) scanFiles() {
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

func (p *Program) parseFile(pkg Package, file os.DirEntry) {
	p.w.Add(1)
	// 这里使用阻塞队列，队列不限制容量，确保文件不会被丢弃
	p.executor.Submit(func() {
		defer p.w.Done()
		fileAbsPath := pkg.Dir + string(os.PathSeparator) + file.Name()
		content, err := os.ReadFile(fileAbsPath)
		asserts.IsNil(err)
		log.Infof("dir %s,file %s", pkg.Dir, file.Name())
		visitor := &FileVisitor{content: content}
		fset := token.NewFileSet()
		visitor.fset = fset
		parsedFile, err := parser.ParseFile(fset, fileAbsPath, content, parser.ParseComments)
		ast.Walk(visitor, parsedFile)
	})
}
