package utils

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"os"
	"strings"
)

var OS = &osUtils{}

type osUtils struct {
}

func (osUtils) GetWorkDir() string {
	dir, _ := os.Getwd()
	return dir
}

func (osUtils) WriteTo(file string, content []byte) {
	osUtils{}.MkDirs(file[0:strings.LastIndex(file, "/")])
	err := os.WriteFile(file, content, os.ModePerm)
	asserts.IsNil(err)
}

func (osUtils) MkDirs(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	asserts.IsNil(err)
}

func (osUtils) CopyFile(source, target string) {
	file, err := os.ReadFile(source)
	asserts.IsNil(err)
	err = os.WriteFile(target, file, os.ModePerm)
	asserts.IsNil(err)
}
