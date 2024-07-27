package collections

import (
	"math/rand"
	"testing"
	"unsafe"
)

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
