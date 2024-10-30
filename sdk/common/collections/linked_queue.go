package collections

import (
	"math"
	"sync"
)

type linkedNode[T any] struct {
	val  T
	next *linkedNode[T]
}

type LinkedQueue[T any] struct {
	head     *linkedNode[T]
	tail     *linkedNode[T]
	lock     *sync.Mutex
	putCond  *sync.Cond
	takeCond *sync.Cond
	capacity int
	count    int
}

func NewLinkedQueue[T any](capacity int) *LinkedQueue[T] {
	if capacity <= 0 {
		capacity = math.MaxInt
	}
	lock := &sync.Mutex{}
	return &LinkedQueue[T]{
		lock:     lock,
		putCond:  sync.NewCond(lock),
		takeCond: sync.NewCond(lock),
		capacity: capacity,
	}
}

func (l *LinkedQueue[T]) Enqueue(value T) (success bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	for l.count == l.capacity {
		l.putCond.Wait()
	}

	n := &linkedNode[T]{value, nil}
	if l.tail != nil {
		l.tail.next = n
	}
	l.tail = n
	if l.count == 0 {
		l.head = l.tail
	}
	l.count++

	l.takeCond.Signal()
	return true
}

func (l *LinkedQueue[T]) Dequeue() (value T, success bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	for l.count == 0 {
		l.takeCond.Wait()
	}

	if l.head == nil {
		success = false
		return
	}
	value = l.head.val
	success = true
	l.head = l.head.next
	l.count--

	l.putCond.Signal()
	return
}
