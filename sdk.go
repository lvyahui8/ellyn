package ellyn

import (
	"embed"
)

//go:embed ellyn_agent ellyn_common
var SdkFs embed.FS
