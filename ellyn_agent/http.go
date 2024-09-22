package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"net/http"
)

func newServer() {
	http.HandleFunc("/list", func(writer http.ResponseWriter, request *http.Request) {
		// 流量列表
	})
	http.HandleFunc("/traffic", func(writer http.ResponseWriter, request *http.Request) {
		// 单个流量明细
	})
	err := http.ListenAndServe(":19898", nil)
	asserts.IsNil(err)
}
