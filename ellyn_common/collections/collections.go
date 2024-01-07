package collections

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type keyHasher func(key interface{}) int

type mapApi interface {
	Put(key, val interface{})
	Get(key interface{}) interface{}
	Del(key interface{})
}

type concurrentMap struct {
	segMask  int
	segments []*mapSegment
	hasher   keyHasher
	size     uint64
}

func NewConcurrentMap(segSize int, hasher keyHasher) *concurrentMap {
	segSize = int(roundingToPowerOfTwo(uint64(segSize)))
	m := &concurrentMap{
		segMask:  segSize - 1,
		segments: make([]*mapSegment, segSize),
		hasher:   hasher,
	}
	for i := 0; i < segSize; i++ {
		m.segments[i] = &mapSegment{
			entries: make(map[interface{}]interface{}),
		}
	}
	return m
}

func (m *concurrentMap) Put(key, val interface{}) {
	m.getSegment(key).Put(key, val)
}

func (m *concurrentMap) Get(key interface{}) interface{} {
	return m.getSegment(key).Get(key)
}

func (m *concurrentMap) Del(key interface{}) {
	m.getSegment(key).Del(key)
}

func (m *concurrentMap) getSegment(key interface{}) *mapSegment {
	return m.segments[m.hasher(key)&m.segMask]
}

func (m *concurrentMap) Size() int {
	size := 0
	for _, s := range m.segments {
		size += s.Size()
	}
	return size
}

type mapSegment struct {
	sync.RWMutex
	entries map[interface{}]interface{}
	size    int
}

func (s *mapSegment) Put(key, val interface{}) {
	s.Lock()
	defer s.Unlock()
	s.entries[key] = val
	s.size = len(s.entries)
}

func (s *mapSegment) Get(key interface{}) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.entries[key]
}

func (s *mapSegment) Del(key interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.entries, key)
	s.size = len(s.entries)
}

func (s *mapSegment) Size() int {
	return s.size
}

// unsafeStack 非并发安全的Stack
type unsafeStack struct {
	elements *list.List
}

func NewStack() *unsafeStack {
	return &unsafeStack{
		elements: list.New(),
	}
}

func (s *unsafeStack) Push(val interface{}) {
	s.elements.PushBack(val)
}

func (s *unsafeStack) Pop() interface{} {
	e := s.elements.Back()
	if e != nil {
		s.elements.Remove(e)
	}
	return e.Value
}

func (s *unsafeStack) Top() interface{} {
	e := s.elements.Back()
	if e != nil {
		return e.Value
	}
	return nil
}

func (s *unsafeStack) Size() int {
	return s.elements.Len()
}

// https://www.lenshood.dev/2021/04/19/lock-free-ring-buffer/

func roundingToPowerOfTwo(size uint64) uint64 {
	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size |= size >> 32
	size++
	return size
}

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
