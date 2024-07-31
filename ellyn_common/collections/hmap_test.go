package collections

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"
)

const maxVal = 10000 * 10000
const cnt = 1000 * 10000

func shuffle(nums []int) (res []int) {
	res = make([]int, len(nums))
	for i, v := range nums {
		res[i] = v
	}
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(nums), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return
}

// 生成100w个小于1000w不重复的数字
func randomSeq() (res []int) {
	raw := make([]int, maxVal)
	for i := 0; i < maxVal; i++ {
		raw[i] = i
	}
	res = shuffle(raw[0:cnt])
	return
}

func testReadWrite(b *testing.B, name string, m mapApi, insertSeq, readSeq, deleteSeq []int) {
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

func BenchmarkMap(b *testing.B) {
	insertSeq := randomSeq()
	readSeq := shuffle(insertSeq)
	deleteSeq := shuffle(insertSeq)
	testReadWrite(b, "concurrentMap", NewConcurrentMap(2048, func(key any) int {
		return key.(int)
	}), insertSeq, readSeq, deleteSeq)
	sMap := &sync.Map{}
	testReadWrite(b, "syncMap", sMap, insertSeq, readSeq, deleteSeq)
}

func TestMapPadding(t *testing.T) {
	m := NewConcurrentMap(2048, func(key any) int {
		return key.(int)
	})
	require.Equal(t, 64, int(unsafe.Sizeof(*m)))
	t.Log(unsafe.Sizeof(m.size))
	t.Log(unsafe.Sizeof(m.segMask))
	t.Log(unsafe.Sizeof(m.hasher))
	t.Log(unsafe.Sizeof(m.segments))
	t.Log("=======")
	ms := mapSegment{}
	require.Equal(t, 64, int(unsafe.Sizeof(ms)))
	t.Log(unsafe.Sizeof(ms.RWMutex))
	t.Log(unsafe.Sizeof(ms.entries))
	t.Log(unsafe.Sizeof(ms.size))
}
