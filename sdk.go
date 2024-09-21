package ellyn

import (
	"embed"
)

//go:embed ellyn_agent ellyn_common
var SdkFs embed.FS

var SdkPaths = []string{"ellyn_agent", "ellyn_common"}

const SdkRawRootPkg = "github.com/lvyahui8/ellyn"

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
