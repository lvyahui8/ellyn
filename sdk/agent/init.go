package agent

import (
	"embed"
	"github.com/lvyahui8/ellyn/sdk/common/collections"
)

func initAgent(m embed.FS) {
	meta = m
	initConfig()
	initMetaData()
	if len(blocks) > 0 {
		globalCovered = collections.NewBitMap(uint(len(blocks)))
	}
}
