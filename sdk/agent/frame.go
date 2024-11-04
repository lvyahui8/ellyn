package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
)

type methodFrame struct {
	methodId  uint32
	blocks    []bool
	recursion bool
	begin     int64
	args      []any
	results   []any
}

func (mf *methodFrame) Equals(value collections.Frame) bool {
	f, ok := value.(*methodFrame)
	return ok && mf.methodId == f.methodId
}

func (mf *methodFrame) Init() {
	mf.blocks = newMethodBlockFlags(mf.methodId)
	mf.begin = currentTime().UnixMilli()
}

func (mf *methodFrame) ReEnter() {
	if !mf.recursion {
		mf.recursion = true
	}
}
