package goroutine

import (
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRoutinePool(t *testing.T) {
	w := sync.WaitGroup{}
	cnt := 100
	w.Add(cnt)
	pool := NewRoutinePool(10)
	var sum int64
	for i := 0; i < cnt; i++ {
		pool.Submit(func() {
			defer w.Done()
			t.Logf("goid:%d execute\n", GetGoId())
			atomic.AddInt64(&sum, 1)
		})
	}
	w.Wait()
	require.Equal(t, int64(cnt), sum)
}
