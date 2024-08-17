package collections

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

//var buffer =

const capacity = 100000

type channelQueue struct {
	q chan any
}

func newChannelQueue() *channelQueue {
	return &channelQueue{
		q: make(chan any, capacity),
	}
}

func (c channelQueue) Enqueue(value any) (success bool) {
	c.q <- value
	return true
}

func (c channelQueue) Dequeue() (value any, success bool) {
	select {
	case value, success = <-c.q:
	default:
	}
	return
}

func (c channelQueue) Close() {
	close(c.q)
}

func queueReadWrite(queue Queue, cnt int) bool {
	pNum := runtime.NumCPU()
	cNum := runtime.NumCPU()
	group := sync.WaitGroup{}
	group.Add(pNum + cNum)
	produceGroup := sync.WaitGroup{}
	produceGroup.Add(pNum)
	var produceCnt uint64 = 0
	for i := 0; i < pNum; i++ {
		go func() {
			defer group.Done()
			defer produceGroup.Done()
			for k := 0; k < cnt; k++ {
				suc := queue.Enqueue(1)
				if suc {
					atomic.AddUint64(&produceCnt, 1)
				}
			}
		}()
	}
	var consumeCnt uint64 = 0
	stop := false
	for i := 0; i < cNum; i++ {
		go func() {
			defer group.Done()
			for {
				val, suc := queue.Dequeue()
				if suc {
					if val == 0 {
						stop = true
					} else {
						// 消费有效值
						atomic.AddUint64(&consumeCnt, 1)
					}
				}
				if stop {
					break
				}
			}
		}()
	}
	produceGroup.Wait()
	queue.Enqueue(0) // 队列最后一个元素标识
	group.Wait()
	return produceCnt == consumeCnt
}

func TestMaxIdx(t *testing.T) {
	s := math.MaxInt64 / int64(100*10000)
	year := s / int64(60*60*24*365)
	t.Log(year)
}

func TestRingBufferBasic(t *testing.T) {
	buffer := NewRingBuffer(100000)
	n := 4
	for i := 1; i <= n; i++ {
		buffer.Enqueue(i)
	}
	sum := 0
	for {
		if v, ok := buffer.Dequeue(); ok {
			sum += v.(int)
		} else {
			break
		}
	}
	require.Equal(t, (1+n)*n/2, sum)
}

func TestRingBufferConcurrent(t *testing.T) {
	require.True(t, queueReadWrite(NewRingBuffer(100), 10000))
}

// BenchmarkRingBuffer10000 go test -v -run ^$  -bench BenchmarkRingBuffer10000 -benchtime=5s -benchmem -memprofile memprofile.pprof -cpuprofile profile.pprof
// $ go tool pprof -http=":8081" memprofile.pprof
func BenchmarkRingBuffer10000(b *testing.B) {
	q := NewRingBuffer(capacity)
	for i := 0; i < b.N; i++ {
		queueReadWrite(q, 10000)
	}
}

// BenchmarkRingBuffer go test -v -run ^$  -bench BenchmarkRingBuffer -benchtime=5s -benchmem
func BenchmarkRingBuffer(b *testing.B) {
	// 读写元素个数
	cntList := []int{1, 10, 100, 1000, 10000} // 每个线程读写次数
	for _, cnt := range cntList {
		k := cnt
		b.Run(fmt.Sprintf("RingBuffer readWrite_%d", k), func(b *testing.B) {
			q := NewRingBuffer(capacity)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				queueReadWrite(q, k)
			}
		})
		b.Run(fmt.Sprintf("channelBuffer readWrite_%d", k), func(b *testing.B) {
			q := newChannelQueue()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				queueReadWrite(q, k)
			}
			q.Close()
		})
	}
}
