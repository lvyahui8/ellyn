package ellyn_ast

import (
	"github.com/lvyahui8/ellyn"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"github.com/lvyahui8/ellyn/ellyn_testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProgram(t *testing.T) {
	t.Skip()
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	prog.Visit()
	utils.Go.Build(prog.mainPkg.Dir)
	prog.rollbackAll()
}

func TestExample(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	prog := NewProgram(ellyn_testing.GetTestProjPath())
	prog.specifySdkDir = ellyn_testing.GetRepoRootPath()
	prog.sdkImportPkgPath = ellyn.SdkAgentPkg
	prog.Visit()
}
