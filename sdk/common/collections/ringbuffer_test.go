package collections

import (
	"github.com/stretchr/testify/require"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

//var buffer =

const capacity = 100000

type channelQueue[T any] struct {
	q chan T
}

func newChannelQueue[T any](capacity int) *channelQueue[T] {
	return &channelQueue[T]{
		q: make(chan T, capacity),
	}
}

func (c channelQueue[T]) Enqueue(value T) (success bool) {
	c.q <- value
	return true
}

func (c channelQueue[T]) Dequeue() (value T, success bool) {
	select {
	case value, success = <-c.q:
	default:
	}
	return
}

func (c channelQueue[T]) Close() {
	close(c.q)
}

func TestChannelQueue(t *testing.T) {
	queue := newChannelQueue[int](100)
	for i := 0; i < 10; i++ {
		total, produceCnt, consumeCnt := queueReadWrite(queue, 1000, 10, 10, true)
		// 阻塞队列一定会写入成功、消费成功，不会丢元素
		require.Equal(t, total, produceCnt)
		require.Equal(t, produceCnt, consumeCnt)
	}
}

func TestMaxIdx(t *testing.T) {
	s := math.MaxInt64 / int64(100*10000)
	year := s / int64(60*60*24*365)
	t.Log(year)
}

func TestRingBufferFull(t *testing.T) {
	buffer := NewRingBuffer[int](4)
	success := buffer.Enqueue(1)
	require.True(t, success)
	success = buffer.Enqueue(2)
	require.True(t, success)
	success = buffer.Enqueue(3)
	require.True(t, success)
	success = buffer.Enqueue(4)
	require.True(t, success)
	success = buffer.Enqueue(5)
	require.False(t, success)
	success = buffer.Enqueue(6)
	require.False(t, success)
	value, success := buffer.Dequeue()
	require.True(t, success)
	require.Equal(t, 1, value)
	value, success = buffer.Dequeue()
	require.True(t, success)
	require.Equal(t, 2, value)
	value, success = buffer.Dequeue()
	require.True(t, success)
	require.Equal(t, 3, value)
	value, success = buffer.Dequeue()
	require.True(t, success)
	require.Equal(t, 4, value)
	value, success = buffer.Dequeue()
	require.False(t, success)
	require.Equal(t, 0, value)
}

func TestRingBufferBasic(t *testing.T) {
	buffer := NewRingBuffer[int](100000)
	n := 4
	for i := 1; i <= n; i++ {
		buffer.Enqueue(i)
	}
	sum := 0
	for {
		if v, ok := buffer.Dequeue(); ok {
			sum += v
		} else {
			break
		}
	}
	require.Equal(t, (1+n)*n/2, sum)
}

func TestRingBufferConcurrent(t *testing.T) {
	queue := NewRingBuffer[int](100)
	for i := 0; i < 10; i++ {
		target, produceCnt, consumeCnt := queueReadWrite(queue, 100000,
			10, 10, false)
		t.Logf("Round #%d, target %d,p %d,c %d\n", i, target, produceCnt, consumeCnt)
		require.Equal(t, produceCnt, consumeCnt)
	}
}

func TestRingBufferCorrect(t *testing.T) {
	buf := NewRingBuffer[int](100)
	var pCnt, pSum, cCnt, cSum uint64
	pWait := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		pWait.Add(1)
		go func() {
			defer pWait.Done()
			for k := 0; k < 100000; k++ {
				suc := buf.Enqueue(k)
				if suc {
					atomic.AddUint64(&pCnt, 1)
					atomic.AddUint64(&pSum, uint64(k))
				}
			}
		}()
	}
	cWait := &sync.WaitGroup{}
	stopped := false
	for i := 0; i < 10; i++ {
		cWait.Add(1)
		go func() {
			defer cWait.Done()
			for {
				val, suc := buf.Dequeue()
				if suc {
					if val == -1 {
						time.Sleep(100 * time.Millisecond)
						stopped = true
						break
					}
					//t.Logf("c %d\n", val)
					atomic.AddUint64(&cCnt, 1)
					atomic.AddUint64(&cSum, uint64(val))
				} else {
					if stopped {
						break
					}
				}
			}
		}()
	}
	pWait.Wait()
	for !buf.Enqueue(-1) {
	}
	cWait.Wait()
	t.Logf("pCnt %d,pSum %d,cCnt %d,cSum %d\n", pCnt, pSum, cCnt, cSum)
	require.Equal(t, pCnt, cCnt)
	require.True(t, pSum >= cSum)
	require.Equal(t, pSum, cSum)
}

func TestRingBuffer_Dequeue_Loop(t *testing.T) {
	t.Skip()
	queue := NewRingBuffer[int](100)
	//for {
	//	res, ok := queue.Dequeue()
	//	if ok {
	//		t.Log(res)
	//	} else {
	//		runtime.Gosched()
	//	}
	//}
	// 空转会导致CPU 100%， 这里需要考虑元素为空时的等待策略，参考 http://www.enmalvi.com/2023/03/22/disruptor/
	w := &sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			for {
				res, ok := queue.Dequeue()
				if ok {
					t.Log(res)
				} else {
					// runtime.Gosched() // 让出cpu还是会跑满全部CPU
					time.Sleep(time.Nanosecond * 100)
				}
			}
		}()
	}
	w.Wait()
}

// BenchmarkRingBuffer10000
// go test -v -run ^$  -bench BenchmarkRingBuffer10000 -benchtime=5s -benchmem -memprofile memprofile.pprof -cpuprofile profile.pprof
// go tool pprof -http=":8081" memprofile.pprof
func BenchmarkRingBuffer10000(b *testing.B) {
	q := NewRingBuffer[int](capacity)
	for i := 0; i < b.N; i++ {
		queueReadWrite(q, 10000, 10, 5, false)
	}
}

// BenchmarkRingBufferAndMap
// go test -v -run ^$  -bench BenchmarkRingBufferAndMap -benchtime=10s -benchmem
func BenchmarkRingBufferAndMap(b *testing.B) {
	b.Run("ringBuffer", func(b *testing.B) {
		q := NewRingBuffer[int](capacity)
		for i := 0; i < b.N; i++ {
			q.Enqueue(1)
		}
	})
	b.Run("mapSeqWrite", func(b *testing.B) {
		m := make(map[int]struct{})
		for i := 0; i < b.N; i++ {
			m[i] = struct{}{}
		}
	})
	b.Run("mapEachWrite", func(b *testing.B) {
		m := make(map[int]struct{})
		idx := 0
		maxMask := 1<<8 - 1
		for i := 0; i < b.N; i++ {
			idx++
			m[idx&maxMask] = struct{}{}
		}
	})
	b.Run("mapNormalWrite", func(b *testing.B) {
		m := make(map[int]struct{})
		for i := 0; i < 300; i++ {
			m[i] = struct{}{}
		}
		for i := 0; i < b.N; i++ {
			m[150] = struct{}{}
		}
	})
}
