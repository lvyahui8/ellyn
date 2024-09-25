package ellyn_agent

import (
	"embed"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"
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
		header := response.Header()
		header.Set("Content-Type", "application/json")
		handler(response, request)
	}
}

func register(path string, handler func(writer http.ResponseWriter, request *http.Request)) {
	http.HandleFunc(path, wrapper(handler))
}

func metaMethods(writer http.ResponseWriter, request *http.Request) {
	// 元数据检索配置，配置方法采集，配置mock等
	m, err := json.Marshal(methods)
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
	}
	_, _ = writer.Write(m)
}

func trafficList(writer http.ResponseWriter, request *http.Request) {
	// 流量列表
	allGraphs := graphCache.Values()
	var res []*Traffic
	for _, g := range allGraphs {
		res = append(res, toTraffic(g.(*graph)))
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
	}
	_, _ = writer.Write(bytes)
}

func trafficDetail(writer http.ResponseWriter, request *http.Request) {
	// 单个流量明细
}

func sourceFile(writer http.ResponseWriter, request *http.Request) {
	bytes, err := targetSources.ReadFile(
		filepath.ToSlash(filepath.Join(SourcesDir, files[0].RelativePath)) + SourcesFileExt)
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
	}
	_, _ = writer.Write(bytes)
}

func newServer() {
	register("/meta/methods", metaMethods)
	register("/traffic/list", trafficList)
	register("/traffic/detail", trafficDetail)
	register("/source/0", sourceFile)

	err := http.ListenAndServe(":19898", nil)
	if err != nil {
		log.Error("elly server start failed.err: %v", err)
	}
}

type Node struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	File     string `json:"file"`
	BlockCnt int    `json:"block_cnt"`
	Begin    Pos    `json:"begin"`
	End      Pos    `json:"end"`
}

type Edge struct {
	Source uint32 `json:"source"`
	Target uint32 `json:"target"`
}

type Traffic struct {
	Id    string    `json:"id"`
	Time  time.Time `json:"time"`
	Nodes []*Node   `json:"nodes"`
	Edges []*Edge   `json:"edges"`
}

func toTraffic(g *graph) *Traffic {
	t := &Traffic{}
	t.Id = strconv.FormatUint(g.id, 10) // uint64转成字符串发给前端显示，否则前端会精度丢失
	t.Time = time.UnixMilli(g.time)
	for _, n := range g.nodes {
		method := methods[n.methodId]
		file := files[method.FileId]
		t.Nodes = append(t.Nodes, &Node{
			Id:       n.methodId,
			Name:     method.FullName,
			File:     file.RelativePath,
			BlockCnt: method.BlockCnt,
			Begin:    *method.Begin,
			End:      *method.End,
		})
	}
	for edge := range g.edges {
		t.Edges = append(t.Edges, &Edge{
			Source: uint32(edge >> 32),
			Target: uint32(edge),
		})
	}
	return t
}
