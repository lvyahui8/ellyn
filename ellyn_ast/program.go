package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"io/ioutil"
	"path"
	"strings"
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
		pkgMap: make(map[string]Package),
	}
	return prog
}

// Visit 触发项目扫描处理动作
func (p *Program) Visit() {
	p._init()

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
	})
}

func (p *Program) listFiles() {
	for pkgDir, pkg := range p.pkgMap {
		if !strings.HasPrefix(pkg.Path, p.rootPkg.Path) {
			continue
		}
		files, err := ioutil.ReadDir(pkgDir)
		asserts.IsNil(err)
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if utils.Go.IsTestFile(file.Name()) {
				continue
			}

			if utils.Go.IsAutoGenFile(file.Name()) {
				continue
			}
			// 加入遍历集合
		}
	}
}
