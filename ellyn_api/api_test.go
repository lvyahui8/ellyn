package ellyn_api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAgent(t *testing.T) {
	require.Nil(t, Agent.GetGraph())
	Agent.Init(nil)
	require.Nil(t, Agent.target)
	Agent.SetAutoClear(false)
}
