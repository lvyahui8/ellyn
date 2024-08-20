package guid

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetSvrId(t *testing.T) {
	res := make(map[uint64]struct{})
	for i := 0; i < 1000; i++ {
		res[getSvrId()] = struct{}{}
	}
	t.Log(len(res))
}

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

func TestUint64GUIDGenerator_GenGUID_Concurrent(t *testing.T) {
	g := NewGuidGenerator()
	cnt := 600
	m := collections.NewConcurrentMap(1<<10, func(key any) int {
		return int(key.(uint64))
	})
	gCnt := 10
	w := &sync.WaitGroup{}
	w.Add(gCnt)
	for i := 0; i < gCnt; i++ {
		go func() {
			defer w.Done()
			for i := 0; i < cnt; i++ {
				m.Store(g.GenGUID(), struct{}{})
			}
		}()
	}
	w.Wait()
	require.Equal(t, gCnt*cnt, m.Size())
}

func TestGuidCycleSeq(t *testing.T) {
	generator := NewGuidGenerator()
	generator.cycleSeq()
	gNum := 32
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
	t.Logf("sum = %d\n", sum)
	require.Equal(t, uint64(10263490816), sum)
}

func BenchmarkGuid(b *testing.B) {
	g := NewGuidGenerator()
	for i := 0; i < b.N; i++ {
		g.GenGUID()
	}
}
