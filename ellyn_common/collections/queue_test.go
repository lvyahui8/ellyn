package collections

import "testing"

func BenchmarkQueue(b *testing.B) {
	b.Run("linkedQueue", func(b *testing.B) {
		queue := NewLinkedQueue[int](0)
		for i := 0; i < b.N; i++ {
			queue.Enqueue(1)
			queue.Dequeue()
		}
	})
	b.Run("ringBuffer", func(b *testing.B) {
		queue := NewRingBuffer(1028)
		for i := 0; i < b.N; i++ {
			queue.Enqueue(1)
			queue.Dequeue()
		}
	})
}
