package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/goroutine"
	"sync"
)

var ctxPool = &sync.Pool{
	New: func() any {
		return &EllynCtx{
			stack:     collections.NewUnsafeCompressedStack[*methodFrame](),
			g:         newGraph(0),
			autoClear: true,
		}
	},
}
var ctxLocal = &goroutine.RoutineLocal[*EllynCtx]{}

type EllynCtx struct {
	id        uint64
	stack     collections.Stack[*methodFrame]
	g         *graph
	autoClear bool
}

func (c *EllynCtx) Snapshot() (id uint64, currentMethodId uint32) {
	return c.id, c.stack.Top().methodId
}

func (c *EllynCtx) recycle() {
	c.stack.Clear()
	c.g = newGraph(0) // g 暂时还不能复用，后续还需要使用
	ctxPool.Put(c)
	ctxLocal.Clear()
}
