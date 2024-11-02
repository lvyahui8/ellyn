package ellyn

import (
	"embed"
	"path"
	"runtime"
)

const (
	SdkPkgDir  = "sdk"
	SdkPkgPath = "github.com/lvyahui8/ellyn/sdk"
)

const AgentPkg = "ellyn_agent"

const ApiPackage = "github.com/lvyahui8/ellyn/api"

//go:embed sdk
var SdkFs embed.FS

var RepoRootPath = func() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Dir(b)
}()
