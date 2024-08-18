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

func TestLinkedQueue(t *testing.T) {
	queue := NewLinkedQueue[int](100)
	num := 100
	producerCnt := 10
	total := num * producerCnt
	producerWait := sync.WaitGroup{}
	producerWait.Add(producerCnt)
	for i := 0; i < producerCnt; i++ {
		go func() {
			defer producerWait.Done()
			for i := 0; i < num; i++ {
				_ = queue.Enqueue(1)
			}
		}()
	}
	consumerCnt := 5
	var sum int64 = 0
	consumerWait := sync.WaitGroup{}
	consumerWait.Add(consumerCnt)
	for i := 0; i < consumerCnt; i++ {
		go func() {
			defer consumerWait.Done()
			for {
				val, _ := queue.Dequeue()
				if val == 0 {
					break
				}
				atomic.AddInt64(&sum, int64(val))
			}
		}()
	}
	producerWait.Wait()
	for i := 0; i < consumerCnt; i++ {
		queue.Enqueue(0)
	}
	consumerWait.Wait()
	require.Equal(t, int64(total), sum)
}
