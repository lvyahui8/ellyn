package ellyn_agent

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"runtime"
	"time"
)

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
				updateGlobalCovered(g)
				// 消费链路数据，这里缓存到本地用于显示
				saveToDisplayCache(g)
			}
		}()
	}
}

func saveToDisplayCache(g *graph) {
	// 同一个id可能因为异步产生多个graph，需要merge. 这里放到展示的时候再merge
	graphCache.GetWithDefault(g.id, func() any {
		return collections.NewLinkedList[*graph]()
	}).(*collections.LinkedList[*graph]).Add(g)
}

func updateGlobalCovered(g *graph) {
	for _, n := range g.nodes {
		m := methods[n.methodId]
		for _, block := range m.Blocks {
			if n.blocks.Get(uint(block.MethodOffset)) {
				globalCovered.Set(uint(block.Id))
			}
		}
	}
}
