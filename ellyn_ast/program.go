package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"io/ioutil"
	"os"
	"path"
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

type Program struct {
	mainPkg   Package
	rootPkg   Package
	pkgMap    map[string]Package
	modFile   string
	allFuncs  []*GoFunc
	allBlocks []*Block
	processor Processor
	progCtx   *ProgramContext
	initOnce  sync.Once
}

func NewProgram(mainPkgDir string, processor Processor) *Program {
	if processor == nil {
		processor = DefaultProcessor{}
	}
	prog := &Program{
		mainPkg: Package{
			Dir: mainPkgDir,
		},
		pkgMap:    make(map[string]Package),
		processor: processor,
	}
	return prog
}

// Visit 触发项目扫描处理动作
func (p *Program) Visit() {
	p._init()
	p.handleFiles()
}

func (p *Program) _init() {
	p.initOnce.Do(func() {
		packages := utils.Go.AllPackages(p.mainPkg.Dir)
		for pkgPath, pkgDir := range packages {
			p.pkgMap[pkgDir] = NewPackage(pkgDir, pkgPath)
		}
		p.mainPkg.Name = p.pkgMap[p.mainPkg.Dir].Name
		p.modFile = utils.Go.GetModFile(p.mainPkg.Dir)
		rootPkgPath := utils.Go.GetRootPkgPath(p.modFile)
		p.rootPkg = NewPackage(path.Dir(p.modFile), rootPkgPath)
		p.progCtx = &ProgramContext{rootPkgPath: rootPkgPath}
	})
}

func (p *Program) handleFiles() {
	for pkgDir, pkg := range p.pkgMap {
		if !p.processor.FilterPackage(p.progCtx, pkg) {
			continue
		}
		files, err := ioutil.ReadDir(pkgDir)
		asserts.IsNil(err)
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".go") {
				continue
			}
			if !p.processor.FilterFile(p.progCtx, pkg, file) {
				continue
			}
			// 加入遍历集合
			fileAbsPath := pkg.Dir + string(os.PathSeparator) + file.Name()
			content, err := ioutil.ReadFile(fileAbsPath)
			asserts.IsNil(err)
			p.processor.HandleFile(p.progCtx, pkg, file, content)
		}
	}
}
