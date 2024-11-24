package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/logging"
	"sync/atomic"
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

var graphCnt uint64

func (c *collector) start() {
	go func() {
		for {
			g, ok := c.buffer.Dequeue()
			if !ok {
				// 避免取不到数据CPU空转
				time.Sleep(time.Microsecond)
				continue
			}
			// fmt.Printf("g:%d\n", g.id)
			atomic.AddUint64(&graphCnt, 1)
			if !conf.NoDemo {
				// 消费链路数据，这里缓存到本地用于demo显示
				log.InfoKV(logging.Code("g_collect").Int("n", len(g.nodes)).
					Int("e", len(g.edges)).Bool("c", g.origin != nil))
				saveToDisplayCache(g)
			} else {
				// 省略实际消费g，比如写磁盘，或者上报MQ
				// 消费完成后回收g
				g.Recycle()
			}
		}
	}()
}

func saveToDisplayCache(g *graph) {
	// 同一个id可能因为异步产生多个graph，需要merge. 这里放到展示的时候再merge
	graphCache.GetWithDefault(g.id, func() *graphGroup {
		return &graphGroup{
			list: collections.NewLinkedList[*graph](),
		}
	}).list.Add(g)
}
