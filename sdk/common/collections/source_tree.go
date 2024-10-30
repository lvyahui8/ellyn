package collections

import (
	"path/filepath"
	"strings"
)

type SourceTree struct {
	nodeMap map[string]*SourceNode
	root    []*SourceNode
}

func NewSourceTree() *SourceTree {
	return &SourceTree{
		nodeMap: make(map[string]*SourceNode),
	}
}

func (st *SourceTree) Root() []*SourceNode {
	return st.root
}

func (st *SourceTree) Add(path string, key string) *SourceNode {
	path = strings.Trim(filepath.ToSlash(path), "/")
	n := st.nodeMap[path]
	if n != nil {
		return n
	}
	items := strings.Split(path, "/")
	n = &SourceNode{
		Title: items[len(items)-1],
		Key:   key,
	}

	if len(items) >= 2 {
		pPath := strings.Join(items[0:len(items)-1], "/")
		parent := st.Add(pPath, pPath)
		parent.Children = append(parent.Children, n)
	} else {
		st.root = append(st.root, n)
	}

	st.nodeMap[path] = n
	return n
}

type SourceNode struct {
	Title    string        `json:"title"`
	Key      string        `json:"key"`
	Children []*SourceNode `json:"children"`
	IsLeaf   bool          `json:"isLeaf"`
}
