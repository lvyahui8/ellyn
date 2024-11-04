package agent

import (
	"github.com/lvyahui8/ellyn/api"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
)

func init() {
	api.Init(&ellynApiImpl{})
}

var _ api.EllynApi = (*ellynApiImpl)(nil)

type ellynApiImpl struct {
	Agent *ellynAgent
}

func (e *ellynApiImpl) ClearCtx() {
	ctx := e.Agent.GetCtx()
	ctx.recycle()
}

func (e *ellynApiImpl) GetGraphId() uint64 {
	ctx := e.Agent.GetCtx()
	g := ctx.g
	return g.id
}

func (e *ellynApiImpl) SetAutoClear(auto bool) {
	ctx := e.Agent.GetCtx()
	ctx.autoClear = auto
}

func (e *ellynApiImpl) GetGraph() *api.Graph {
	ctx := e.Agent.GetCtx()
	g := ctx.g
	res := &api.Graph{
		Nodes: make(map[uint32]*api.Node),
		Edges: make(map[uint64]struct{}),
	}
	res.Edges = utils.CopyMap(g.edges)
	for methodId, n := range g.nodes {
		method := methods[n.methodId]
		file := files[method.FileId]
		pkg := packages[method.PackageId]
		res.Nodes[methodId] = &api.Node{
			MethodId:   methodId,
			MethodName: method.FullName,
			File:       file.RelativePath,
			Package:    pkg.Path,
			Begin:      api.Pos{Line: method.Begin.Line, Column: method.Begin.Column, Offset: method.Begin.Offset},
			End:        api.Pos{Line: method.End.Line, Column: method.End.Column, Offset: method.End.Offset},
		}
	}
	return res
}
