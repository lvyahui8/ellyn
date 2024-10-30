package utils

import (
	"github.com/lvyahui8/ellyn/test"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestGoUtils_AllPackages(t *testing.T) {
	pkgMap := Go.AllPackages(test.GetTestProjPath())
	require.NotNil(t, pkgMap)
	require.True(t, len(pkgMap) > 0)
	for pkg, dir := range pkgMap {
		t.Logf("pkg %s, dir %s\n", pkg, dir)
	}
}

func TestGoUtils_GetGoEnv(t *testing.T) {
	modFilePath := Go.GetGoEnv(test.GetTestProjPath(), "GOMOD")
	require.NotEmpty(t, modFilePath)
	t.Log(modFilePath)
}

func TestGoUtils_GetModFile(t *testing.T) {
	modFilePath := Go.GetModFile(test.GetTestProjPath())
	require.Equal(t, filepath.Join(test.GetTestProjPath(), "go.mod"), modFilePath)
}

func TestGoUtils_IsAutoGenFile(t *testing.T) {
	t.Skip()
	goRootDir := Go.GetGoRootDir()
	require.True(t, Go.IsAutoGenFile(goRootDir+"\\src\\cmd\\go\\internal\\test\\flagdefs.go"))
	require.False(t, Go.IsAutoGenFile(goRootDir+"\\src\\cmd\\go\\internal\\test\\cover.go"))
}

// TestGoUtils_IsUnittestEnv
// go test -v
func TestGoUtils_IsUnittestEnv(t *testing.T) {
	require.True(t, Go.IsUnittestEnv())
}
