package ellyn_agent

var allMethods []*method

type method struct {
	id     uint32
	blocks []*block
}
