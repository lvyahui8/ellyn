package collections

import "sync"

type keyHasher func(key interface{}) int

type mapApi interface {
	Store(key, val interface{})
	Load(key interface{}) (interface{}, bool)
	Delete(key interface{})
}

type concurrentMap struct {
	segMask  int           // 8
	segments []*mapSegment // 24字节，底层是一个slice struct，里面有array指针+两个int字段，3 * 8
	hasher   keyHasher     // 8
	size     uint64        // 8 字节
	_padding [16]byte
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

func (m *concurrentMap) Store(key, val interface{}) {
	m.getSegment(key).Store(key, val)
}

func (m *concurrentMap) Load(key interface{}) (interface{}, bool) {
	return m.getSegment(key).Load(key)
}

func (m *concurrentMap) Delete(key interface{}) {
	m.getSegment(key).Delete(key)
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
	sync.RWMutex                             // 24
	entries      map[interface{}]interface{} // 8
	size         int                         // 8
	_padding     [24]byte                    //
}

func (s *mapSegment) Store(key, val interface{}) {
	s.Lock()
	defer s.Unlock()
	s.entries[key] = val
	s.size = len(s.entries)
}

func (s *mapSegment) Load(key interface{}) (res interface{}, ok bool) {
	s.RLock()
	defer s.RUnlock()
	res, ok = s.entries[key]
	return
}

func (s *mapSegment) Delete(key interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.entries, key)
	s.size = len(s.entries)
}

func (s *mapSegment) Size() int {
	return s.size
}
