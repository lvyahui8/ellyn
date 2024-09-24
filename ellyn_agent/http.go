package ellyn_agent

import (
	"embed"
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"
)

//go:embed sources
var targetSources embed.FS

const (
	SourcesDir          = "sources"
	SourcesRelativePath = "ellyn_agent/sources"
	SourcesFileExt      = ".src"
)

var serviceOnce sync.Once

func init() {
	serviceOnce.Do(func() {
		go func() {
			newServer()
		}()
	})
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func wrapper(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		setupCORS(&response, request)
		if request.Method == "OPTIONS" {
			return
		}
		handler(response, request)
	}
}

func register(path string, handler func(writer http.ResponseWriter, request *http.Request)) {
	http.HandleFunc(path, wrapper(handler))
}

func newServer() {
	register("/meta/methods", func(writer http.ResponseWriter, request *http.Request) {
		// 元数据检索配置，配置方法采集，配置mock等
		header := writer.Header()
		header.Set("Content-Type", "application/json")
		m, err := json.Marshal(methods)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
		}
		_, _ = writer.Write(m)
	})
	register("/list", func(writer http.ResponseWriter, request *http.Request) {
		// 流量列表
	})
	register("/traffic", func(writer http.ResponseWriter, request *http.Request) {
		// 单个流量明细
	})
	register("/source/0", func(writer http.ResponseWriter, request *http.Request) {
		bytes, err := targetSources.ReadFile(filepath.ToSlash(filepath.Join(SourcesDir, files[0].RelativePath)) + SourcesFileExt)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
		}
		_, _ = writer.Write(bytes)
	})
	err := http.ListenAndServe(":19898", nil)
	if err != nil {
		log.Error("elly server start failed.err: %v", err)
	}
}
