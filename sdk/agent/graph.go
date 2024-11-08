package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"sync"
)

var graphPool = &sync.Pool{New: func() any {
	return newGraph(0)
}}

var emptyVal = struct{}{}

type graphGroup struct {
	list *collections.LinkedList[*graph]
}

func (g *graphGroup) Recycle() {
	for _, subGraph := range g.list.Values() {
		subGraph.Recycle()
	}
}

type graph struct {
	id uint64
	// time 发生时间
	time int64
	// origin 调用链触发来源，当异步触发时有此字段值
	origin *uint64
	nodes  map[uint32]*node
	edges  map[uint64]struct{}
}

func newGraph(id uint64) *graph {
	return &graph{
		id:    id,
		time:  currentTime().UnixMilli(),
		nodes: make(map[uint32]*node, 8),
		edges: make(map[uint64]struct{}, 8),
	}
}

func (g *graph) Recycle() {
	//log.Info("recycle g:%d", g.id)
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

func (g *graph) addEdge(from, to uint32) {
	g.edges[toEdge(from, to)] = emptyVal
}

func toEdge(from, to uint32) uint64 {
	return uint64(from)<<32 | uint64(to)
}

func splitEdge(edge uint64) (from, to uint32) {
	to = uint32(edge)
	from = uint32(edge >> 32)
	return
}
