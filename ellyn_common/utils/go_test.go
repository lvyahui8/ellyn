package utils

import (
	"github.com/lvyahui8/ellyn/ellyn_testing"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestGoUtils_AllPackages(t *testing.T) {
	pkgMap := Go.AllPackages(ellyn_testing.GetTestProjPath())
	require.NotNil(t, pkgMap)
	require.True(t, len(pkgMap) > 0)
	for pkg, dir := range pkgMap {
		t.Logf("pkg %s, dir %s\n", pkg, dir)
	}
}

func TestGoUtils_GetGoEnv(t *testing.T) {
	modFilePath := Go.GetGoEnv(ellyn_testing.GetTestProjPath(), "GOMOD")
	require.NotEmpty(t, modFilePath)
	t.Log(modFilePath)
}

func TestGoUtils_GetModFile(t *testing.T) {
	modFilePath := Go.GetModFile(ellyn_testing.GetTestProjPath())
	require.Equal(t, filepath.Join(ellyn_testing.GetTestProjPath(), "go.mod"), modFilePath)
}

func TestGoUtils_GetRootPkg(t *testing.T) {
	rootPkg := Go.GetRootPkgPath(Go.GetModFile(ellyn_testing.GetTestProjPath()))
	require.Equal(t, "test_proj", rootPkg)
}

func TestGoUtils_IsAutoGenFile(t *testing.T) {
	require.True(t, Go.IsAutoGenFile("D:\\software\\go\\go1.18\\src\\cmd\\go\\internal\\test\\flagdefs.go"))
	require.False(t, Go.IsAutoGenFile("D:\\software\\go\\go1.18\\src\\cmd\\go\\internal\\test\\cover.go"))
}
