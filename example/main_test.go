package main

import (
	"github.com/lvyahui8/ellyn/ellyn_api"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func getNode(g *ellyn_api.Graph, name string) *ellyn_api.Node {
	for _, n := range g.Nodes {
		if strings.Contains(n.MethodName, name) {
			return n
		}
	}
	return nil
}

func TestSum(t *testing.T) {
	ellyn_api.Agent.SetAutoClear(false)
	Sum(1, 2)
	graph := ellyn_api.Agent.GetGraph()
	require.NotNil(t, graph)
	node := getNode(graph, "Sum")
	require.NotNil(t, node)
	require.True(t, len(graph.Edges) == 0)
}

func TestTrade(t *testing.T) {
	ellyn_api.Agent.SetAutoClear(false)
	Trade()
	graph := ellyn_api.Agent.GetGraph()
	require.NotNil(t, graph)
	node := getNode(graph, "Trade")
	require.NotNil(t, node)
}

func TestN(t *testing.T) {
	ellyn_api.Agent.SetAutoClear(false)
	N(4)
	graph := ellyn_api.Agent.GetGraph()
	require.NotNil(t, graph)
	require.True(t, len(graph.Nodes) == 1)
	node := getNode(graph, "N")
	require.NotNil(t, node)
	require.True(t, len(graph.Edges) > 0)
	_, exist := graph.Edges[uint64(node.MethodId)<<32|uint64(node.MethodId)]
	require.True(t, exist)
}

// Run With 'CPU Profiler'
//func TestRunMain(t *testing.T) {
//	go func() {
//		main()
//	}()
//	time.Sleep(10 * time.Second)
//}
