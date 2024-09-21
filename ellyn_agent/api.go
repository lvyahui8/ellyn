package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_api"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
)

func init() {
	ellyn_api.Agent.Init(&ellynApiImpl{})
}

var _ ellyn_api.EllynApi = (*ellynApiImpl)(nil)

type ellynApiImpl struct {
}

func (e *ellynApiImpl) SetAutoClear(auto bool) {
	ctx := Agent.GetCtx()
	ctx.autoClear = auto
}

func (e *ellynApiImpl) GetGraph() *ellyn_api.Graph {
	ctx := Agent.GetCtx()
	g := ctx.g
	res := &ellyn_api.Graph{
		Nodes: make(map[uint32]*ellyn_api.Node),
		Edges: make(map[uint64]struct{}),
	}
	res.Edges = utils.CopyMap(g.edges)
	for methodId, n := range g.nodes {
		method := methods[n.methodId]
		file := files[method.FileId]
		pkg := packages[method.PackageId]
		res.Nodes[methodId] = &ellyn_api.Node{
			MethodId:   methodId,
			MethodName: method.FullName,
			File:       file.RelativePath,
			Package:    pkg.Path,
			Begin:      ellyn_api.Pos{Line: method.Begin.Line, Column: method.Begin.Column, Offset: method.Begin.Offset},
			End:        ellyn_api.Pos{Line: method.End.Line, Column: method.End.Column, Offset: method.End.Offset},
		}
	}
	return res
}
