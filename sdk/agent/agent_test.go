package agent

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAgentPanic(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	defer handleSelfError()()
	items := []int{0, 1, 2, 3}
	t.Log(items[6])
}
