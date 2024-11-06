package main

import (
	_ "github.com/lvyahui8/ellyn"
	"github.com/lvyahui8/ellyn/api"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func getNode(g *api.Graph, name string) *api.Node {
	for _, n := range g.Nodes {
		if strings.Contains(n.MethodName, name) {
			return n
		}
	}
	return nil
}

func TestAutoClean(t *testing.T) {
	Sum(1, 2)
	graph := api.Agent.GetGraph()
	require.NotNil(t, graph.Nodes)
	require.True(t, len(graph.Nodes) == 0)
}

func TestSum(t *testing.T) {
	api.Agent.SetAutoClear(false)
	Sum(1, 2)
	graph := api.Agent.GetGraph()
	require.NotNil(t, graph)
	node := getNode(graph, "Sum")
	require.NotNil(t, node)
	require.True(t, len(graph.Edges) == 0)
}

func TestTrade(t *testing.T) {
	api.Agent.SetAutoClear(false)
	Trade()
	graph := api.Agent.GetGraph()
	require.NotNil(t, graph)
	node := getNode(graph, "Trade")
	require.NotNil(t, node)
}

func TestN(t *testing.T) {
	api.Agent.SetAutoClear(false)
	N(4)
	graph := api.Agent.GetGraph()
	require.NotNil(t, graph)
	require.True(t, len(graph.Nodes) == 1)
	node := getNode(graph, "N")
	require.NotNil(t, node)
	// 递归调用，存在递归边
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
