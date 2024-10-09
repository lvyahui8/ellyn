package collections

import (
	"sync"
)

type LinkedList[T any] struct {
	lock  sync.RWMutex
	head  *linkedNode[T]
	tail  *linkedNode[T]
	count int
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (ll *LinkedList[T]) Add(val T) {
	ll.lock.Lock()
	defer ll.lock.Unlock()

	newNode := &linkedNode[T]{
		val:  val,
		next: nil,
	}
	if ll.tail != nil {
		ll.tail.next = newNode
	}
	ll.tail = newNode
	if ll.head == nil {
		ll.head = newNode
	}
	ll.count++
}

func (ll *LinkedList[T]) Values() []T {
	ll.lock.RLock()
	defer ll.lock.RUnlock()

	var res []T
	cur := ll.head
	for cur != nil {
		res = append(res, cur.val)
		cur = cur.next
	}
	return res
}

func (ll *LinkedList[T]) Clear() {
	ll.lock.Lock()
	defer ll.lock.Unlock()

	ll.head = nil
	ll.tail = nil
	ll.count = 0
}

func (ll *LinkedList[T]) IsEmpty() bool {
	ll.lock.RLock()
	defer ll.lock.RUnlock()

	return ll.count == 0
}

func (ll *LinkedList[T]) Size() int {
	ll.lock.RLock()
	defer ll.lock.RUnlock()

	return ll.count
}

func (ll *LinkedList[T]) Remove(index int) {
	ll.lock.Lock()
	defer ll.lock.Unlock()

	if index < 0 || index >= ll.count {
		return
	}
	cur := ll.head
	var prev *linkedNode[T]
	for i := 0; cur != nil; i, cur = i+1, cur.next {
		if index == i {
			if prev != nil {
				prev.next = cur.next
			}
			if ll.head == cur {
				ll.head = cur.next
			}
			if ll.tail == cur {
				ll.tail = prev
			}
			ll.count--
			return
		}
		prev = cur
	}
}
