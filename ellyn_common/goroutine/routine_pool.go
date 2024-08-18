package goroutine

import "github.com/lvyahui8/ellyn/ellyn_common/collections"

type Runnable func()

type RoutinePool struct {
	buf        *collections.LinkedQueue[Runnable]
	routineNum int
	closed     bool
}

func NewRoutinePool(routineNum int) *RoutinePool {
	pool := &RoutinePool{
		buf:        collections.NewLinkedQueue[Runnable](0),
		routineNum: routineNum,
	}
	pool.init()
	return pool
}

func (p *RoutinePool) init() {
	for i := 0; i < p.routineNum; i++ {
		go func() {
			for {
				r, success := p.buf.Dequeue()
				if !success {
					continue
				}
				if p.closed {
					return
				}
				func() {
					defer func() {
						_ = recover()
					}()
					r()
				}()
			}
		}()
	}
}

func (p *RoutinePool) Submit(r Runnable) {
	_ = p.buf.Enqueue(r)
}

func (p *RoutinePool) Shutdown() {
	p.closed = true
}
