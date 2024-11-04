package agent

import (
	"embed"
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/guid"
)

var idGenerator = guid.NewGuidGenerator()

var globalCovered []bool

type Api interface {
	InitCtx(ctxId uint64, from uint32)
	GetCtx() *EllynCtx
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

func (agent *ellynAgent) InitCtx(ctxId uint64, from uint32) {
	ctx := ctxPool.Get().(*EllynCtx)
	ctx.id = ctxId
	ctx.g.id = ctxId
	origin := toEdge(from, 0)
	ctx.g.origin = &origin
	ctxLocal.Set(ctx)
}

func (agent *ellynAgent) GetCtx() *EllynCtx {
	res, exist := ctxLocal.Get()
	if !exist {
		trafficId := idGenerator.GenGUID()
		res = &EllynCtx{
			id:        trafficId,
			stack:     collections.NewUnsafeCompressedStack[*methodFrame](),
			g:         newGraph(trafficId),
			autoClear: true,
		}
		ctxLocal.Set(res)
	}
	return res
}

func (agent *ellynAgent) Push(ctx *EllynCtx, methodId uint32, params []any) {
	// 压栈
	if ctx.g.origin != nil && ctx.stack.Size() == 0 {
		*(ctx.g.origin) |= uint64(methodId)
	}
	f := &methodFrame{methodId: methodId}

	ctx.stack.Push(f)
	if params != nil && f.data != nil {
		// 只记录首次入栈的参数
		f.data.args = EncodeVars(params)
	}
}

func (agent *ellynAgent) Pop(ctx *EllynCtx, results []any) {
	// 弹栈，加到调用链
	pop := ctx.stack.Pop()
	top := ctx.stack.Top()
	if top != nil && pop.methodId == top.methodId {
		// 方法递归中，未完全弹出
		return
	}
	// 只记录首次入栈的参数
	pop.data.results = EncodeVars(results)
	// 记录调用链
	ctx.g.add(top, pop)
	if top == nil {
		// 已经完全弹空， 调用链路追加到队列
		coll.add(ctx.g)
		if ctx.autoClear {
			ctx.recycle()
		}
	}
}

func (agent *ellynAgent) SetBlock(ctx *EllynCtx, blockOffset, blockId int) {
	// 取栈顶元素，标记block覆盖请求
	top := ctx.stack.Top()
	(*top.data.blocks)[blockOffset] = true
	globalCovered[blockId] = true
}
