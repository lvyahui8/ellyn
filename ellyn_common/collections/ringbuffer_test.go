package collections

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"testing"
)

//var buffer =

const capacity = 100000

type Queue interface {
	Enqueue(value interface{}) (success bool)
	Dequeue() (value interface{}, success bool)
}

type channelQueue struct {
	q chan interface{}
}

func newChannelQueue() *channelQueue {
	return &channelQueue{
		q: make(chan interface{}, capacity),
	}
}

func (c channelQueue) Enqueue(value interface{}) (success bool) {
	c.q <- value
	return true
}

func (c channelQueue) Dequeue() (value interface{}, success bool) {
	select {
	case value, success = <-c.q:
	default:
	}
	return
}

func (c channelQueue) Close() {
	close(c.q)
}

func queueReadWrite(queue Queue, cnt int) {
	pNum := runtime.NumCPU()
	cNum := runtime.NumCPU()
	group := sync.WaitGroup{}
	group.Add(pNum + cNum)

	for i := 0; i < pNum; i++ {
		go func() {
			defer group.Done()
			for k := 0; k < cnt; k++ {
				_ = queue.Enqueue(1)
			}
		}()
	}

	for i := 0; i < cNum; i++ {
		go func() {
			defer group.Done()
			for {
				_, suc := queue.Dequeue()
				if !suc {
					return
				}
			}
		}()
	}
	group.Wait()
	return
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

func TestRingBuffer(t *testing.T) {
	queueReadWrite(NewRingBuffer(2000), 100)
}

// BenchmarkRingBuffer go test -v -run ^$  -bench BenchmarkRingBuffer -benchtime=5s -benchmem
func BenchmarkRingBuffer(b *testing.B) {
	// 读写元素个数
	cntList := []int{1, 10, 100, 1000, 10000}
	for _, cnt := range cntList {
		k := cnt
		b.Run(fmt.Sprintf("ringBuffer readWrite_%d", k), func(b *testing.B) {
			q := NewRingBuffer(capacity)
			for i := 0; i < b.N; i++ {
				queueReadWrite(q, k)
			}
		})
		b.Run(fmt.Sprintf("channelBuffer readWrite_%d", k), func(b *testing.B) {
			q := newChannelQueue()
			for i := 0; i < b.N; i++ {
				queueReadWrite(q, k)
			}
			q.Close()
		})
	}
}
