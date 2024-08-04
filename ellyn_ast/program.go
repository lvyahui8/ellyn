package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"github.com/lvyahui8/ellyn/ellyn_common/log"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
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
	mainPkg      Package
	rootPkg      Package
	pkgMap       map[string]Package
	modFile      string
	allFuncs     []*GoFunc
	allBlocks    []*Block
	progCtx      *ProgramContext
	initOnce     sync.Once
	funcCounter  atomic.Int32
	blockCounter atomic.Int32
}

func NewProgram(mainPkgDir string) *Program {
	prog := &Program{
		mainPkg: Package{
			Dir: mainPkgDir,
		},
		pkgMap: make(map[string]Package),
	}
	return prog
}

// Visit 触发项目扫描处理动作
func (p *Program) Visit() {
	p._init()
	p.handleFiles()
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

func (p *Program) handleFiles() {
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

			// 加入遍历集合
			fileAbsPath := pkg.Dir + string(os.PathSeparator) + file.Name()
			content, err := os.ReadFile(fileAbsPath)
			asserts.IsNil(err)
			log.Infof("dir %s,file %s", pkg.Dir, file.Name())
			visitor := &FileVisitor{content: content}
			fset := token.NewFileSet()
			visitor.fset = fset
			parsedFile, err := parser.ParseFile(fset, fileAbsPath, content, parser.ParseComments)
			ast.Walk(visitor, parsedFile)
		}
	}
}
