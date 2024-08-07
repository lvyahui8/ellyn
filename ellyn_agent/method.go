package ellyn_agent

import "github.com/lvyahui8/ellyn/ellyn_common/collections"

var allMethods []*method

type method struct {
	id     uint32
	blocks []*block
}

func newMethodBlockBits(methodId uint32) *collections.BitMap {
	return collections.NewBitMap(uint(len(allMethods[methodId].blocks)))
}
