package agent

import (
	"encoding/json"
	"path/filepath"
)

var conf Configuration

func initConfig() {
	configContent, _ := meta.ReadFile(filepath.ToSlash(filepath.Join(MetaRelativePath, RuntimeConfFile)))
	if len(configContent) > 0 {
		err := json.Unmarshal(configContent, &conf)
		if err != nil {
			log.Error("config init failed. err %v", err)
		}
	}
}

func NewDefaultConf() *Configuration {
	return &Configuration{
		NoArgs:       false,
		NoDemo:       false,
		SamplingRate: 1,
	}
}

type Configuration struct {
	// 是否采集参数
	NoArgs bool

	// 是否收集演示数据
	NoDemo bool

	// 采样率， 0-1(0%-100%)
	SamplingRate float64
}
