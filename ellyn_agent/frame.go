package ellyn_agent

import "github.com/lvyahui8/ellyn/ellyn_common/collections"

type methodFrame struct {
	methodId  uint32
	blocks    *collections.BitMap
	recursion bool
}

func (mf *methodFrame) Equals(value collections.Frame) bool {
	f, ok := value.(*methodFrame)
	return ok && mf.methodId == f.methodId
}

func (mf *methodFrame) Init() {
	mf.blocks = newMethodBlockBits(mf.methodId)
}

func (mf *methodFrame) ReEnter() {
	if !mf.recursion {
		mf.recursion = true
	}
}
