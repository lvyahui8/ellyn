package ellyn_testing

import (
	"testing"
)

func TestGetTestProjPath(t *testing.T) {
	t.Logf(GetTestProjPath())
}

func TestGetRepoRootPath(t *testing.T) {
	t.Log(GetRepoRootPath())
}
