package instr

import (
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
	prog := NewProgram(test.GetTestProjPath(), true, nil)
	defer prog.Destroy()
	prog.useRawSdk = true
	prog.RollbackAll()
	prog.Visit()
	utils.Go.Build(prog.mainPkg.Dir)
	prog.RollbackAll()
}

func TestCleanBackupFiles(t *testing.T) {
	prog := NewProgram(test.GetTestProjPath(), true, nil)
	defer prog.Destroy()
	prog.cleanBackupFiles()
}

func TestFileEach(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(test.GetTestProjPath(), true, nil)
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
	prog := NewProgram(test.GetTestProjPath(), true, nil)
	defer prog.Destroy()
	prog.RollbackAll()
	prog.useRawSdk = true
	prog.Visit()

}

func TestRollbackExample(t *testing.T) {
	prog := NewProgram(test.GetTestProjPath(), true, nil)
	defer prog.Destroy()
	prog.RollbackAll()
}

func TestUpdateBenchmark(t *testing.T) {
	prog := NewProgram(test.GetBenchmarkPath(), true, &agent.Configuration{
		NoArgs:       true,
		NoDemo:       true,
		SamplingRate: 0.001,
	})
	defer prog.Destroy()
	prog.RollbackAll()
	prog.Visit()
}

func TestRollbackBenchmark(t *testing.T) {
	prog := NewProgram(test.GetBenchmarkPath(), true, nil)
	defer prog.Destroy()
	prog.RollbackAll()
}
