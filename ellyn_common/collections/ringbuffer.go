package collections

import "sync/atomic"

type ringBuffer struct {
	head      uint64
	_padding0 [56]byte
	tail      uint64
	_padding1 [56]byte
	mask      uint64
	_padding2 [56]byte
	element   []*node
}

type node struct {
	step     uint64
	value    interface{} // 64位处理器上 size 16字节  Interface consumes 2 words of memory: 1 word for runtime type and 1 word for data pointer.
	_padding [40]byte
}

func NewRingBuffer(capacity uint64) *ringBuffer {
	capacity = roundingToPowerOfTwo(capacity)
	nodes := make([]*node, capacity)
	for i := uint64(0); i < capacity; i++ {
		nodes[i] = &node{step: i}
	}

	return &ringBuffer{
		head:    uint64(0),
		tail:    uint64(0),
		mask:    capacity - 1,
		element: nodes,
	}
}

// Offer a value pointer.
func (r *ringBuffer) Offer(value interface{}) (success bool) {
	oldTail := atomic.LoadUint64(&r.tail)
	tailNode := r.element[oldTail&r.mask]
	oldStep := atomic.LoadUint64(&tailNode.step)
	// not published yet
	if oldStep != oldTail {
		return false
	}

	if !atomic.CompareAndSwapUint64(&r.tail, oldTail, oldTail+1) {
		return false
	}

	tailNode.value = value
	atomic.StoreUint64(&tailNode.step, tailNode.step+1)
	return true
}

// Poll head value pointer.
func (r *ringBuffer) Poll() (value interface{}, success bool) {
	oldHead := atomic.LoadUint64(&r.head)
	headNode := r.element[oldHead&r.mask]
	oldStep := atomic.LoadUint64(&headNode.step)
	// not published yet
	if oldStep != oldHead+1 {
		return
	}

	if !atomic.CompareAndSwapUint64(&r.head, oldHead, oldHead+1) {
		return
	}

	value = headNode.value
	atomic.StoreUint64(&headNode.step, oldStep+r.mask)
	return value, true
}
