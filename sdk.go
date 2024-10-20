package ellyn

import (
	"embed"
	"strings"
)

//go:embed ellyn_agent ellyn_common
var SdkFs embed.FS

var SdkPaths = []string{"ellyn_agent", "ellyn_common"}

const (
	SdkRawRootPkg = "github.com/lvyahui8/ellyn"
	SdkAgentPkg   = "github.com/lvyahui8/ellyn/ellyn_agent"
)

const (
	MetaRelativePath = "ellyn_agent/meta"
	MetaBlocks       = "blocks.dat"
	MetaFiles        = "files.dat"
	MetaMethods      = "methods.dat"
	MetaPackages     = "packages.dat"
)

const (
	ApiFile = "ellyn_agent/api.go"
	ApiPkg  = "github.com/lvyahui8/ellyn/ellyn_api"
)

func IsSdkPkg(pkgPath string) bool {
	for _, sdkPath := range SdkPaths {
		if strings.Contains(pkgPath, sdkPath) {
			return true
		}
	}
	return false
}
