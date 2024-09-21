package ellyn_api

import "sync"

type EllynApi interface {
	// SetAutoClear
	// 是否在收集完毕之后自动清理掉ctx
	SetAutoClear(auto bool)
	// GetGraph 获取当前的链路数据
	// 手动获取链路，必须先设置AutoClear=false
	GetGraph() *Graph
}

var Agent *agentProxy = &agentProxy{}

type agentProxy struct {
	sync.Once
	target EllynApi
}

func (a *agentProxy) Init(target EllynApi) {
	if target == nil {
		return
	}
	a.Once.Do(func() {
		a.target = target
	})
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
