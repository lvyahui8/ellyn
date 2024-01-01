package utils

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestGoUtils_AllPackages(t *testing.T) {
	pkgMap := Go.AllPackages(TestProjPath)
	require.NotNil(t, pkgMap)
	require.True(t, len(pkgMap) > 0)
	for pkg, dir := range pkgMap {
		t.Logf("pkg %s, dir %s\n", pkg, dir)
	}
}

func TestGoUtils_GetGoEnv(t *testing.T) {
	modFilePath := Go.GetGoEnv(TestProjPath, "GOMOD")
	require.NotEmpty(t, modFilePath)
	t.Log(modFilePath)
}

func TestGoUtils_GetModFile(t *testing.T) {
	modFilePath := Go.GetModFile(TestProjPath)
	require.Equal(t, filepath.Join(TestProjPath, "go.mod"), modFilePath)
}

func TestGoUtils_GetRootPkg(t *testing.T) {
	rootPkg := Go.GetRootPkg(Go.GetModFile(TestProjPath))
	require.Equal(t, "test_proj", rootPkg)
}
