package collections

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSourceTree1(t *testing.T) {
	tree := NewSourceTree()
	tree.Add("./a/b/c.txt", "1")
	tree.Add("./a/d.txt", "2")
	tree.Add("./a/b/h.txt", "3")
	tree.Add("./x/m.txt", "4")
	require.Equal(t, 1, len(tree.Root()))
}

func TestSourceTree2(t *testing.T) {
	tree := NewSourceTree()
	tree.Add("a/b/c.txt", "1")
	tree.Add("a/d.txt", "2")
	tree.Add("a/b/h.txt", "3")
	tree.Add("x/m.txt", "4")
	require.Equal(t, 2, len(tree.Root()))
}
