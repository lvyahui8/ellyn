package agent

import (
	"embed"
)

// InitAgent 基于元数据初始化agent
func InitAgent(meta embed.FS) Api {
	initAgent(meta)
	return &ellynAgent{}
}

func initAgent(m embed.FS) {
	meta = m
	initConfig()
	initMetaData()
	if len(blocks) > 0 {
		globalCovered = make([]bool, int(len(blocks)))
	}
	initSampling()
}
