package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"runtime"
	"time"
)

var graphCache *collections.LRUCache[uint64] = collections.NewLRUCache[uint64](1024)

var coll *collector = newCollector()

func newCollector() *collector {
	c := &collector{
		buffer: collections.NewRingBuffer[*graph](2048),
	}
	c.start()
	return c
}

type collector struct {
	buffer *collections.RingBuffer[*graph]
}

func (c *collector) add(g *graph) {
	c.buffer.Enqueue(g)
}

func (c *collector) start() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				g, ok := c.buffer.Dequeue()
				if !ok {
					// 避免取不到数据CPU空转
					time.Sleep(1 * time.Nanosecond)
					continue
				}
				graphCache.Set(g.id, g)
			}
		}()
	}
}
