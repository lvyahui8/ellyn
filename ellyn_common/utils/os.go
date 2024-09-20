package utils

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"os"
	"path/filepath"
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
	content, err := os.ReadFile(source)
	asserts.IsNil(err)
	osUtils{}.MkDirs(filepath.Dir(target))
	err = os.WriteFile(target, content, os.ModePerm)
	asserts.IsNil(err)
}
