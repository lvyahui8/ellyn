package utils

import (
	"path"
	"runtime"
	"strings"
)

var File = &fileUtils{}

type fileUtils struct {
}

func (fileUtils) GetProjectRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Dir(path.Dir(path.Dir(b)))
}

func (fileUtils) GetSimpleFileName(file string) {

}

func (fileUtils) FormatFilePath(file string) string {
	return strings.ReplaceAll(file, "\\", "/")
}
