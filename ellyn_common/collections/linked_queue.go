package collections

type linkedNode[T any] struct {
	val  T
	next *linkedNode[T]
}

type LinkedQueue[T any] struct {
	head *linkedNode[T]
	tail *linkedNode[T]
}

func (l LinkedQueue[T]) Enqueue(value T) (success bool) {
	n := &linkedNode[T]{value, nil}
	if l.tail != nil {
		l.tail.next = n
	}
	l.tail = n
	return true
}

func (l LinkedQueue[T]) Dequeue() (value T, success bool) {
	if l.head == nil {
		success = false
		return
	}
	value = l.head.val
	success = true
	l.head = l.head.next
	return
}
