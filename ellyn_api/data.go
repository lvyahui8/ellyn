package ellyn_api

type Pos struct {
	Line   int
	Column int
	Offset int
}

type Graph struct {
	Nodes map[uint32]*Node
	Edges map[uint64]struct{}
}

type Node struct {
	MethodId   uint32
	MethodName string
	File       string
	Package    string
	Begin      Pos
	End        Pos
}
