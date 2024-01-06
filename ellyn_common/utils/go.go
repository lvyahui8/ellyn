package utils

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var Go *goUtils = &goUtils{}

type goUtils struct {
}

// AllPackages 返回pkgPath -> pkgDir
func (g *goUtils) AllPackages(mainPkgPath string) map[string]string {
	out := Shell.Exec(mainPkgPath, "go", "list", "-test=false", "-deps=true", "-f", "{{.ImportPath}}|#|{{.Dir}}")
	lines := strings.Split(out, "\n")
	res := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cols := strings.Split(line, "|#|")
		res[cols[0]] = cols[1]
	}
	return res
}

// GetGoEnv 获取指定的go env
func (g *goUtils) GetGoEnv(mainPkgPath, envKey string) string {
	out := Shell.Exec(mainPkgPath, "go", "env", envKey)
	return strings.TrimSpace(out)
}

// GetModFile 获取go.mod的绝对路径
func (g *goUtils) GetModFile(mainPkgPath string) string {
	return g.GetGoEnv(mainPkgPath, "GOMOD")
}

// GetRootPkgPath 获取项目go.mod文件所在的package name
func (g *goUtils) GetRootPkgPath(modFilePath string) string {
	content, err := os.ReadFile(modFilePath)
	asserts.IsNil(err)
	modFile, err := modfile.Parse("go.mod", content, nil)
	asserts.IsNil(err)
	return modFile.Module.Mod.Path
}

func (g *goUtils) IsTestFile(file string) bool {
	return strings.HasSuffix(file, "_test.go")
}

func (g *goUtils) IsAutoGenFile(file string) bool {
	content, err := ioutil.ReadFile(file)
	asserts.IsNil(err)
	return g.IsAutoGenContent(content)
}

var autoGenMarkReg = regexp.MustCompile("[\r\n]// Code generated .* DO NOT EDIT.[\r\n]")

func (g *goUtils) IsAutoGenContent(content []byte) bool {
	idx := strings.Index(string(content), "package")
	if idx < 0 {
		return false
	}
	head := content[0:idx]
	return autoGenMarkReg.Match(head)
}
