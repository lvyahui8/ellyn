package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLRUCache_Basic(t *testing.T) {
	cache := NewLRUCache[int](3)
	cache.Set(1, "a")
	res, ok := cache.Get(1)
	require.True(t, ok)
	val := res.(string)
	require.Equal(t, "a", val)
	cache.Set(2, "b")
	res, ok = cache.Get(2)
	require.True(t, ok)
	val = res.(string)
	require.Equal(t, "b", val)
	cache.Set(3, "c")
	cache.Set(4, "d")
	res, ok = cache.Get(1) // key 1 最先被淘汰
	require.False(t, ok)
	require.Equal(t, "d", cache.head.val)
	require.Equal(t, "b", cache.tail.val)
	cache.Remove(3)
	require.Equal(t, "d", cache.head.val)
	require.Equal(t, "b", cache.tail.val)
	require.Equal(t, "b", cache.head.next.val)
	require.Nil(t, cache.head.prev)
	require.Equal(t, "d", cache.tail.prev.val)
	require.Nil(t, cache.tail.next)
}
