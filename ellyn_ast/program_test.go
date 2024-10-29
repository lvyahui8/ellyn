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
	defer prog.Destroy()
	prog.specifySdkDir = ellyn_testing.GetRepoRootPath()
	prog.sdkImportPkgPath = ellyn.SdkAgentPkg
	prog.RollbackAll()
	prog.Visit()
	utils.Go.Build(prog.mainPkg.Dir)
	prog.RollbackAll()
}

func TestCleanBackupFiles(t *testing.T) {
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	defer prog.Destroy()
	prog.cleanBackupFiles()
}

func TestFileEach(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	defer prog.Destroy()
	prog.scanSourceFiles(func(pkg *ellyn_agent.Package, fileAbsPath string) {
		t.Log(fileAbsPath)
	})
}

func TestExample(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	defer prog.Destroy()
	prog.RollbackAll()
	prog.specifySdkDir = ellyn_testing.GetRepoRootPath()
	prog.sdkImportPkgPath = ellyn.SdkAgentPkg
	prog.Visit()

}

func TestRollbackExample(t *testing.T) {
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	defer prog.Destroy()
	prog.RollbackAll()
}

func TestUpdateBenchmark(t *testing.T) {
	prog := NewProgram2(ellyn_testing.GetBenchmarkPath(), ellyn_agent.Configuration{
		NoArgs: true,
		NoDemo: true,
	})
	defer prog.Destroy()
	prog.RollbackAll()
	prog.Visit()
}

func TestRollbackBenchmark(t *testing.T) {
	prog := NewProgram(ellyn_testing.GetBenchmarkPath())
	defer prog.Destroy()
	prog.RollbackAll()
}
