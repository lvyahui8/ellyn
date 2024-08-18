package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/lvyahui8/ellyn/ellyn_common/goroutine"
)

var CtxLocal = &goroutine.RoutineLocal[*EllynCtx]{}

type EllynCtx struct {
	stack collections.Stack
	g     *graph
}
