package ellyn_agent

import "github.com/lvyahui8/ellyn/ellyn_common/asserts"

type graph struct {
	nodes map[uint32]*node
	edges map[uint64]struct{}
}

func (g *graph) add(from *methodFrame, to *methodFrame) {
	fromNode := g.draw(from)
	toNode := g.draw(to)
	g.addEdges(fromNode, toNode)
	if to.recursion {
		g.addEdges(toNode, toNode)
	}
}

func (g *graph) draw(f *methodFrame) *node {
	if n, ok := g.nodes[f.methodId]; ok {
		err := n.blocks.Merge(f.blocks)
		asserts.IsNil(err)
		return n
	} else {
		n = &node{
			methodId: f.methodId,
			blocks:   f.blocks,
		}
		g.nodes[f.methodId] = n
		return n
	}
}

func (g *graph) addEdges(from *node, to *node) {
	g.edges[uint64(from.methodId)<<32|uint64(to.methodId)] = struct{}{}
}
