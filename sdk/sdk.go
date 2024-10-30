package sdk

import (
	"embed"
	"strings"
)

//go:embed agent common
var SdkFs embed.FS

var SdkPaths = []string{"agent", "common"}

const (
	SdkRawRootPkg = "github.com/lvyahui8/ellyn"
	SdkAgentPkg   = "github.com/lvyahui8/ellyn/sdk/agent"
)

const (
	MetaRelativePath = "ellyn_agent/meta"
	MetaBlocks       = "blocks.dat"
	MetaFiles        = "files.dat"
	MetaMethods      = "methods.dat"
	MetaPackages     = "packages.dat"
)

const RuntimeConfFile = "config.json"

const (
	ApiFile = "ellyn_agent/api.go"
	ApiPkg  = "github.com/lvyahui8/ellyn/api"
)

func IsSdkPkg(pkgPath string) bool {
	for _, sdkPath := range SdkPaths {
		if strings.Contains(pkgPath, sdkPath) {
			return true
		}
	}
	return false
}
