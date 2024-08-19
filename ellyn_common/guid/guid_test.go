package guid

import (
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestUint64GUIDGenerator_GenGUID(t *testing.T) {
	generator := NewGuidGenerator()
	cnt := 100
	idMap := map[uint64]struct{}{}
	for i := 0; i < cnt; i++ {
		id := generator.GenGUID()
		t.Log(id)
		idMap[id] = struct{}{}
	}
	require.Equal(t, cnt, len(idMap))
}

func TestGuidCycleSeq(t *testing.T) {
	generator := NewGuidGenerator()
	generator.cycleSeq()
	gNum := runtime.NumCPU() << 1
	w := &sync.WaitGroup{}
	w.Add(gNum)
	var sum uint64 = 0
	k := 10000
	for i := 0; i < gNum; i++ {
		go func() {
			defer w.Done()
			for i := 0; i < k; i++ {
				seq := generator.cycleSeq()
				atomic.AddUint64(&sum, seq)
			}
		}()
	}
	w.Wait()
	require.Equal(t, uint64(10263490816), sum)
}

func BenchmarkGuid(b *testing.B) {
	g := NewGuidGenerator()
	for i := 0; i < b.N; i++ {
		g.GenGUID()
	}
}
