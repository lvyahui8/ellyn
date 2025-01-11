package agent

import (
	"context"
	"github.com/lvyahui8/ellyn/sdk/common/guid"
	"github.com/lvyahui8/ellyn/sdk/common/logging"
	"unsafe"
)

// log
// log写盘逻辑：定时、大小溢出、日期切换
var log = logging.GetLogger()

// idGenerator 用于生成流量id
var idGenerator = guid.NewGuidGenerator()

// globalCovered 记录全局覆盖，长度为全局的块信息
var globalCovered []bool

// Api agent API,插桩到目标项目中使用的代码
type Api interface {
	// InitCtx 手动创建一个ctx，用于异步协程指定ctx来源
	InitCtx(traceId uint64, from uint32)
	// GetCtx 获取当前协程的ctx，如果没有则会初始化一个
	GetCtx() (ctx *EllynCtx, collect bool, cleaner func())
	// Push 方法压栈
	Push(ctx *EllynCtx, goCtx *context.Context, methodId uint32, params []any)
	// Pop 方法弹栈
	Pop(ctx *EllynCtx, results []any)
	// Mark 标记覆盖的块
	Mark(ctx *EllynCtx, blockOffset, blockId int)
}

// 用于限制ellynAgent必须实现Api接口
var _ Api = (*ellynAgent)(nil)

// ellynAgent agent api实现
type ellynAgent struct {
}

func (agent *ellynAgent) InitCtx(traceId uint64, from uint32) {
	defer handleSelfError()()
	if traceId == 0 {
		setEllynCtx(discardedCtx)
		return
	}
	ctx := ctxPool.Get().(*EllynCtx)
	ctx.id = traceId
	ctx.g = graphPool.Get().(*graph)
	ctx.g.id = traceId
	origin := toEdge(from, 0)
	ctx.g.origin = &origin
	setEllynCtx(ctx)
}

func (agent *ellynAgent) GetCtx() (ctx *EllynCtx, collect bool, cleaner func()) {
	defer handleSelfError()()
	ctx, exist := getEllynCtx()
	if !exist {
		if !sampling.hit() {
			ctx = discardedCtx
			setEllynCtx(ctx)
			collect = false
			cleaner = func() {
				clearEllynCtx()
			}
			return
		}
		ctx = ctxPool.Get().(*EllynCtx)
		traceId := idGenerator.GenGUID()
		ctx.id = traceId
		ctx.g = graphPool.Get().(*graph)
		ctx.g.id = traceId
		setEllynCtx(ctx)
	}
	collect = ctx != discardedCtx
	return
}

func (agent *ellynAgent) Push(ctx *EllynCtx, goCtx *context.Context, methodId uint32, params []any) {
	defer handleSelfError()()
	// 压栈
	if ctx.g.origin != nil && ctx.stack.Empty() {
		*(ctx.g.origin) |= uint64(methodId)
	}
	newElem := ctx.stack.Push(methodId)
	if newElem {
		// 当前method最新一次重入栈
		var n *node
		var exist bool
		if n, exist = ctx.g.nodes[methodId]; !exist {
			n = nodePool.Get().(*node)
			n.methodId = methodId
			n.blocks = newMethodBlockFlags(methodId)
			// 方法多次调用只记录第一次参数
			if !conf.NoArgs {
				n.args = EncodeVars(params)
			}
			//logging.Info("ctx.g %x", uintptr(unsafe.Pointer(ctx.g)))
			ctx.g.nodes[methodId] = n
		}
		// 后续使用不用再查找
		ctx.stack.SetTopExtra(uintptr(unsafe.Pointer(n)))
	}
}

func (agent *ellynAgent) Pop(ctx *EllynCtx, results []any) {
	defer handleSelfError()()
	// 弹栈，加到调用链
	pop, extra, _ := ctx.stack.PopWithExtra()
	top, ok := ctx.stack.Top()
	n := (*node)(unsafe.Pointer(extra))
	if ok && pop == top {
		// 方法递归中，未完全弹出
		if !n.recursion {
			n.recursion = true
			ctx.g.addEdge(pop, pop)
		}
		return
	}

	if !conf.NoArgs && n.results == nil {
		// 只记录首次产生节点的参数
		n.results = EncodeVars(results)
	}

	// 记录调用链
	if ok {
		ctx.g.addEdge(top, pop)
	} else {
		// 已经完全弹空， 调用链路追加到队列
		coll.add(ctx.g)
		if ctx.autoClear {
			ctx.Recycle()
		}
	}
}

func (agent *ellynAgent) Mark(ctx *EllynCtx, blockOffset, blockId int) {
	defer handleSelfError()()
	// 取栈顶元素，标记block覆盖请求
	extra := ctx.stack.GetTopExtra()
	((*node)(unsafe.Pointer(extra))).blocks[blockOffset] = true
	globalCovered[blockId] = true
}

func handleSelfError() func() {
	return func() {
		err := recover()
		if err != nil {
			// 可以在滑动时间窗口内，按照指数避退策略，减少错误日志输出，避免出现异常时打满磁盘
			log.Error("agent panic err: %v", err)
		}
	}
}
