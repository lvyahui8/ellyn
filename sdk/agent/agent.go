package agent

import (
	"embed"
	"github.com/lvyahui8/ellyn/sdk/common/guid"
	"unsafe"
)

var idGenerator = guid.NewGuidGenerator()

var globalCovered []bool

type Api interface {
	InitCtx(ctxId uint64, from uint32)
	GetCtx() (ctx *EllynCtx, collect bool, cleaner func())
	Push(ctx *EllynCtx, methodId uint32, params []any)
	Pop(ctx *EllynCtx, results []any)
	SetBlock(ctx *EllynCtx, blockOffset, blockId int)
}

var _ Api = (*ellynAgent)(nil)

// ellynAgent 实例
type ellynAgent struct {
}

func InitAgent(meta embed.FS) Api {
	initAgent(meta)
	return &ellynAgent{}
}

func (agent *ellynAgent) InitCtx(trafficId uint64, from uint32) {
	ctx := ctxPool.Get().(*EllynCtx)

	ctx.id = trafficId
	ctx.g = graphPool.Get().(*graph)
	ctx.g.id = trafficId
	origin := toEdge(from, 0)
	ctx.g.origin = &origin
	setEllynCtx(ctx)
}

func (agent *ellynAgent) GetCtx() (ctx *EllynCtx, collect bool, cleaner func()) {
	ctx, exist := getEllynCtx()
	if !exist {
		defer func() {
			setEllynCtx(ctx)
		}()
		if !sampling.hit() {
			ctx = discardedCtx
			return ctx, false, func() {
				clearEllynCtx()
			}
		}
		ctx = ctxPool.Get().(*EllynCtx)
		trafficId := idGenerator.GenGUID()
		ctx.id = trafficId
		ctx.g = graphPool.Get().(*graph)
		ctx.g.id = trafficId
	}
	return ctx, ctx != discardedCtx, nil
}

func (agent *ellynAgent) Push(ctx *EllynCtx, methodId uint32, params []any) {
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
			//log.Info("ctx.g %x", uintptr(unsafe.Pointer(ctx.g)))
			ctx.g.nodes[methodId] = n
		}
		// 后续使用不用再查找
		ctx.stack.SetTopExtra(uintptr(unsafe.Pointer(n)))
	}
}

func (agent *ellynAgent) Pop(ctx *EllynCtx, results []any) {
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

func (agent *ellynAgent) SetBlock(ctx *EllynCtx, blockOffset, blockId int) {
	// 取栈顶元素，标记block覆盖请求
	extra := ctx.stack.GetTopExtra()
	((*node)(unsafe.Pointer(extra))).blocks[blockOffset] = true
	globalCovered[blockId] = true
}
