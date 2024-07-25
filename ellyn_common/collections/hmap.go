package collections

import "sync"

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
