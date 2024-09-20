package collections

import "container/list"

type Stack[T any] interface {
	Push(val T)
	Pop() T
	Top() T
	Size() int
}

type UnsafeStack[T any] struct {
	elements  *list.List
	_padding0 [56]byte
}

func NewUnsafeStack[T any]() *UnsafeStack[T] {
	return &UnsafeStack[T]{
		elements: list.New(),
	}
}

func (u UnsafeStack[T]) Push(val T) {
	u.elements.PushBack(val)
}

func (u UnsafeStack[T]) Pop() (t T) {
	v := u.elements.Back()
	if v != nil {
		u.elements.Remove(v)
		t = v.Value.(T)
	}
	return
}

func (u UnsafeStack[T]) Top() (t T) {
	v := u.elements.Back()
	if v != nil {
		t = v.Value.(T)
	}
	return
}

func (u UnsafeStack[T]) Size() int {
	return u.elements.Len()
}

// UnsafeCompressedStack 非并发安全的Stack
type UnsafeCompressedStack[T Frame] struct {
	elements  *list.List
	_padding0 [56]byte
	count     int
	_padding1 [56]byte
}

type Frame interface {
	Equals(value Frame) bool
	Init()
	ReEnter()
}

type stackElement[T Frame] struct {
	val   T
	max   int
	count int
}

func NewUnsafeCompressedStack[T Frame]() *UnsafeCompressedStack[T] {
	return &UnsafeCompressedStack[T]{
		elements: list.New(),
	}
}

func (s *UnsafeCompressedStack[T]) Push(val T) {
	s.count++
	back := s.elements.Back()
	if back != nil {
		ele := back.Value.(*stackElement[T])
		if val.Equals(ele.val) {
			val.ReEnter()
			ele.max++
			ele.count++
			return
		}
	}
	val.Init()
	s.elements.PushBack(&stackElement[T]{val: val, max: 1, count: 1})
}

func (s *UnsafeCompressedStack[T]) Pop() (t T) {
	e := s.elements.Back()
	if e != nil {
		ele := e.Value.(*stackElement[T])
		if ele.count == 1 {
			s.elements.Remove(e)
		} else {
			ele.count--
		}
		s.count--
		t = ele.val
	}
	return
}

func (s *UnsafeCompressedStack[T]) Top() (t T) {
	e := s.elements.Back()
	if e != nil {
		t = e.Value.(*stackElement[T]).val
	}
	return
}

func (s *UnsafeCompressedStack[T]) Size() int {
	return s.count
}
