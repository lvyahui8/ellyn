package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_testing"
	"testing"
)

func TestProgram(t *testing.T) {
	prog := NewProgram(ellyn_testing.GetTestProjPath(), nil)
	prog.Visit()
}
