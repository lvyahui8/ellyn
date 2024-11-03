package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
	"unsafe"
)

func TestNewBitMap(t *testing.T) {
	m := NewBitMap(1)
	require.Equal(t, 64, int(unsafe.Sizeof(*m)))
	require.Equal(t, 1, len(m.slots))
	m = NewBitMap(64)
	require.Equal(t, 1, len(m.slots))
	m = NewBitMap(65)
	require.Equal(t, 2, len(m.slots))
	m = NewBitMap(128)
	require.Equal(t, 2, len(m.slots))
	m = NewBitMap(129)
	require.Equal(t, 3, len(m.slots))
}

func TestBitMapBasic(t *testing.T) {
	m := NewBitMap(10)
	m.Set(1)
	m.Set(2)
	require.True(t, m.Get(1))
	require.True(t, m.Get(2))
	require.Equal(t, 2, m.Size())
	require.False(t, m.Get(3))
	m.Clear(1)
	require.False(t, m.Get(1))
	require.Equal(t, 1, m.Size())
}

func TestBitMap_FieldPadding(t *testing.T) {
	m := NewBitMap(1024)
	t.Logf("slots %v", unsafe.Sizeof(m.slots))
	t.Logf("size %v", unsafe.Sizeof(m.size))
	t.Logf("cap %v", unsafe.Sizeof(m.cap))
}

func BenchmarkBitMap(b *testing.B) {
	m := NewBitMap(1024)
	for i := 0; i < b.N; i++ {
		m.Set(uint(i) & 1023)
	}
}

// go test -v -run ^$  -bench 'BenchmarkBitMapAndArrayAndMap' -benchtime=5s -benchmem -cpuprofile profile.pprof
// go tool pprof -http=":8081" profile.pprof
func BenchmarkBitMapAndArrayAndMap(b *testing.B) {
	size := 100 * 10000
	bitMap := NewBitMap(uint(size))
	arr := make([]bool, size)
	m := make(map[int]struct{}, size)
	for i := 0; i < size; i++ {
		bitMap.Set(uint(i))
		arr[i] = false
		m[i] = struct{}{}
	}
	target := size / 2
	b.Run("bitMap_read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = bitMap.Get(uint(target))
		}
	})
	b.Run("array_read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = arr[target]
		}
	})
	b.Run("map_read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[target]
		}
	})

	b.Run("bitMap_write", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bitMap.Set(uint(target))
		}
	})
	b.Run("array_write", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arr[target] = true
		}
	})
	b.Run("map_write", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m[target] = struct{}{}
		}
	})
}
