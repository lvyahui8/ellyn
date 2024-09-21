package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/lvyahui8/ellyn/ellyn_common/guid"
)

var Agent *ellynAgent = &ellynAgent{}

var idGenerator = guid.NewGuidGenerator()

// 实例
type ellynAgent struct {
}

func init() {
	initMetaData()
}

func (agent *ellynAgent) GetCtx() *EllynCtx {
	res, exist := CtxLocal.Get()
	if !exist {
		res = &EllynCtx{
			id:        idGenerator.GenGUID(),
			stack:     collections.NewUnsafeCompressedStack[*methodFrame](),
			g:         newGraph(),
			autoClear: true,
		}
		CtxLocal.Set(res)
	}
	return res
}

func (agent *ellynAgent) Push(ctx *EllynCtx, methodId uint32) {
	// 压栈
	ctx.stack.Push(&methodFrame{methodId: methodId})
}

func (agent *ellynAgent) Pop(ctx *EllynCtx) {
	// 弹栈，加到调用链
	pop := ctx.stack.Pop()
	top := ctx.stack.Pop()
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
