package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type poolVal int

func (poolVal) Recycle() {

}

func TestLRUCache_Basic(t *testing.T) {
	cache := NewLRUCache[int, poolVal](3)
	cache.Set(1, 1)
	res, ok := cache.Get(1)
	require.True(t, ok)
	require.Equal(t, poolVal(1), res)
	cache.Set(2, 2)
	res, ok = cache.Get(2)
	require.True(t, ok)

	require.Equal(t, poolVal(2), res)
	cache.Set(3, 3)
	cache.Set(4, 4)
	res, ok = cache.Get(1) // key 1 最先被淘汰
	require.False(t, ok)
	require.Equal(t, poolVal(4), cache.head.val)
	require.Equal(t, poolVal(2), cache.tail.val)
	cache.Remove(3)
	require.Equal(t, poolVal(4), cache.head.val)
	require.Equal(t, poolVal(2), cache.tail.val)
	require.Equal(t, poolVal(2), cache.head.next.val)
	require.Nil(t, cache.head.prev)
	require.Equal(t, poolVal(4), cache.tail.prev.val)
	require.Nil(t, cache.tail.next)
}
