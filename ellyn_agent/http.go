package ellyn_agent

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
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
		defer func() {
			err := recover()
			if err != nil {
				responseError(response, errors.New(fmt.Sprintf("panic: %v", err)))
			}
		}()
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
	var res []*MethodInfo
	for _, m := range methods {
		file := files[m.FileId]
		pkg := packages[m.PackageId]
		res = append(res, &MethodInfo{
			Method:  m,
			File:    file.RelativePath,
			Package: pkg.Path,
		})
	}
	responseJson(writer, res)
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
	id := queryVal[uint](request, "id")
	g, ok := graphCache.Get(uint64(id))
	if !ok {
		responseError(writer, errors.New("traffic not found"))
		return
	}
	responseJson(writer, toTraffic(g.(*graph), true))
}

func sourceFile(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write(readCode(0))
}

func readCode(fileId uint32) []byte {
	bytes, err := targetSources.ReadFile(
		filepath.ToSlash(filepath.Join(SourcesDir, files[fileId].RelativePath)) + SourcesFileExt)
	asserts.IsNil(err)
	return bytes
}

func nodeDetail(writer http.ResponseWriter, request *http.Request) {
	graphId := uint64(queryVal[uint](request, "graphId"))
	nodeId := uint32(queryVal[uint](request, "nodeId"))
	g, ok := graphCache.Get(graphId)
	asserts.True(ok)
	n := (g.(*graph)).nodes[nodeId]
	resNode := transferNode(n, true)
	mtd := methods[n.methodId]
	code := readCode(mtd.FileId)
	//funcCode := string(code[mtd.Begin.Offset:mtd.End.Offset])
	funcCode := strings.Join(strings.Split(string(code), "\n")[mtd.Begin.Line-1:mtd.End.Line], "\n")
	responseJson(writer, map[string]any{
		"resNode":  resNode,
		"funcCode": funcCode,
	})
}

// queryVal 工具方法
func queryVal[T string | int | uint](request *http.Request, key string) T {
	query := request.URL.Query()
	valStr := query.Get(key)

	v := reflect.ValueOf(new(T))

	switch v.Type().Elem().Kind() {
	case reflect.String:
		v.Elem().Set(reflect.ValueOf(valStr))
	case reflect.Int:
		val, err := strconv.Atoi(valStr)
		asserts.IsNil(err)
		v.Elem().Set(reflect.ValueOf(val))
	case reflect.Uint:
		val, err := strconv.ParseUint(valStr, 10, 64)
		asserts.IsNil(err)
		v.Elem().Set(reflect.ValueOf(uint(val)))
	default:
		panic("invalid type")
	}
	return v.Elem().Interface().(T)
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
	register("/node/detail", nodeDetail)

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
		t.Nodes = append(t.Nodes, transferNode(n, withDetail))
	}
	for edge := range g.edges {
		t.Edges = append(t.Edges, &Edge{
			Source: strconv.Itoa(int(uint32(edge >> 32))),
			Target: strconv.Itoa(int(uint32(edge))),
		})
	}
	return t
}

func transferNode(n *node, withDetail bool) *Node {
	method := methods[n.methodId]
	file := files[method.FileId]
	item := &Node{
		Id:       strconv.Itoa(int(n.methodId)),
		Name:     method.FullName,
		File:     file.RelativePath,
		BlockCnt: method.BlockCnt,
		Begin:    *method.Begin,
		End:      *method.End,
		Cost:     int32(n.cost),
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
	return item
}

type MethodInfo struct {
	Method  *Method `json:"method"`
	File    string  `json:"file"`
	Package string  `json:"package"`
}
