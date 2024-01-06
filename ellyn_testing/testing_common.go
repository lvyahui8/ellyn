package ellyn_testing

import (
	"path"
	"runtime"
)

const TestProjPath = ""

func GetTestProjPath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(path.Dir(b)), "test_proj")
}
