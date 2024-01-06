package utils

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestFileUtils_GetProjectRootPath(t *testing.T) {
	require.True(t, strings.HasSuffix(File.GetProjectRootPath(), "/ellyn"))
}
