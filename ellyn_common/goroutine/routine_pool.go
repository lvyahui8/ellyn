package goroutine

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
)

type Runnable func()

type RoutinePool struct {
	buf         *collections.LinkedQueue[Runnable]
	routineNum  int
	closed      bool
	ignorePanic bool
}

func NewRoutinePool(routineNum int, ignorePanic bool) *RoutinePool {
	pool := &RoutinePool{
		buf:         collections.NewLinkedQueue[Runnable](0),
		routineNum:  routineNum,
		ignorePanic: ignorePanic,
	}
	pool.init()
	return pool
}

func (p *RoutinePool) init() {
	for i := 0; i < p.routineNum; i++ {
		go func() {
			for {
				if p.closed {
					return
				}
				r, success := p.buf.Dequeue()
				if !success {
					continue
				}
				func() {
					if p.ignorePanic {
						defer func() {
							_ = recover()
						}()
					}
					r()
				}()
			}
		}()
	}
}

func (p *RoutinePool) Submit(r Runnable) {
	if p.closed {
		panic("routine pool has closed")
	}
	_ = p.buf.Enqueue(r)
}

func (p *RoutinePool) Shutdown() {
	p.closed = true
	go func() {
		// 清空队列
		for {
			_, _ = p.buf.Dequeue()
		}
	}()
}
