package sdk

import (
	"embed"
	"github.com/lvyahui8/ellyn/sdk/agent"
)

const AgentPkg = "ellyn_agent"

const (
	MetaRelativePath = ".meta"
	MetaBlocks       = "blocks.dat"
	MetaFiles        = "files.dat"
	MetaMethods      = "methods.dat"
	MetaPackages     = "packages.dat"
)

const RuntimeConfFile = "config.json"

const (
	AgentApiFile = "agent/api.go"
)

//go:embed .meta
var meta embed.FS

var Agent = agent.InitAgent(meta)
