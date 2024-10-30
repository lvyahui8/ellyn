package test

import (
	"path"
	"runtime"
)

const TestProjPath = ""

func GetTestProjPath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(path.Dir(b)), "example")
}

func GetBenchmarkPath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(path.Dir(b)), "benchmark")
}

func GetRepoRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Dir(path.Dir(b))
}
