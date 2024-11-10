package collections

import (
	"github.com/lvyahui8/ellyn/sdk/common/definitions"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"sort"
	"sync"
)

var _ mapApi[int, any] = (*ConcurrentMap[int, any])(nil)

type keyHasher[T comparable] func(t T) int

type mapApi[K comparable, V any] interface {
	Store(key K, val V)
	Load(key K) (V, bool)
	Delete(key K)
}

type ConcurrentMap[K comparable, V any] struct {
	segMask   int // 8
	_padding0 [56]byte
	segments  []*mapSegment[K, V] // 24字节，底层是一个slice struct，里面有array指针+两个int字段，3 * 8
	_padding1 [40]byte
	hasher    keyHasher[K] // 8
	_padding2 [56]byte
	size      uint64 // 8 字节
	_padding3 [56]byte
}

func NewNumberKeyConcurrentMap[K definitions.Number, V any](segSize int) *ConcurrentMap[K, V] {
	return NewConcurrentMap[K, V](segSize, func(t K) int {
		return int(t)
	})
}

func NewConcurrentMap[K comparable, V any](segSize int, hasher keyHasher[K]) *ConcurrentMap[K, V] {
	segSize = int(roundingToPowerOfTwo(uint64(segSize)))
	m := &ConcurrentMap[K, V]{
		segMask:  segSize - 1,
		segments: make([]*mapSegment[K, V], segSize),
		hasher:   hasher,
	}
	for i := 0; i < segSize; i++ {
		m.segments[i] = &mapSegment[K, V]{
			entries: make(map[K]V),
		}
	}
	return m
}

func (m *ConcurrentMap[K, V]) Store(key K, val V) {
	m.getSegment(key).Store(key, val)
}

func (m *ConcurrentMap[K, V]) Load(key K) (V, bool) {
	return m.getSegment(key).Load(key)
}

func (m *ConcurrentMap[K, V]) Delete(key K) {
	m.getSegment(key).Delete(key)
}

func (m *ConcurrentMap[K, V]) getSegment(key K) *mapSegment[K, V] {
	return m.segments[m.hasher(key)&m.segMask]
}

func (m *ConcurrentMap[K, V]) Size() int {
	size := 0
	for _, s := range m.segments {
		size += s.Size()
	}
	return size
}

func (m *ConcurrentMap[K, V]) Values() (res []V) {
	for _, seg := range m.segments {
		seg.RLock()
		res = append(res, utils.GetMapValues(seg.entries)...)
		seg.RUnlock()
	}
	return
}

func (m *ConcurrentMap[K, V]) SortedValues(compare func(a, b V) bool) (res []V) {
	for _, seg := range m.segments {
		seg.RLock()
		res = append(res, utils.GetMapValues(seg.entries)...)
		seg.RUnlock()
	}
	sort.Slice(res, func(i, j int) bool {
		return compare(res[i], res[j])
	})
	return
}

type mapSegment[K comparable, V any] struct {
	sync.RWMutex // 24
	_padding0    [40]byte
	entries      map[K]V // 8
	_padding1    [56]byte
	size         int      // 8
	_padding2    [56]byte //
}

func (s *mapSegment[K, V]) Store(key K, val V) {
	s.Lock()
	s.entries[key] = val
	s.size = len(s.entries)
	s.Unlock()
}

func (s *mapSegment[K, V]) Load(key K) (res V, ok bool) {
	s.RLock()
	res, ok = s.entries[key]
	s.RUnlock()
	return
}

func (s *mapSegment[K, V]) Delete(key K) {
	s.Lock()
	delete(s.entries, key)
	s.size = len(s.entries)
	s.Unlock()
}

func (s *mapSegment[K, V]) Size() int {
	return s.size
}
