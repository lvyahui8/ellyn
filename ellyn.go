package ellyn

import "embed"

const (
	SdkDir        = "sdk"
	SdkRawRootPkg = "github.com/lvyahui8/ellyn/sdk"
)

const AgentPkg = "ellyn_agent"

const ApiPackage = "github.com/lvyahui8/ellyn/api"

//go:embed sdk
var SdkFs embed.FS
