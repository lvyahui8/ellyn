package instr

import (
	"flag"
	"github.com/lvyahui8/ellyn/sdk/agent"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"github.com/lvyahui8/ellyn/test"
	"github.com/stretchr/testify/require"
	"strconv"
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
	args := flag.Args()
	rate := 0.0001
	if len(args) == 1 {
		val, err := strconv.ParseFloat(args[0], 10)
		require.Nil(t, err)
		require.True(t, val >= 0 && val <= 1)
		rate = val
	}
	t.Logf("samplingRate: %f\n", rate)
	prog := NewProgram(test.GetBenchmarkPath(), true, &agent.Configuration{
		NoArgs:       true,
		NoDemo:       true,
		SamplingRate: rate,
	})
	defer prog.Destroy()
	prog.RollbackAll()
	prog.Visit()

	t.Log()
}

func TestRollbackBenchmark(t *testing.T) {
	prog := NewProgram(test.GetBenchmarkPath(), true, nil)
	defer prog.Destroy()
	prog.RollbackAll()
}
