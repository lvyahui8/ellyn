package ellyn_ast

import (
	"github.com/lvyahui8/ellyn"
	"github.com/lvyahui8/ellyn/ellyn_agent"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"github.com/lvyahui8/ellyn/ellyn_testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProgramAll(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	prog.specifySdkDir = ellyn_testing.GetRepoRootPath()
	prog.sdkImportPkgPath = ellyn.SdkAgentPkg
	prog.RollbackAll()
	prog.Visit()
	utils.Go.Build(prog.mainPkg.Dir)
	prog.RollbackAll()
	prog.Destroy()
}

func TestCleanBackupFiles(t *testing.T) {
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	prog.cleanBackupFiles()
	prog.Destroy()
}

func TestFileEach(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	prog.scanSourceFiles(func(pkg *ellyn_agent.Package, fileAbsPath string) {
		t.Log(fileAbsPath)
	})
	prog.Destroy()
}

func TestExample(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	prog.RollbackAll()
	prog.specifySdkDir = ellyn_testing.GetRepoRootPath()
	prog.sdkImportPkgPath = ellyn.SdkAgentPkg
	prog.Visit()
	prog.Destroy()
}
