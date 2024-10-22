package utils

import (
	"context"
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"
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
		res[cols[0]] = filepath.ToSlash(cols[1])
	}
	return res
}

// GetGoEnv 获取指定的go env
func (g *goUtils) GetGoEnv(workDir, envKey string) string {
	out := Shell.Exec(workDir, "go", "env", envKey)
	return strings.TrimSpace(out)
}

func (g *goUtils) ModTidy(workDir string) {
	Shell.Exec(workDir, "go", "mod", "tidy")
}

// GetModFile 获取go.mod的绝对路径
func (g *goUtils) GetModFile(mainPkgPath string) string {
	return g.GetGoEnv(mainPkgPath, "GOMOD")
}

func (g *goUtils) GetGoRootDir() string {
	wd, err := os.Getwd()
	asserts.IsNil(err)
	return g.GetGoEnv(wd, "GOROOT")
}

func (g *goUtils) IsTestFile(file string) bool {
	return strings.HasSuffix(file, "_test.go")
}

func (g *goUtils) IsSourceFile(file string) bool {
	return strings.HasSuffix(file, ".go") && !g.IsTestFile(file)
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

func (g *goUtils) Build(mainPkgDir string) {
	Shell.Exec(mainPkgDir, "go", "build", "./...")
}

func (g *goUtils) IsUnittestEnv() bool {
	return len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test")
}

type iface struct {
	itab, data uintptr
}

type valueCtx struct {
	context.Context
	key, val interface{}
}

func (g *goUtils) GetKeyValues(ctx context.Context) map[any]any {
	res := make(map[any]any)
	g.getKeyValue(ctx, res)
	return res
}

func (g *goUtils) getKeyValue(ctx context.Context, res map[any]any) {
	iCtx := *(*iface)(unsafe.Pointer(&ctx))
	if iCtx.data == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(iCtx.data))
	if valCtx.Context == nil {
		return
	}
	if valCtx != nil && valCtx.key != nil && valCtx.val != nil {
		res[valCtx.key] = valCtx.val
	}
	g.getKeyValue(valCtx.Context, res)
}
