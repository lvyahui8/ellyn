package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
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

//var ctxLocal = &goroutine.RoutineLocal[*EllynCtx]{}

var ctxLocal = collections.NewNumberKeyConcurrentMap[uint64, *EllynCtx](4096)

type EllynCtx struct {
	goid      uint64 // 当前协程id
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
	c.g = graphPool.Get().(*graph)
	c.autoClear = true
	ctxPool.Put(c)
	// todo 删除
}
