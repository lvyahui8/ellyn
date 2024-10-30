package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/goroutine"
)

var CtxLocal = &goroutine.RoutineLocal[*EllynCtx]{}

type EllynCtx struct {
	id        uint64
	stack     collections.Stack[*methodFrame]
	g         *graph
	autoClear bool
}

func (c *EllynCtx) Snapshot() (id uint64, currentMethodId uint32) {
	return c.id, c.stack.Top().methodId
}
