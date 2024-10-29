package ellyn_agent

import (
	_ "embed"
	"encoding/json"
)

//go:embed meta/config.json
var configContent []byte

var conf Configuration

func configInit() {
	if len(configContent) > 0 {
		err := json.Unmarshal(configContent, &conf)
		if err != nil {
			log.Error("config init failed. err %v", err)
		}
	}
}

type Configuration struct {
	// 是否采集参数
	NoArgs bool

	// 是否收集演示数据
	NoDemo bool
}
