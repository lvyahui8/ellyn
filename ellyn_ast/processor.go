package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_common/log"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"os"
	"strings"
)

type Processor interface {
	FilterPackage(progCtx *ProgramContext, pkg Package) bool
	FilterFile(progCtx *ProgramContext, pkg Package, file os.DirEntry) bool
	HandleFile(progCtx *ProgramContext, pkg Package, file os.DirEntry, content []byte)
	HandlePackage(progCtx *ProgramContext, pkg Package)
	HandleFunc(progCtx *ProgramContext, pkg Package, file os.DirEntry, goFunc *GoFunc)
}

type DefaultProcessor struct {
}

func (d DefaultProcessor) FilterPackage(progCtx *ProgramContext, pkg Package) bool {
	return strings.HasPrefix(pkg.Path, progCtx.RootPkgPath())
}

func (d DefaultProcessor) FilterFile(progCtx *ProgramContext, pkg Package, file os.DirEntry) bool {
	if file.IsDir() {
		return false
	}

	if utils.Go.IsTestFile(file.Name()) {
		return false
	}

	if utils.Go.IsAutoGenFile(pkg.Dir + string(os.PathSeparator) + file.Name()) {
		return false
	}

	return true
}

func (d DefaultProcessor) HandlePackage(progCtx *ProgramContext, pkg Package) {

}

func (d DefaultProcessor) HandleFile(progCtx *ProgramContext, pkg Package, file os.DirEntry, content []byte) {
	log.Infof("dir %s,file %s", pkg.Dir, file.Name())
}

func (d DefaultProcessor) HandleFunc(progCtx *ProgramContext, pkg Package, file os.DirEntry, goFunc *GoFunc) {

}
