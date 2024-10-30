package agent

import "github.com/lvyahui8/ellyn/sdk/common/collections"

type node struct {
	methodId uint32
	blocks   *collections.BitMap
	cost     int64
	args     []any
	results  []any
}
