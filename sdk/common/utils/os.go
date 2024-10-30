package utils

import (
	"errors"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"os"
	"path/filepath"
)

var OS = &osUtils{}

type osUtils struct {
}

func (osUtils) GetWorkDir() string {
	dir, _ := os.Getwd()
	return dir
}

func (osUtils) WriteTo(file string, content []byte) {
	osUtils{}.MkDirs(filepath.Dir(file))
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

func (osUtils) Remove(file string) {
	err := os.RemoveAll(file)
	asserts.IsNil(err)
}

func (osUtils) NotExists(file string) bool {
	_, err := os.Stat(file)
	return errors.Is(err, os.ErrNotExist)
}
