package collections

import (
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"testing"
	"unsafe"
)

func testReadWrite(b *testing.B, name string, m mapApi[int, any], insertSeq, readSeq, deleteSeq []int) {
	// 测试纯粹写入的并发性能
	routineSize := runtime.NumCPU()
	handleSize := len(insertSeq) / routineSize
	b.Run(name+"_put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := sync.WaitGroup{}
			w.Add(routineSize)
			for j := 0; j < routineSize; j++ {
				go func(offset int) {
					defer w.Done()
					for i := offset; i < handleSize+offset; i++ {
						m.Store(insertSeq[i], struct{}{})
					}
				}(handleSize * j)
			}
			w.Wait()
		}
	})
	// 测试纯粹读取的并发性能
	b.Run(name+"_read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := sync.WaitGroup{}
			w.Add(routineSize)
			for j := 0; j < routineSize; j++ {
				go func(offset int) {
					defer w.Done()
					for i := offset; i < handleSize+offset; i++ {
						_, _ = m.Load(readSeq[i])
					}
				}(handleSize * j)
			}
			w.Wait()
		}
	})
	// 测试纯粹删除的并发性能
	b.Run(name+"_delete", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := sync.WaitGroup{}
			w.Add(routineSize)
			for j := 0; j < routineSize; j++ {
				go func(offset int) {
					defer w.Done()
					for i := offset; i < handleSize+offset; i++ {
						m.Delete(deleteSeq[i])
					}
				}(handleSize * j)
			}
			w.Wait()
		}
	})
	// 测试同时写入、读取的并发性能
	b.Run(name+"_R&W", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := sync.WaitGroup{}
			w.Add(routineSize * 2)
			for j := 0; j < routineSize; j++ {
				go func(offset int) {
					defer w.Done()
					for i := offset; i < handleSize+offset; i++ {
						m.Store(insertSeq[i], struct{}{})
					}
				}(handleSize * j)
			}
			for j := 0; j < routineSize; j++ {
				go func(offset int) {
					defer w.Done()
					for i := offset; i < handleSize+offset; i++ {
						_, _ = m.Load(readSeq[i])
					}
				}(handleSize * j)
			}
			w.Wait()
		}
	})
}

const maxVal = 1000 * 10000
const cnt = 100 * 10000

type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func (s *SyncMap[K, V]) Store(key K, val V) {
	s.m.Store(key, val)
}

func (s *SyncMap[K, V]) Load(key K) (v V, exist bool) {
	value, ok := s.m.Load(key)
	if !ok {
		exist = false
		return
	}
	return value.(V), true
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.m.Delete(key)
}

func BenchmarkMap(b *testing.B) {
	insertSeq := randomSeq(cnt, maxVal)
	readSeq := shuffle(insertSeq)
	deleteSeq := shuffle(insertSeq)
	testReadWrite(b, "ConcurrentMap", NewNumberKeyConcurrentMap[int, any](2048), insertSeq, readSeq, deleteSeq)
	testReadWrite(b, "syncMap", &SyncMap[int, any]{}, insertSeq, readSeq, deleteSeq)
}

func TestMapPadding(t *testing.T) {
	m := NewNumberKeyConcurrentMap[int, struct{}](2048)
	require.Equal(t, 256, int(unsafe.Sizeof(*m)))
	t.Log(unsafe.Sizeof(m.size))
	t.Log(unsafe.Sizeof(m.segMask))
	t.Log(unsafe.Sizeof(m.hasher))
	t.Log(unsafe.Sizeof(m.segments))
	t.Log("=======")
	ms := mapSegment[int, struct{}]{}
	require.Equal(t, 64, int(unsafe.Sizeof(ms)))
	t.Log(unsafe.Sizeof(ms.RWMutex))
	t.Log(unsafe.Sizeof(ms.entries))
	t.Log(unsafe.Sizeof(ms.size))
}

func TestConcurrentMap(t *testing.T) {
	m := NewNumberKeyConcurrentMap[int, string](10)
	str := "xxx"
	m.Store(1, str)
	s, ok := m.Load(1)
	require.True(t, ok)
	require.Equal(t, str, s)
	m.Delete(1)
	_, ok = m.Load(1)
	require.False(t, ok)
}
