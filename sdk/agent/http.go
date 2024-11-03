package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

const feAddr = ":19898"

const (
	SourcesDir          = "sources"
	SourcesRelativePath = MetaRelativePath + "/" + SourcesDir
	SourcesFileExt      = ".src"
)

var serviceOnce sync.Once

var graphCache *collections.LRUCache[uint64] = collections.NewLRUCache[uint64](1024)

var StartBackend = true

func init() {
	serviceOnce.Do(func() {
		go func() {
			if utils.Go.IsUnittestEnv() || !StartBackend {
				return
			}
			newServer()
		}()
	})
}

func handle404(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := fs.Open(path.Clean(request.URL.Path))
		if os.IsNotExist(err) {
			// 返回index.html的内容
			indexHtml, _ := meta.ReadFile(filepath.ToSlash(filepath.Join(MetaRelativePath, "page/index.html")))
			http.ServeContent(writer, request, "", time.Now(), bytes.NewReader(indexHtml))
			return
		}
		fileServer.ServeHTTP(writer, request)
	})
}

func newServer() {
	// api
	registerApi("/meta/methods", metaMethods)
	registerApi("/traffic/list", trafficList)
	registerApi("/traffic/detail", trafficDetail)
	registerApi("/source/file", sourceFile)
	registerApi("/node/detail", nodeDetail)
	registerApi("/source/tree", sourceTree)
	registerApi("/target/info", targetInfo)

	// page
	staticResources, _ := fs.Sub(meta, filepath.ToSlash(filepath.Join(MetaRelativePath, "page")))
	http.Handle("/", handle404(http.FS(staticResources)))

	// open browser
	autoOpenBrowser()

	// start
	err := http.ListenAndServe(feAddr, nil)
	if err != nil {
		log.Error("elly server start failed.err: %v", err)
	}
}

func autoOpenBrowser() {
	defer func() {
		_ = recover()
	}()
	cmd := map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}[runtime.GOOS]
	if cmd != "" {
		utils.Shell.Exec("", cmd, "http://localhost"+feAddr+"/test")
	}
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
				responseError(response, errors.New(fmt.Sprintf("panic: %v, stack: %v", err, string(debug.Stack()))))
			}
		}()
		header := response.Header()
		header.Set("Content-Type", "application/json")

		handler(response, request)
	}
}

func registerApi(path string, handler func(writer http.ResponseWriter, request *http.Request)) {
	http.HandleFunc("/api"+path, wrapper(handler))
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
	for _, gList := range allGraphs {
		g := mergeGraphs(gList.(*collections.LinkedList[*graph]).Values())
		res = append(res, toTraffic(g, false))
	}
	responseJson(writer, res)
}

func trafficDetail(writer http.ResponseWriter, request *http.Request) {
	// 单个流量明细
	id := queryVal[uint](request, "id")
	gList, ok := graphCache.Get(uint64(id))
	asserts.True(ok)
	g := mergeGraphs(gList.(*collections.LinkedList[*graph]).Values())
	responseJson(writer, toTraffic(g, true))
}

