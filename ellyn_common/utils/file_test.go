package utils

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
	"testing"
)

var TestProjPath = filepath.Join(File.GetProjectRootPath(), "test_proj")

func TestFileUtils_GetProjectRootPath(t *testing.T) {
	require.True(t, strings.HasSuffix(File.GetProjectRootPath(), "/ellyn"))
}
