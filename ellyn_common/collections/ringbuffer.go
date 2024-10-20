package collections

import (
	"sync/atomic"
)

type RingBuffer[T any] struct {
	// dequeuePos 指向下一个可消费点位，一直累加然后对capacity取模（&mask）取值
	dequeuePos uint64
	_padding0  [56]byte
	// dequeuePos 指向下一个可写入点位，一直累加然后对capacity取模（&mask）取值
	enqueuePos uint64
	_padding1  [56]byte
	// mask capacity - 1
	mask      uint64
	_padding2 [56]byte
	elements  []*node[T]
	_padding3 [56]byte
}

type node[T any] struct {
	// seq一直累加，用来标记生产、消费状态
	// 当seq=enqueuePos时，表示位置为空可写入
	// 当seq=dequeuePos+1时，表示当前位置可以消费
	// 将seq设置为dequeuePos+capacity时，说明消费掉这个元素，将他设置到下一个可写入的窗口（enqueuePos会循环追上这个值）
	seq       uint64
	_padding0 [56]byte
	// 具体元素值
	value     T // 64位处理器上 size 16字节
	_padding1 [48]byte
}

func NewRingBuffer[T any](capacity uint64) *RingBuffer[T] {
	capacity = roundingToPowerOfTwo(capacity)
	nodes := make([]*node[T], capacity)
	for i := uint64(0); i < capacity; i++ {
		nodes[i] = &node[T]{seq: i}
	}

	return &RingBuffer[T]{
		dequeuePos: uint64(0),
		enqueuePos: uint64(0),
		mask:       capacity - 1,
		elements:   nodes,
	}
}

// Enqueue 非阻塞式写入，当缓冲区满时，返回失败
func (r *RingBuffer[T]) Enqueue(value T) (success bool) {
	pos := atomic.LoadUint64(&r.enqueuePos)
	var element *node[T]
	for {
		element = r.elements[pos&r.mask]
		seq := atomic.LoadUint64(&element.seq)
		diff := int64(seq) - int64(pos)
		if diff == 0 {
			// 可以尝试写入
			if atomic.CompareAndSwapUint64(&r.enqueuePos, pos, pos+1) {
				// 写入成功
				break
			}
		} else if diff < 0 {
			// 缓冲区满，写入失败
			return false
		} else {
			// pos由于并发已经过期，需要重新读取
			pos = atomic.LoadUint64(&r.enqueuePos)
		}
	}
	element.value = value
	// 将seq设置为待消费标识
	atomic.StoreUint64(&element.seq, pos+1)
	return true
}

// Dequeue 非阻塞式出队，当没有元素可以取时（缓冲区空）返回失败
func (r *RingBuffer[T]) Dequeue() (value T, success bool) {
	var element *node[T]
	pos := atomic.LoadUint64(&r.dequeuePos)
	for {
		element = r.elements[pos&r.mask]
		seq := atomic.LoadUint64(&element.seq)
		diff := int64(seq) - int64(pos+1)
		if diff == 0 {
			// 可以尝试读取
			if atomic.CompareAndSwapUint64(&r.dequeuePos, pos, pos+1) {
				break
			}
		} else if diff < 0 {
			// 缓冲区为空，没有元素可以消费
			return
		} else {
			// pos由于并发已经过期，需要重新读取
			pos = atomic.LoadUint64(&r.dequeuePos)
		}
	}
	// Dequeue中将seq更新一个容量跨度，相当于设置本轮循环的数据已经消费掉了
	atomic.StoreUint64(&element.seq, pos+r.mask+1)
	return element.value, true
}