func nodeDetail(writer http.ResponseWriter, request *http.Request) {
	graphId := uint64(queryVal[uint](request, "graphId"))
	nodeId := uint32(queryVal[uint](request, "nodeId"))
	gList, ok := graphCache.Get(graphId)
	asserts.True(ok)
	g := mergeGraphs(gList.(*collections.LinkedList[*graph]).Values())
	n := g.nodes[nodeId]
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

func targetInfo(writer http.ResponseWriter, request *http.Request) {
	info := &TargetInfo{}
	for _, f := range files {
		info.TotalLineNum += f.LineNum
	}
	for _, b := range blocks {
		info.TargetLineNum += b.End.Line - b.Begin.Line + 1
		if globalCovered.Get(uint(b.Id)) {
			info.CoveredLineNum += b.End.Line - b.Begin.Line + 1
		}
	}
	info.CoveredRate = float32(info.CoveredLineNum) / float32(info.TargetLineNum) * 100
	responseJson(writer, info)
}

func sourceTree(writer http.ResponseWriter, request *http.Request) {
	tree := collections.NewSourceTree()

	for _, f := range files {
		n := tree.Add(f.RelativePath, strconv.Itoa(int(f.FileId)))
		n.IsLeaf = true
	}
	responseJson(writer, tree.Root())
}

func sourceFile(writer http.ResponseWriter, request *http.Request) {
	fileId := uint32(queryVal[int](request, "id"))
	lineMap := make(map[int]int)
	for _, b := range blocks {
		if b.FileId != fileId {
			continue
		}
		bCovered := globalCovered.Get(uint(b.Id))
		for i := b.Begin.Line; i <= b.End.Line; i++ {
			if bCovered {
				lineMap[i] = 2
			} else {
				lineMap[i] = 1
			}
		}
	}

	responseJson(writer, map[string]any{
		"code":    string(readCode(fileId)),
		"lineMap": lineMap,
	})
}

///  API结束

// readCode 读取指定代码id文件
func readCode(fileId uint32) []byte {
	content, err := meta.ReadFile(
		filepath.ToSlash(filepath.Join(MetaRelativePath, SourcesDir, files[fileId].RelativePath)) + SourcesFileExt)
	asserts.IsNil(err)
	return content
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
	content, err := json.Marshal(res)
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
	}
	_, _ = writer.Write(content)
}

type CoveredBlock struct {
	Begin Pos `json:"begin"`
	End   Pos `json:"end"`
}

type VarItem struct {
	Idx  int    `json:"idx"`
	Type string `json:"type"`
	Name string `json:"name"`
	Val  any    `json:"val"`
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
	ArgsList      *VarDefList    `json:"args_list"`
	ReturnList    *VarDefList    `json:"return_list"`
	Args          []*VarItem     `json:"args"`
	Returns       []*VarItem     `json:"returns"`
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
		Id:         strconv.Itoa(int(n.methodId)),
		Name:       method.FullName,
		File:       file.RelativePath,
		BlockCnt:   method.BlockCnt,
		Begin:      *method.Begin,
		End:        *method.End,
		Cost:       int32(n.cost),
		ArgsList:   method.ArgsList,
		ReturnList: method.ReturnList,
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
		item.Args = toVarItemList(n.args, method.ArgsList)
		item.Returns = toVarItemList(n.results, method.ReturnList)
	}
	return item
}

func toVarItemList(vars []any, defList *VarDefList) (res []*VarItem) {
	for i, val := range vars {
		res = append(res, &VarItem{
			Idx:  i,
			Type: defList.Type(i),
			Name: defList.Name(i),
			Val:  val,
		})
	}
	return
}

// 合并多个图
func mergeGraphs(list []*graph) *graph {
	if len(list) == 0 {
		return nil
	}
	if len(list) == 1 {
		return list[0]
	}
	res := newGraph(list[0].id)
	origins := make(map[uint64]struct{})
	// 遍历合并多张子图
	for _, g := range list {
		for id, n := range g.nodes {
			old := res.nodes[id]
			if old == nil {
				res.nodes[id] = n
			} else {
				// 现有节点合并块覆盖信息
				asserts.IsNil(old.blocks.Merge(n.blocks))
				old.cost += n.cost
			}
		}
		for k, v := range g.edges {
			res.edges[k] = v
		}
		if g.origin != nil {
			origins[*g.origin] = struct{}{}
		}
		res.time = utils.Min(res.time, g.time)
	}
	// 如果来源节点在图中存在，将来源信息加到边集合
	for origin := range origins {
		from, _ := splitEdge(origin)
		if _, exist := res.nodes[from]; exist {
			res.edges[origin] = struct{}{}
		}
	}
	return res
}

type MethodInfo struct {
	Method  *Method `json:"method"`
	File    string  `json:"file"`
	Package string  `json:"package"`
}

type TargetInfo struct {
	TotalLineNum   int     `json:"totalLineNum"`
	TargetLineNum  int     `json:"targetLineNum"`
	CoveredLineNum int     `json:"coveredLineNum"`
	CoveredRate    float32 `json:"coveredRate"`
}
