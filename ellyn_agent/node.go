package ellyn_agent

import "github.com/lvyahui8/ellyn/ellyn_common/collections"

type node struct {
	methodId uint32
	blocks   *collections.BitMap
	cost     int64
	args     []any
	results  []any
}
