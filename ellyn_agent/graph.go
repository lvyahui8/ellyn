package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"time"
)

type graph struct {
	id uint64
	// time 发生时间
	time   int64
	nodes  map[uint32]*node
	edges  map[uint64]struct{}
	origin *uint64
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

func (g *graph) Merge(o *graph) {
	asserts.Equals(g.id, o.id)
	for id, n := range o.nodes {
		old := g.nodes[id]
		if old == nil {
			g.nodes[id] = n
		} else {
			asserts.IsNil(old.blocks.Merge(n.blocks))
			old.cost += n.cost
		}
	}
	for k, v := range o.edges {
		g.edges[k] = v
	}
	// todo 处理origin的merge
}

func toEdge(from, to uint32) uint64 {
	return uint64(from)<<32 | uint64(to)
}
