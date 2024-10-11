package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/lvyahui8/ellyn/ellyn_common/guid"
)

var Agent *ellynAgent = &ellynAgent{}

var idGenerator = guid.NewGuidGenerator()

var globalCovered *collections.BitMap

// 实例
type ellynAgent struct {
}

func init() {
	initMetaData()
	if len(blocks) > 0 {
		globalCovered = collections.NewBitMap(uint(len(blocks)))
	}
}

func (agent *ellynAgent) InitCtx(ctxId uint64, from uint32) {
	ctx := &EllynCtx{
		id:        ctxId,
		stack:     collections.NewUnsafeCompressedStack[*methodFrame](),
		g:         newGraph(ctxId),
		autoClear: true,
	}
	origin := toEdge(from, 0)
	ctx.g.origin = &origin
	CtxLocal.Set(ctx)
}

func (agent *ellynAgent) GetCtx() *EllynCtx {
	res, exist := CtxLocal.Get()
	if !exist {
		trafficId := idGenerator.GenGUID()
		res = &EllynCtx{
			id:        trafficId,
			stack:     collections.NewUnsafeCompressedStack[*methodFrame](),
			g:         newGraph(trafficId),
			autoClear: true,
		}
		CtxLocal.Set(res)
	}
	return res
}

func (agent *ellynAgent) Push(ctx *EllynCtx, methodId uint32) {
	// 压栈
	if ctx.g.origin != nil && ctx.stack.Size() == 0 {
		*(ctx.g.origin) |= uint64(methodId)
	}
	ctx.stack.Push(&methodFrame{methodId: methodId})
}

func (agent *ellynAgent) Pop(ctx *EllynCtx) {
	// 弹栈，加到调用链
	pop := ctx.stack.Pop()
	top := ctx.stack.Top()
	if top != nil && pop.methodId == top.methodId {
		// 方法递归中，未完全弹出
		return
	}
	// 记录调用链
	ctx.g.add(top, pop)
	if top == nil {
		// 已经完全弹空， 调用链路追加到队列
		coll.add(ctx.g)
		if ctx.autoClear {
			CtxLocal.Clear()
		}
	}
}

func (agent *ellynAgent) VisitBlock(ctx *EllynCtx, blockOffset int) {
	// 取栈顶元素，标记block覆盖请求
	top := ctx.stack.Top()
	top.blocks.Set(uint(blockOffset))
}
