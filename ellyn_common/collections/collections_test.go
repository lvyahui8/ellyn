package collections

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

var buffer = NewRingBuffer(2000)

func queueReadWrite(cnt int, produce func() interface{}, consume func(val interface{})) {
	pNum := runtime.NumCPU()
	cNum := runtime.NumCPU()
	group := sync.WaitGroup{}
	group.Add(pNum + cNum)

	for i := 0; i < pNum; i++ {
		go func() {
			defer group.Done()
			for k := 0; k < cnt; k++ {
				var suc bool
				if produce != nil {
					suc = buffer.Offer(produce())
				} else {
					suc = buffer.Offer(1)
				}
				if !suc {
					fmt.Printf("offer failed.\n")
				}
			}
		}()
	}

	for i := 0; i < cNum; i++ {
		go func() {
			defer group.Done()
			for {
				val, suc := buffer.Poll()
				if consume != nil && suc {
					consume(val)
				}
			}
		}()
	}
	group.Wait()
	return
}

func TestRingBuffer(t *testing.T) {
	var sum uint64 = 0
	queueReadWrite(100,
		func() interface{} {
			time.Sleep(10 * time.Millisecond)
			return 1
		},
		func(val interface{}) {
			t.Logf("%d\n", val.(int))
			atomic.AddUint64(&sum, uint64(val.(int)))
		})
	t.Logf("sum %d\n", sum)
}

// BenchmarkRingBuffer go test -bench BenchmarkRingBuffer -benchtime=5s
func BenchmarkRingBuffer(b *testing.B) {
	cntList := []int{1, 10, 100, 1000, 10000}
	for _, cnt := range cntList {
		k := cnt
		b.Run(fmt.Sprintf("readWrite_%d", k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				queueReadWrite(k, nil, nil)
			}
		})
	}
}

func TestInterfaceSize(t *testing.T) {
	var a interface{}
	a = float64(10)
	var b interface{}
	b = float32(10)
	t.Logf("interface{} size=%d\n", unsafe.Sizeof(a))
	t.Logf("interface{} size=%d\n", unsafe.Sizeof(b))
	var c *interface{}
	c = &a
	t.Logf("*interface{} size=%d\n", unsafe.Sizeof(c))
}

// go test -v -run ^$  -bench BenchmarkIndex -benchtime=5s
const listLen = 10000

func BenchmarkIndexMap(b *testing.B) {
	m := make(map[int]int, listLen)
	for i := 0; i < listLen; i++ {
		m[i] = i
	}
	for i := 0; i < b.N; i++ {
		_ = m[rand.Intn(listLen)]
	}
}

func BenchmarkIndexSlice(b *testing.B) {
	s := make([]int, listLen)
	for i := 0; i < listLen; i++ {
		s[i] = i
	}
	for i := 0; i < b.N; i++ {
		_ = s[rand.Intn(listLen)]
	}
}

func BenchmarkIndexArray(b *testing.B) {
	var a [listLen]int
	for i := 0; i < listLen; i++ {
		a[i] = i
	}
	for i := 0; i < b.N; i++ {
		_ = a[rand.Intn(listLen)]
	}
}
