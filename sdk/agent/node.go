package agent

type node struct {
	recursion bool
	methodId  uint32
	cost      int64
	blocks    []bool
	args      *[]any
	results   *[]any
}

func newNode(methodId uint32) *node {
	return &node{
		methodId: methodId,
		blocks:   newMethodBlockFlags(methodId),
	}
}
