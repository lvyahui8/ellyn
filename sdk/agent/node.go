package agent

import "sync"

var nodePool = &sync.Pool{
	New: func() any {
		return &node{}
	},
}

type node struct {
	recursion bool
	methodId  uint32
	cost      int64
	blocks    []bool
	args      *[]any
	results   *[]any
}

func (n *node) Recycle() {
	n.recursion = false
	n.cost = 0
	n.blocks = n.blocks[:0]
	n.args = nil
	n.results = nil
	nodePool.Put(n)
}
