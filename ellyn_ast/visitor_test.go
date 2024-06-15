package ellyn_ast

import (
	"github.com/lvyahui8/ellyn/ellyn_testing"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"
)

func TestFileVisitor_Visit(t *testing.T) {
	fileName := ellyn_testing.GetTestProjPath() + "/main.go"
	content, err := ioutil.ReadFile(fileName)
	require.Nil(t, err)
	visitor := &FileVisitor{content: content}
	fset := token.NewFileSet()
	visitor.fset = fset
	parsedFile, err := parser.ParseFile(fset, fileName, content, parser.ParseComments)
	ast.Walk(visitor, parsedFile)
}
