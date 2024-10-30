package instr

import (
	"github.com/lvyahui8/ellyn/sdk"
	"github.com/lvyahui8/ellyn/sdk/agent"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"github.com/lvyahui8/ellyn/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProgramAll(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(test.GetTestProjPath())
	defer prog.Destroy()
	prog.specifySdkDir = test.GetRepoRootPath()
	prog.sdkImportPkgPath = sdk.SdkAgentPkg
	prog.RollbackAll()
	prog.Visit()
	utils.Go.Build(prog.mainPkg.Dir)
	prog.RollbackAll()
}

func TestCleanBackupFiles(t *testing.T) {
	prog := NewProgram(test.GetTestProjPath())
	defer prog.Destroy()
	prog.cleanBackupFiles()
}

func TestFileEach(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(test.GetTestProjPath())
	defer prog.Destroy()
	prog.scanSourceFiles(func(pkg *agent.Package, fileAbsPath string) {
		t.Log(fileAbsPath)
	})
}

func TestExample(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(test.GetTestProjPath())
	defer prog.Destroy()
	prog.RollbackAll()
	prog.specifySdkDir = test.GetRepoRootPath()
	prog.sdkImportPkgPath = sdk.SdkAgentPkg
	prog.Visit()

}

func TestRollbackExample(t *testing.T) {
	prog := NewProgram(test.GetTestProjPath())
	defer prog.Destroy()
	prog.RollbackAll()
}

func TestUpdateBenchmark(t *testing.T) {
	prog := NewProgram2(test.GetBenchmarkPath(), agent.Configuration{
		NoArgs: true,
		NoDemo: true,
	})
	defer prog.Destroy()
	prog.RollbackAll()
	prog.Visit()
}

func TestRollbackBenchmark(t *testing.T) {
	prog := NewProgram(test.GetBenchmarkPath())
	defer prog.Destroy()
	prog.RollbackAll()
}
