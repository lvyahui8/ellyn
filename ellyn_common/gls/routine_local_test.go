package gls

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestRoutineLocalBasic(t *testing.T) {
	local := NewRoutineLocal()
	local.Set(1)
	val, ok := local.Get()
	require.True(t, ok)
	require.Equal(t, 1, val)
	local.Clear()
	val, ok = local.Get()
	require.False(t, ok)
	require.Nil(t, val)
}

func TestRoutineLocalConcurrent(t *testing.T) {
	local := NewRoutineLocal()
	local.Set(1)
	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		defer w.Done()
		val, ok := local.Get()
		require.False(t, ok)
		require.Nil(t, val)
		local.Set(100)
		val, ok = local.Get()
		require.True(t, ok)
		require.Equal(t, 100, val)
	}()
	w.Wait()
	val, ok := local.Get()
	require.True(t, ok)
	require.Equal(t, 1, val)
}
