package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/goroutine"
	"sync"
	"unsafe"
)

var ctxPool = &sync.Pool{
	New: func() any {
		return newEllynCtx()
	},
}

var discardedCtx = &EllynCtx{}

//var ctxLocal = &goroutine.RoutineLocal[*EllynCtx]{}
// ctxLocal map[goid]*EllynCtx
// getCtx高频调用，使用map存对性能影响比较明显
var ctxLocal = collections.NewNumberKeyConcurrentMap[uint64, *EllynCtx](4096)

// getEllynCtx
// 如果使用map实现，对应 ctx, exist := ctxLocal.Load(goid)
func getEllynCtx() (ctx *EllynCtx, exist bool) {
	ptr := goroutine.GetRoutineCtx()
	if ptr == 0 {
		return nil, false
	}
	return (*EllynCtx)(unsafe.Pointer(ptr)), true
}

// setEllynCtx
// 对应ctxLocal.Store(goid, ctx)
func setEllynCtx(ctx *EllynCtx) {
	goroutine.SetRoutineCtx(uintptr(unsafe.Pointer(ctx)))
}

// clearEllynCtx
// 对应 ctxLocal.Delete(goid)
func clearEllynCtx() {
	goroutine.SetRoutineCtx(0)
}

func newEllynCtx() *EllynCtx {
	return &EllynCtx{
		stack:     collections.NewUnsafeUint32Stack(),
		autoClear: true,
	}
}

type EllynCtx struct {
	autoClear bool
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
	clearEllynCtx()
	ctxPool.Put(c)
}
