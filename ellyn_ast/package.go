package ellyn_ast

import "strings"

type Package struct {
	// Dir pkg在本地文件系统的绝对路径
	Dir string
	// Name pkg名，如ellyn
	Name string
	// Path Pkg全路径，即写代码时的Import path. 如：github.com/lvyahui8/ellyn
	Path string
}

func NewPackage(dir, path string) Package {
	items := strings.Split(path, "/")
	name := items[len(items)-1]
	return Package{Dir: dir, Name: name, Path: path}
}
