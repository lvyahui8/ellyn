package collections

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
	"unsafe"
)

// go test -v -run ^$  -bench BenchmarkIndex -benchtime=5s

const listLen = 1 << 20

// BenchmarkMapAndArray 数据量大时，map的性能与array/slice的性能差距甚远，基本不能算o(1)的操作
func BenchmarkMapAndArray(b *testing.B) {
	mask := listLen - 1
	m := make(map[int]int, listLen)
	for i := 0; i < listLen; i++ {
		m[i] = i
	}
	s := make([]int, listLen)
	for i := 0; i < listLen; i++ {
		s[i] = i
	}
	var a [listLen]int
	for i := 0; i < listLen; i++ {
		a[i] = i
	}
	seq := randomSeq(listLen, listLen)
	b.Run("map", func(b *testing.B) {
		idx := 0
		for i := 0; i < b.N; i++ {
			idx = (idx + 1) & mask
			_ = m[seq[idx]]
		}
	})
	b.Run("slice", func(b *testing.B) {
		idx := 0
		for i := 0; i < b.N; i++ {
			idx = (idx + 1) & mask
			_ = s[seq[idx]]
		}
	})
	b.Run("array", func(b *testing.B) {
		idx := 0
		for i := 0; i < b.N; i++ {
			idx = (idx + 1) & mask
			_ = a[seq[idx]]
		}
	})
}

func TestInterfaceSize(t *testing.T) {
	var a any
	a = float64(10)
	var b any
	b = float32(10)
	t.Logf("any size=%d\n", unsafe.Sizeof(a))
	t.Logf("any size=%d\n", unsafe.Sizeof(b))
	var c *any
	c = &a
	t.Logf("*any size=%d\n", unsafe.Sizeof(c))
}

func shuffle(nums []int) (res []int) {
	res = make([]int, len(nums))
	copy(res, nums)
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(nums), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return
}

// 从0-maxVal之间随机挑选cnt个数字，并打乱顺序
func randomSeq(cnt, maxVal int) (res []int) {
	if cnt > maxVal {
		panic("cnt > maxVal")
	}
	raw := make([]int, maxVal)
	for i := 0; i < maxVal; i++ {
		raw[i] = i
	}
	raw = shuffle(raw)
	res = raw[0:cnt]
	return
}

// BenchmarkDeferOrDirectCall在明确不会panic的情况下，直接调用比defer调用性能更好
func BenchmarkDeferOrDirectCall(b *testing.B) {
	b.Run("defer", func(b *testing.B) {
		lock := &sync.RWMutex{}
		for i := 0; i < b.N; i++ {
			func() {
				lock.RLock()
				defer lock.RUnlock()
			}()
		}
	})
	b.Run("direct", func(b *testing.B) {
		lock := &sync.RWMutex{}
		for i := 0; i < b.N; i++ {
			func() {
				lock.RLock()
				lock.RUnlock()
			}()
		}
	})
}

func BenchmarkMapInter(b *testing.B) {
	type user struct {
		Id   int
		Name string
	}
	structMap := make(map[int]user)
	ptrMap := make(map[int]*user)
	emptyValMap := make(map[int]struct{})
	total := 100
	for i := 0; i < total; i++ {
		name := strconv.Itoa(i)
		structMap[i] = user{i, name}
		ptrMap[i] = &user{i, name}
		emptyValMap[i] = struct{}{}
	}

	b.Run("structMapNoVal", func(b *testing.B) {
		sum := 0
		for i := 0; i < b.N; i++ {
			for k := range structMap {
				sum += k
			}
		}
	})
	b.Run("structMapWithVal", func(b *testing.B) {
		sum := 0
		for i := 0; i < b.N; i++ {
			for k, v := range structMap {
				sum += k
				if v.Id == 0 {
				}
			}
		}
	})
	b.Run("ptrMapNoVal", func(b *testing.B) {
		sum := 0
		for i := 0; i < b.N; i++ {
			for k := range ptrMap {
				sum += k
			}
		}
	})
	b.Run("ptrMapWithVal", func(b *testing.B) {
		sum := 0
		for i := 0; i < b.N; i++ {
			for k, v := range ptrMap {
				sum += k
				if v.Id == 0 {
				}
			}
		}
	})

	b.Run("emptyMapNoVal", func(b *testing.B) {
		sum := 0
		for i := 0; i < b.N; i++ {
			for k := range emptyValMap {
				sum += k
			}
		}
	})
	b.Run("emptyMapWithVal", func(b *testing.B) {
		sum := 0
		for i := 0; i < b.N; i++ {
			for k, v := range emptyValMap {
				sum += k
				if v == struct{}{} {
				}
			}
		}
	})
}
