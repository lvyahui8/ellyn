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

// UnsafeCompressedStack 非并发安全的Stack
type UnsafeCompressedStack struct {
	elements *list.List
	count    int
}

type stackElement struct {
	val   interface{}
	max   int
	count int
}

func NewStack() *UnsafeCompressedStack {
	return &UnsafeCompressedStack{
		elements: list.New(),
	}
}

func (s *UnsafeCompressedStack) Push(val interface{}) {
	s.count++
	back := s.elements.Back()
	if back != nil {
		ele := back.Value.(*stackElement)
		if ele.val == val {
			ele.max++
			ele.count++
			return
		}
	}
	s.elements.PushBack(&stackElement{val: val, max: 1, count: 1})
}

func (s *UnsafeCompressedStack) Pop() interface{} {
	e := s.elements.Back()
	if e != nil {
		ele := e.Value.(*stackElement)
		if ele.count == 1 {
			s.elements.Remove(e)
		} else {
			ele.count--
		}
		s.count--
		return ele.val
	} else {
		return nil
	}
}

func (s *UnsafeCompressedStack) Top() interface{} {
	e := s.elements.Back()
	if e != nil {
		return e.Value.(*stackElement).val
	}
	return nil
}

func (s *UnsafeCompressedStack) Size() int {
	return s.count
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
