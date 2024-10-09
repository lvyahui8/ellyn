package collections

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestLinkedList(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)
	list.Add(2)
	require.Equal(t, 2, list.Size())
	list.Remove(0)
	require.Equal(t, 1, list.Size())
	require.True(t, reflect.DeepEqual(list.Values(), []int{2}))
	require.False(t, list.IsEmpty())
	list.Clear()
	require.True(t, list.IsEmpty())
}
