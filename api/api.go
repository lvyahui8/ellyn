package api

import "sync"

type EllynApi interface {
	// SetAutoClear
	// 是否在收集完毕之后自动清理掉ctx
	SetAutoClear(auto bool)
	// GetGraph 获取当前的链路数据
	// 手动获取链路，必须先设置AutoClear=false
	GetGraph() *Graph
	// GetGraphId 获取当前链路id
	GetGraphId() uint64
	// ClearCtx 手动清理ctx
	ClearCtx()
	// GetGraphCnt 获取自启动后累计收集的链路数
	GetGraphCnt() uint64
}

var Agent *agentProxy = &agentProxy{}

var _ EllynApi = (*agentProxy)(nil)

type agentProxy struct {
	initOnce sync.Once
	target   EllynApi
}

func (a *agentProxy) GetGraphCnt() uint64 {
	if a.target == nil {
		return 0
	}
	return a.target.GetGraphCnt()
}

func Init(target EllynApi) {
	if target == nil {
		return
	}
	Agent.initOnce.Do(func() {
		Agent.target = target
	})
}

func (a *agentProxy) ClearCtx() {
	if a.target == nil {
		return
	}
	a.target.ClearCtx()
}

func (a *agentProxy) GetGraphId() uint64 {
	if a.target != nil {
		return a.target.GetGraphId()
	}
	return 0
}

func (a *agentProxy) SetAutoClear(auto bool) {
	if a.target != nil {
		a.target.SetAutoClear(auto)
	}
}

func (a *agentProxy) GetGraph() *Graph {
	if a.target != nil {
		return a.target.GetGraph()
	}
	return nil
}
