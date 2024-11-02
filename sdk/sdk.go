package ellyn_agent

import (
	"embed"
	"github.com/lvyahui8/ellyn/sdk/agent"
)

//go:embed meta
var meta embed.FS

var NotCollected = agent.NotCollected

var Agent = agent.InitAgent(meta)
