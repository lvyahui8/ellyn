package utils

import (
	"path"
	"runtime"
)

var File = &fileUtils{}

type fileUtils struct {
}

func (fileUtils) GetProjectRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Dir(path.Dir(path.Dir(b)))
}
