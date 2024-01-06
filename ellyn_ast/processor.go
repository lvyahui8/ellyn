package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_common/log"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"io/fs"
	"os"
	"strings"
)

type Processor interface {
	FilterPackage(progCtx *ProgramContext, pkg Package) bool
	FilterFile(progCtx *ProgramContext, pkg Package, file fs.FileInfo) bool
	HandleFile(progCtx *ProgramContext, pkg Package, file fs.FileInfo, content []byte)
}

type DefaultProcessor struct {
}

func (d DefaultProcessor) FilterPackage(progCtx *ProgramContext, pkg Package) bool {
	return strings.HasPrefix(pkg.Path, progCtx.RootPkgPath())
}

func (d DefaultProcessor) FilterFile(progCtx *ProgramContext, pkg Package, file fs.FileInfo) bool {
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

func (d DefaultProcessor) HandleFile(progCtx *ProgramContext, pkg Package, file fs.FileInfo, content []byte) {
	log.Infof("dir %s,file %s", pkg.Dir, file.Name())
}
