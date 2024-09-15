package ellyn

import (
	"embed"
)

//go:embed ellyn_agent ellyn_common
var SdkFs embed.FS

var SdkPaths = []string{"ellyn_agent", "ellyn_common"}

const SdkRawRootPkg = "github.com/lvyahui8/ellyn"
