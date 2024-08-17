package collections

type Runnable interface {
	Run()
}

type RoutinePool struct {
	buf         *RingBuffer
	consumerNum int
	closed      bool
}

func (p *RoutinePool) init() {
	for i := 0; i < p.consumerNum; i++ {
		go func() {
			for {
				item, success := p.buf.Dequeue()
				if !success {
					continue
				}
				if p.closed {
					return
				}
				r := item.(Runnable)
				func() {
					defer func() {
						_ = recover()
					}()
					r.Run()
				}()
			}
		}()
	}
}

func (p *RoutinePool) Submit(r Runnable) bool {
	return p.buf.Enqueue(r)
}

func (p *RoutinePool) Shutdown() {
	p.closed = true
}
