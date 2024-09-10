package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntMapWrap_Get(t *testing.T) {
	m := map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	}
	w := NewIntMapWrap(m)
	_, exist := w.Get(0)
	require.False(t, exist)
	_, exist = w.Get(1)
	require.True(t, exist)
	_, exist = w.Get(2)
	require.True(t, exist)
	_, exist = w.Get(3)
	require.True(t, exist)
	_, exist = w.Get(4)
	require.True(t, exist)
	_, exist = w.Get(5)
	require.False(t, exist)
	_, exist = w.Get(6)
	require.False(t, exist)
}

/*

从结果可以看到快了一个数量级

goos: windows
goarch: amd64
pkg: github.com/lvyahui8/ellyn/ellyn_common/collections
cpu: AMD Ryzen 7 4800U with Radeon Graphics
BenchmarkIntMapWrap_Get
BenchmarkIntMapWrap_Get/rawMap
BenchmarkIntMapWrap_Get/rawMap-16                  10000            112451 ns/op

BenchmarkIntMapWrap_Get/wrapMap
BenchmarkIntMapWrap_Get/wrapMap-16                101760             10077 ns/op

PASS
*/
func BenchmarkIntMapWrap_Get(b *testing.B) {
	list := randomSeq(10000, 1000000)
	read := randomSeq(10000, 1000000)
	m := make(map[int]struct{})
	for _, val := range list {
		m[val] = struct{}{}
	}
	b.Run("rawMap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, idx := range read {
				_, _ = m[idx]
			}
		}
	})
	w := NewIntMapWrap(m)
	b.Run("wrapMap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, idx := range read {
				_, _ = w.Get(idx)
			}
		}
	})
}
