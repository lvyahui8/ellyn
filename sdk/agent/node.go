package agent

import "sync"

// nodePool 函数节点缓存
var nodePool = &sync.Pool{
	New: func() any {
		return newNode()
	},
}

func newNode() *node {
	return &node{}
}

// node 代表一个函数节点，一个函数如果被调用多次（循环、递归等），在graph中也只会有一个节点
type node struct {
	// recursion 函数是否存在递归调用，这里包括直接递归或者间接递归
	recursion bool
	// methodId 函数id
	methodId uint32
	// cost 函数执行耗时, 如果函数执行多次，cost是累加值
	cost int64
	// blocks 函数块的覆盖情况
	blocks []bool
	// args 函数参数列表
	args *[]any
	// results 函数返回值列表
	results *[]any
}

// Recycle 回收重置函数节点
func (n *node) Recycle() {
	n.recursion = false
	n.cost = 0
	n.blocks = n.blocks[:0]
	n.args = nil
	n.results = nil
	nodePool.Put(n)
}
