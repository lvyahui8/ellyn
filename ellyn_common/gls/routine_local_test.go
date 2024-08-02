package gls

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestRoutineLocalBasic(t *testing.T) {
	local := RoutineLocal[int]{}
	local.Set(1)
	val, ok := local.Get()
	require.True(t, ok)
	require.Equal(t, 1, val)
	localPtr := &local
	val, ok = localPtr.Get()
	require.True(t, ok)
	require.Equal(t, 1, val)
	localPtr.Set(2)
	val, ok = local.Get()
	require.True(t, ok)
	require.Equal(t, 2, val)
	local.Clear()
	val, ok = local.Get()
	require.False(t, ok)
	local.Set(4)
	val, ok = local.Get()
	require.True(t, ok)
	require.Equal(t, 4, val)
}

func TestRoutineLocalConcurrent(t *testing.T) {
	local := &RoutineLocal[int]{}
	local.Set(1)
	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		defer w.Done()
		val, ok := local.Get()
		require.False(t, ok)
		require.Zero(t, val)
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

// go test -v -run ^$  -bench 'BenchmarkRoutineLocal/routineLocal' -benchtime=5s -benchmem -cpuprofile profile.pprof
// go tool pprof -http=":8081" profile.pprof
func BenchmarkRoutineLocal(b *testing.B) {
	local := RoutineLocal[int]{}
	for i := 0; i < 1000; i++ { // 1000个RoutineLocal，相当于context.Context中放了1000个key。
		other := RoutineLocal[int]{}
		other.Set(i)
	}
	ctx := context.WithValue(context.Background(), "test", "val")
	// 当ctx有10多个key时，ctx.Value方法实际是一个链表查找，性能还不如routineLocal
	ctx10Keys := ctx
	for i := 0; i < 10; i++ {
		ctx10Keys = context.WithValue(ctx10Keys, i, i)
	}
	ctx20Keys := ctx
	for i := 0; i < 20; i++ {
		ctx20Keys = context.WithValue(ctx20Keys, i, i)
	}
	b.Run("routineLocal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = local.Get()
		}
	})
	b.Run("ctxGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ctx.Value("test")
		}
	})
	b.Run("ctx10KeysGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ctx10Keys.Value(3)
		}
	})
	b.Run("ctx20KeysGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ctx20Keys.Value(3)
		}
	})
}
