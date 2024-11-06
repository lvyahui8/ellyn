package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/goroutine"
	"sync"
)

var ctxPool = &sync.Pool{
	New: func() any {
		return &EllynCtx{
			stack:     collections.NewUnsafeUint32Stack(),
			g:         newGraph(0),
			autoClear: true,
		}
	},
}
var ctxLocal = &goroutine.RoutineLocal[*EllynCtx]{}

type EllynCtx struct {
	id        uint64
	stack     *collections.UnsafeUint32Stack
	g         *graph
	autoClear bool
}

func (c *EllynCtx) Snapshot() (id uint64, currentMethodId uint32) {
	top, _ := c.stack.Top()
	return c.id, top
}

func (c *EllynCtx) recycle() {
	c.stack.Clear()
	c.g = newGraph(0) // g 暂时还不能复用，后续还需要使用
	ctxPool.Put(c)
	ctxLocal.Clear()
}
