package main

import (
	"github.com/lvyahui8/ellyn/ellyn_api"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSum(t *testing.T) {
	ellyn_api.Agent.SetAutoClear(false)
	Sum(1, 2)
	graph := ellyn_api.Agent.GetGraph()
	require.NotNil(t, graph)
}
