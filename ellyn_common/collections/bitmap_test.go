package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBitMap(t *testing.T) {
	m := NewBitMap(1)
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

func BenchmarkBitMap(b *testing.B) {
	m := NewBitMap(1024)
	for i := 0; i < b.N; i++ {
		m.Set(uint(i) & 1023)
	}
}
