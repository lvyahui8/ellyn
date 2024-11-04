package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
)

type frameData struct {
	blocks    *[]bool
	recursion bool
	begin     int64
	args      *[]any
	results   *[]any
}

type methodFrame struct {
	methodId uint32
	data     *frameData
}

func (mf *methodFrame) Equals(value collections.Frame) bool {
	f, ok := value.(*methodFrame)
	return ok && mf.methodId == f.methodId
}

func (mf *methodFrame) Init() {
	mf.data = &frameData{
		blocks: newMethodBlockFlags(mf.methodId),
		begin:  currentTime().UnixMilli(),
	}
}

func (mf *methodFrame) ReEnter() {
	if !mf.data.recursion {
		mf.data.recursion = true
	}
}
