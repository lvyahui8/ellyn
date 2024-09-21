package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/elly_api"
)

func init() {
	elly_api.Agent.Init(&ellynApiImpl{})
}

var _ elly_api.EllynApi = (*ellynApiImpl)(nil)

type ellynApiImpl struct {
}

func (e *ellynApiImpl) SetAutoClear(auto bool) {
	ctx := Agent.GetCtx()
	ctx.autoClear = auto
}

func (e *ellynApiImpl) GetGraph() *elly_api.Graph {
	return nil
}
