package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"sync"
)

var ctxPool = &sync.Pool{
	New: func() any {
		return newEllynCtx()
	},
}

//var ctxLocal = &goroutine.RoutineLocal[*EllynCtx]{}

var ctxLocal = collections.NewNumberKeyConcurrentMap[uint64, *EllynCtx](4096)

var discardedCtx = &EllynCtx{}

func newEllynCtx() *EllynCtx {
	return &EllynCtx{
		stack:     collections.NewUnsafeUint32Stack(),
		autoClear: true,
	}
}

type EllynCtx struct {
	autoClear bool
	goid      uint64 // 当前协程id
	id        uint64
	stack     *collections.UnsafeUint32Stack
	g         *graph
}

func (c *EllynCtx) Snapshot() (id uint64, currentMethodId uint32) {
	top, _ := c.stack.Top()
	return c.id, top
}

func (c *EllynCtx) Recycle() {
	c.stack.Clear()
	c.g = nil
	c.autoClear = true
	ctxLocal.Delete(c.goid)
	ctxPool.Put(c)
}
