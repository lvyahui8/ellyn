package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/ctime"
	"sync"
)

// graphPool 链路池
var graphPool = &sync.Pool{New: func() any {
	return newGraph(0)
}}

// graphGroup 链路组，id相同的链路（即多个异步协程产生的链路）存放在同一个group中。
type graphGroup struct {
	list *collections.LinkedList[*graph]
}

// Recycle 重置&回收所有子链路
func (g *graphGroup) Recycle() {
	for _, subGraph := range g.list.Values() {
		subGraph.Recycle()
	}
}

// graph 链路图，用于描述调用链路
// 链路图并非100%准确，对于有递归或者循环调用的场景，只会记录一次边，运行时数据也只保留第一次调用的数据
type graph struct {
	// 链路id
	id uint64
	// time 发生时间，单位毫秒
	time int64
	// origin 调用链触发来源方法id，当异步触发时有此字段值
	origin *uint64
	// nodes 节点（方法）
	nodes map[uint32]*node
	// edges 边（调用关系），这里用key存储边，高32位为from、低32位为to
	edges map[uint64]struct{}
}

// newGraph 初始化一个空的graph
func newGraph(id uint64) *graph {
	return &graph{
		id:    id,
		time:  ctime.Current().UnixMilli(),
		nodes: make(map[uint32]*node, 10),
		edges: make(map[uint64]struct{}, 10),
	}
}

// Recycle 重置&回收Graph
func (g *graph) Recycle() {
	//logging.Info("recycle g:%d", g.id)
	g.origin = nil
	for k, n := range g.nodes {
		delete(g.nodes, k)
		n.Recycle()
	}
	for k := range g.edges {
		delete(g.edges, k)
	}
	graphPool.Put(g)
}

// addEdge往graph中添加一条边
func (g *graph) addEdge(from, to uint32) {
	g.edges[toEdge(from, to)] = struct{}{}
}

// toEdge 指定from to方法id生成边
func toEdge(from, to uint32) uint64 {
	return uint64(from)<<32 | uint64(to)
}

// splitEdge 将边还原成from、to
func splitEdge(edge uint64) (from, to uint32) {
	to = uint32(edge)
	from = uint32(edge >> 32)
	return
}
