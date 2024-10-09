package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"time"
)

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
		time:  time.Now().UnixMilli(),
		nodes: make(map[uint32]*node),
		edges: make(map[uint64]struct{}),
	}
}

func (g *graph) add(from *methodFrame, to *methodFrame) {
	toNode := g.draw(to)
	if from != nil {
		fromNode := g.draw(from)
		g.addEdges(fromNode, toNode)
	}
	if to.recursion {
		g.addEdges(toNode, toNode)
	}
}

func (g *graph) draw(f *methodFrame) *node {
	cost := time.Now().UnixMilli() - f.begin
	if n, ok := g.nodes[f.methodId]; ok {
		n.cost += cost // 累计耗时、取最大值
		err := n.blocks.Merge(f.blocks)
		asserts.IsNil(err)
		return n
	} else {
		n = &node{
			methodId: f.methodId,
			blocks:   f.blocks,
			cost:     cost,
		}
		g.nodes[f.methodId] = n
		return n
	}
}

func (g *graph) addEdges(from *node, to *node) {
	g.edges[toEdge(from.methodId, to.methodId)] = struct{}{}
}

func toEdge(from, to uint32) uint64 {
	return uint64(from)<<32 | uint64(to)
}

func splitEdge(edge uint64) (from, to uint32) {
	to = uint32(edge)
	from = uint32(edge >> 32)
	return
}
