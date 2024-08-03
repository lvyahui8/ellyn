package ellyn_agent

import (
	"fmt"
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"runtime"
)

var coll *collector = newCollector()

func newCollector() *collector {
	c := &collector{
		buffer: collections.NewRingBuffer(2048),
	}
	c.start()
	return c
}

type collector struct {
	buffer *collections.RingBuffer
}

func (c *collector) add(g *graph) {
	c.buffer.Enqueue(g)
}

func (c *collector) start() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				v, ok := c.buffer.Dequeue()
				if !ok {
					continue
				}
				g := v.(*graph)
				fmt.Printf("graph %s\n", utils.Marshal(g))
			}
		}()
	}
}
