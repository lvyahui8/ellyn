package ellyn_agent

import (
	"embed"
	"encoding/json"
	"errors"
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
	responseJson(writer, methods)
}

func trafficList(writer http.ResponseWriter, request *http.Request) {
	// 流量列表
	allGraphs := graphCache.Values()
	var res []*Traffic
	for _, g := range allGraphs {
		res = append(res, toTraffic(g.(*graph), false))
	}
	responseJson(writer, res)
}

func trafficDetail(writer http.ResponseWriter, request *http.Request) {
	// 单个流量明细
	query := request.URL.Query()
	idStr := query.Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		responseError(writer, err)
		return
	}
	g, ok := graphCache.Get(id)
	if !ok {
		responseError(writer, errors.New("traffic not found"))
		return
	}
	responseJson(writer, toTraffic(g.(*graph), true))
}

func sourceFile(writer http.ResponseWriter, request *http.Request) {
	bytes, err := targetSources.ReadFile(
		filepath.ToSlash(filepath.Join(SourcesDir, files[0].RelativePath)) + SourcesFileExt)
	if err != nil {
		responseError(writer, err)
		return
	}
	_, _ = writer.Write(bytes)
}

func responseError(writer http.ResponseWriter, err error) {
	_, _ = writer.Write([]byte(err.Error()))
}

func responseJson(writer http.ResponseWriter, res any) {
	bytes, err := json.Marshal(res)
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

type CoveredBlock struct {
	Begin Pos `json:"begin"`
	End   Pos `json:"end"`
}

type Node struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	File          string         `json:"file"`
	BlockCnt      int            `json:"block_cnt"`
	Begin         Pos            `json:"begin"`
	End           Pos            `json:"end"`
	CoveredBlocks []CoveredBlock `json:"covered_blocks"`
	CoveredRate   float32        `json:"covered_rate"`
	HasErr        bool           `json:"has_err"`
	Cost          int32          `json:"cost"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Traffic struct {
	Id    string    `json:"id"`
	Time  time.Time `json:"time"`
	Nodes []*Node   `json:"nodes"`
	Edges []*Edge   `json:"edges"`
}

func toTraffic(g *graph, withDetail bool) *Traffic {
	t := &Traffic{
		Nodes: make([]*Node, 0),
		Edges: make([]*Edge, 0),
	}
	t.Id = strconv.FormatUint(g.id, 10) // uint64转成字符串发给前端显示，否则前端会精度丢失
	t.Time = time.UnixMilli(g.time)
	for _, n := range g.nodes {
		method := methods[n.methodId]
		file := files[method.FileId]
		item := &Node{
			Id:       strconv.Itoa(int(n.methodId)),
			Name:     method.FullName,
			File:     file.RelativePath,
			BlockCnt: method.BlockCnt,
			Begin:    *method.Begin,
			End:      *method.End,
		}
		if withDetail {
			coveredNum := 0
			for _, block := range method.Blocks {
				if n.blocks.Get(uint(block.MethodOffset)) {
					coveredNum += block.End.Line - block.Begin.Line + 1
					item.CoveredBlocks = append(item.CoveredBlocks, CoveredBlock{
						Begin: *block.Begin,
						End:   *block.End,
					})
				}
			}
			item.CoveredRate = float32(coveredNum) / float32(method.End.Line-method.Begin.Line+1) * 100
		}
		t.Nodes = append(t.Nodes, item)
	}
	for edge := range g.edges {
		t.Edges = append(t.Edges, &Edge{
			Source: strconv.Itoa(int(uint32(edge >> 32))),
			Target: strconv.Itoa(int(uint32(edge))),
		})
	}
	return t
}
