package collections

import (
	"fmt"
	"testing"
)

// BenchmarkQueue go test -v -run ^$  -bench BenchmarkQueue -benchtime=10s -benchmem
func BenchmarkQueue(b *testing.B) {
	// 读写元素个数
	targetCnt := 100000
	routineCntList := []int{1, 10, 100, 500} // 生产者消费者数量
	for _, cnt := range routineCntList {
		k := cnt
		b.Run(fmt.Sprintf("%d_RingBuffer", k), func(b *testing.B) {
			q := NewRingBuffer[int](capacity)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				queueReadWrite(q, targetCnt, k, k, false)
			}
		})
		b.Run(fmt.Sprintf("%d_channelBuffer", k), func(b *testing.B) {
			q := newChannelQueue[int](capacity)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				queueReadWrite(q, targetCnt, k, k, true)
			}
			q.Close()
		})
		b.Run(fmt.Sprintf("%d_linkedBlockQueue", k), func(b *testing.B) {
			q := NewLinkedQueue[int](capacity)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				queueReadWrite(q, targetCnt, k, k, true)
			}
		})
	}
}
