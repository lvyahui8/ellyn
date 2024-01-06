package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShellUtils_Exec(t *testing.T) {
	output := Shell.Exec(OS.GetWorkDir(), "echo", "1")
	require.Equal(t, "1\n", output)
}
