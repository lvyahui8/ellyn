package ellyn_agent

import (
	"context"
	"embed"
	"github.com/lvyahui8/ellyn/sdk/agent"
)

//go:embed meta
var meta embed.FS

var NotCollected = agent.NotCollected
var NilGoCtx *context.Context = nil

var Agent = agent.InitAgent(meta)
