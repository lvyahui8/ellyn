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
			g:         graphPool.Get().(*graph),
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
	c.g = nil
	c.autoClear = true
	ctxPool.Put(c)
	ctxLocal.Clear()
}
