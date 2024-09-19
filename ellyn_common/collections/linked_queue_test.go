package collections

import (
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLinkedQueueBasic(t *testing.T) {
	queue := NewLinkedQueue[int](100)
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	val, _ := queue.Dequeue()
	require.Equal(t, 1, val)
	val, _ = queue.Dequeue()
	require.Equal(t, 2, val)
	val, _ = queue.Dequeue()
	require.Equal(t, 3, val)
}

func queueReadWrite(queue Queue[int], num, producerCnt, consumerCnt int,
	block bool) (targetTotal int64, actualProduceCnt int64, actualConsumeCnt int64) {
	targetTotal = int64(num * producerCnt)
	producerWait := sync.WaitGroup{}
	producerWait.Add(producerCnt)
	for i := 0; i < producerCnt; i++ {
		go func() {
			defer producerWait.Done()
			for i := 0; i < num; i++ {
				suc := queue.Enqueue(1)
				if suc {
					atomic.AddInt64(&actualProduceCnt, 1)
				}
			}
		}()
	}
	consumerWait := sync.WaitGroup{}
	consumerWait.Add(consumerCnt)
	stoppped := false
	for i := 0; i < consumerCnt; i++ {
		go func() {
			defer consumerWait.Done()
			for {
				val, ok := queue.Dequeue()
				if ok {
					if val == 0 {
						//fmt.Printf("consumer[#%d] stopped\n", runtime.EllynGetGoid())
						if !block {
							stoppped = true
						}
						break
					}
					atomic.AddInt64(&actualConsumeCnt, 1)
				}
				if stoppped {
					break
				}
			}
		}()
	}

	producerWait.Wait()
	if block {
		for i := 0; i < consumerCnt; i++ {
			queue.Enqueue(0)
		}
	} else {
		for !queue.Enqueue(0) {
		}
	}
	consumerWait.Wait()
	return
}

func TestLinkedQueue(t *testing.T) {
	queue := NewLinkedQueue[int](100)
	for i := 0; i < 10; i++ {
		total, produceCnt, consumeCnt := queueReadWrite(queue, 1000, 10, 10, true)
		// 阻塞队列一定会写入成功、消费成功，不会丢元素
		require.Equal(t, total, produceCnt)
		require.Equal(t, produceCnt, consumeCnt)
	}
}
