package ellyn_agent

import "github.com/lvyahui8/ellyn/ellyn_common/collections"

type node struct {
	methodId int
	blocks   *collections.BitMap
}
