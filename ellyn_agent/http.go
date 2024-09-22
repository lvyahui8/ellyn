package ellyn_agent

import (
	"embed"
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"net/http"
)

//go:embed sources
var targetSources embed.FS

func init() {
	go func() {
		newServer()
	}()
}

func newServer() {
	http.HandleFunc("/meta", func(writer http.ResponseWriter, request *http.Request) {
		// 元数据检索配置，配置方法采集，配置mock等
	})
	http.HandleFunc("/list", func(writer http.ResponseWriter, request *http.Request) {
		// 流量列表
	})
	http.HandleFunc("/traffic", func(writer http.ResponseWriter, request *http.Request) {
		// 单个流量明细
	})
	err := http.ListenAndServe(":19898", nil)
	asserts.IsNil(err)
}
