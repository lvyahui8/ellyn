package agent

import (
	"embed"
)

func initAgent(m embed.FS) {
	meta = m
	initConfig()
	initMetaData()
	if len(blocks) > 0 {
		globalCovered = make([]bool, int(len(blocks)))
	}
	initSampling()
}
