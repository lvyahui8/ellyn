package main

import (
	"github.com/lvyahui8/ellyn/ellyn_ast"
	"github.com/lvyahui8/ellyn/ellyn_testing"
)

func main() {
	program := ellyn_ast.NewProgram(ellyn_testing.GetTestProjPath(), nil)
	program.Visit()
}
