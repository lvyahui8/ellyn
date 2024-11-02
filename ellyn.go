package ellyn

import "embed"

const (
	SdkPkgDir        = "sdk"
	SdkPkgPathPrefix = "github.com/lvyahui8/ellyn/sdk"
)

const AgentPkg = "ellyn_agent"

const ApiPackage = "github.com/lvyahui8/ellyn/api"

//go:embed sdk
var SdkFs embed.FS
