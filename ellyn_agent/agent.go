package ellyn_agent

var Agent *ellynAgent = &ellynAgent{}

// 实例
type ellynAgent struct {
}

func (agent *ellynAgent) GetCtx() *EllynCtx {
	res, exist := CtxLocal.Get()
	if !exist {
		res = &EllynCtx{}
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
	pop := ctx.stack.Pop().(*methodFrame)
	top := ctx.stack.Pop().(*methodFrame)
	if top != nil {
		if pop.methodId == top.methodId {
			// 方法递归中，未完全弹出
			return
		}
		// 记录调用链
		ctx.g.add(top, pop)
	} else {
		// 已经完全弹空， 调用链路追加到队列
		coll.add(ctx.g)
	}
}

func (agent *ellynAgent) VisitBlock(ctx *EllynCtx, blockOffset int) {
	// 取栈顶元素，标记block覆盖请求
	top := ctx.stack.Top().(*methodFrame)
	top.blocks.Set(uint(blockOffset))
}
