package ellyn_agent

import (
	"errors"
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var currentPath = func() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Dir(b)
}()

type importChecker struct {
	err error
}

func (c *importChecker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		imports := n.Imports
		for _, item := range imports {
			//fmt.Printf("pkg:%s\n", item.Path.Value)
			pkgPath := strings.Trim(item.Path.Value, "\"")
			if !(utils.Go.IsStdPkg(pkgPath) || strings.Contains(pkgPath, "lvyahui8/ellyn")) {
				c.err = errors.New(fmt.Sprintf(
					"SDK does not allow the use of non-standard libraries. pkg %s", pkgPath))
			}
		}
	}
	return nil
}

// TestSdkDependencies 检查分析sdk包的依赖
func TestSdkDependencies(t *testing.T) {
	err := filepath.Walk(currentPath, func(file string, info fs.FileInfo, err error) error {
		if info.IsDir() || len(strings.TrimSpace(file)) == 0 {
			return nil
		}
		if !(strings.HasSuffix(file, ".go") && !strings.HasSuffix(file, "_test.go")) {
			return nil
		}
		content, err := os.ReadFile(file)
		require.Nil(t, err)
		fset := token.NewFileSet()
		parseFile, err := parser.ParseFile(fset, file, content, parser.ImportsOnly)
		require.Nil(t, err)
		checker := &importChecker{}
		ast.Walk(checker, parseFile)
		return checker.err
	})
	require.Nil(t, err)
}
