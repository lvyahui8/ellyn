package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"sync"
)

type Program struct {
	mainPkg   Package
	rootPkg   Package
	pkgMap    map[string]Package
	modFile   string
	allFuncs  []*GoFunc
	allBlocks []*Block

	initOnce sync.Once
}

func NewProgram(mainPkgDir string) *Program {
	prog := &Program{
		mainPkg: Package{
			Dir: mainPkgDir,
		},
	}
	return prog
}

// Visit 触发项目扫描处理动作
func (p *Program) Visit() {
	p.init()
}

func (p *Program) init() {
	p.initOnce.Do(func() {
		packages := utils.Go.AllPackages(p.mainPkg.Dir)
		for dir, pkg := range packages {
			p.pkgMap[dir] = Package{Dir: dir, Name: pkg}
		}
		p.mainPkg.Name = p.pkgMap[p.mainPkg.Dir].Name
		p.modFile = utils.Go.GetModFile(p.mainPkg.Dir)
		rootPkgName := utils.Go.GetRootPkg(p.modFile)
		p.rootPkg = Package{Dir: p.modFile, Name: rootPkgName}
	})
}
