package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/goroutine"
	"sync"
	"unsafe"
)

// ctxPool ctx池化复用，降低gc对性能的影响
var ctxPool = &sync.Pool{
	New: func() any {
		return newEllynCtx()
	},
}

// discardedCtx 创建监控服务
var discardedCtx = &EllynCtx{}

//var ctxLocal = &goroutine.RoutineLocal[*EllynCtx]{}
// ctxLocal map[goid]*EllynCtx
// getCtx高频调用，使用map存对性能影响比较明显. 这里改为直接将ctx挂在runtime.g上
// var ctxLocal = collections.NewNumberKeyConcurrentMap[uint64, *EllynCtx](4096)

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

// newEllynCtx 初始化一个空的ctx
func newEllynCtx() *EllynCtx {
	return &EllynCtx{
		stack:     collections.NewUnsafeUint32Stack(),
		autoClear: true,
	}
}

// EllynCtx agent数据核心，用于存储当前协程收集到的所有链路数据、模拟栈、状态数据等
type EllynCtx struct {
	// autoClear 是否在栈完全弹空时自动清理
	autoClear bool
	// id 链路id
	id uint64
	// stack 模拟栈，用于构建调用链
	stack *collections.UnsafeUint32Stack
	// g 当前已经采集的调用链数据
	g *graph
}

// Snapshot 当前ctx的快照数据
func (c *EllynCtx) Snapshot() (id uint64, currentMethodId uint32) {
	if c.stack == nil {
		return c.id, 0
	}
	top, _ := c.stack.Top()
	return c.id, top
}

// Recycle 重置&回收ctx
func (c *EllynCtx) Recycle() {
	c.stack.Clear()
	c.g = nil
	c.autoClear = true
	clearEllynCtx()
	ctxPool.Put(c)
}
