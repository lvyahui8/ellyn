package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/lvyahui8/ellyn/ellyn_common/gls"
)

var CtxLocal = &gls.RoutineLocal[*EllynCtx]{}

type EllynCtx struct {
	stack collections.Stack
	g     *graph
}
