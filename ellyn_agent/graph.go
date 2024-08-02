package ellyn_agent

type graph struct {
	nodes map[uint32]*node
	edges map[uint64]struct{}
}
