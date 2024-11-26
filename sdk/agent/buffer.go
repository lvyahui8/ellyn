package agent

import (
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"github.com/lvyahui8/ellyn/sdk/common/logging"
	"sync/atomic"
	"time"
)

// coll 链路数据缓冲
// 基于Lock-Free RingBuffer缓冲链路数据
var coll *collector = newCollector()

// graphCnt 记录启动后累计收集的链路数
var graphCnt uint64

func newCollector() *collector {
	c := &collector{
		buffer: collections.NewRingBuffer[*graph](2048),
	}
	c.start()
	return c
}

// collector 收集器
type collector struct {
	buffer *collections.RingBuffer[*graph]
}

// add 将链路加入到缓冲队列中
func (c *collector) add(g *graph) {
	c.buffer.Enqueue(g)
}

// start 启动消费逻辑
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

// saveToDisplayCache 将链路数据保存到本地的LRU缓存，用于demo程序可视化展示
func saveToDisplayCache(g *graph) {
	// 同一个id可能因为异步产生多个graph，需要merge. 这里放到展示的时候再merge
	graphCache.GetWithDefault(g.id, func() *graphGroup {
		return &graphGroup{
			list: collections.NewLinkedList[*graph](),
		}
	}).list.Add(g)
}
