package collections

import (
	"container/list"
	"sync"
)

type Stack[T any] interface {
	Push(val T) (newElem bool)
	Pop() (T, bool)
	Top() (T, bool)
	Size() int
	Empty() bool
	Clear()
}

var _ Stack[any] = (*UnsafeStack[any])(nil)
var _ Stack[Frame] = (*UnsafeCompressedStack[Frame])(nil)
var _ Stack[uint32] = (*UnsafeUint32Stack)(nil)

type UnsafeStack[T any] struct {
	elements  *list.List
	_padding0 [56]byte
}

func NewUnsafeStack[T any]() *UnsafeStack[T] {
	return &UnsafeStack[T]{
		elements: list.New(),
	}
}

func (u *UnsafeStack[T]) Push(val T) (newElem bool) {
	u.elements.PushBack(val)
	return true
}

func (u *UnsafeStack[T]) Pop() (t T, ok bool) {
	v := u.elements.Back()
	if v != nil {
		u.elements.Remove(v)
		t = v.Value.(T)
		ok = true
	}
	return
}

func (u *UnsafeStack[T]) Top() (t T, ok bool) {
	v := u.elements.Back()
	if v != nil {
		ok = true
		t = v.Value.(T)
	}
	return
}

func (u *UnsafeStack[T]) Size() int {
	return u.elements.Len()
}

func (u *UnsafeStack[T]) Empty() bool {
	return u.elements.Len() == 0
}

func (u *UnsafeStack[T]) Clear() {
	u.elements.Init()
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
	// 栈元素值
	val T
	// 当前在栈中元素个数
	count int
}

func NewUnsafeCompressedStack[T Frame]() *UnsafeCompressedStack[T] {
	return &UnsafeCompressedStack[T]{
		elements: list.New(),
	}
}

func (s *UnsafeCompressedStack[T]) Push(val T) bool {
	s.count++
	back := s.elements.Back()
	if back != nil {
		ele := back.Value.(*stackElement[T])
		if val.Equals(ele.val) {
			ele.val.ReEnter()
			ele.count++
			return false
		}
	}
	val.Init()
	s.elements.PushBack(&stackElement[T]{val: val, count: 1})
	return true
}

func (s *UnsafeCompressedStack[T]) Pop() (t T, ok bool) {
	e := s.elements.Back()
	if e != nil {
		ele := e.Value.(*stackElement[T])
		ele.count--
		if ele.count == 0 {
			s.elements.Remove(e)
		}
		s.count--
		t = ele.val
		ok = true
	}
	return
}

func (s *UnsafeCompressedStack[T]) Top() (t T, ok bool) {
	e := s.elements.Back()
	if e != nil {
		ok = true
		t = e.Value.(*stackElement[T]).val
	}
	return
}

func (s *UnsafeCompressedStack[T]) Size() int {
	return s.count
}

func (s *UnsafeCompressedStack[T]) Empty() bool {
	//return s.elements.Len() == 0
	return s.count == 0
}

func (s *UnsafeCompressedStack[T]) Clear() {
	s.elements.Init()
}

type uint32Node struct {
	prev  *uint32Node
	ele   uint64 // 高32位 val，低32位 count
	extra uintptr
}

var nodePool = &sync.Pool{New: func() any { return &uint32Node{} }}

type UnsafeUint32Stack struct {
	tail *uint32Node
}

func NewUnsafeUint32Stack() *UnsafeUint32Stack {
	return &UnsafeUint32Stack{}
}

func (u *UnsafeUint32Stack) Push(val uint32) (newElem bool) {
	if u.tail != nil {
		if uint32(u.tail.ele>>32) == val {
			// same
			u.tail.ele += 1
			return false
		}
	}
	n := nodePool.Get().(*uint32Node)
	n.ele = uint64(val)<<32 | 1
	n.prev = u.tail
	u.tail = n
	return true
}

func (u *UnsafeUint32Stack) Pop() (val uint32, suc bool) {
	val, _, suc = u.PopWithExtra()
	return
}

func (u *UnsafeUint32Stack) PopWithExtra() (val uint32, extra uintptr, suc bool) {
	if u.tail == nil {
		return 0, 0, false
	}
	val = uint32(u.tail.ele >> 32)
	extra = u.tail.extra
	if uint32(u.tail.ele) == 1 {
		// real pop
		top := u.tail
		u.tail = top.prev
		nodePool.Put(top)
	} else {
		u.tail.ele--
	}
	suc = true
	return
}

func (u *UnsafeUint32Stack) Top() (val uint32, suc bool) {
	val, _, suc = u.TopWithExtra()
	return
}

func (u *UnsafeUint32Stack) TopWithExtra() (val uint32, extra uintptr, suc bool) {
	if u.tail == nil {
		return 0, 0, false
	}
	return uint32(u.tail.ele >> 32), u.tail.extra, true
}

func (u *UnsafeUint32Stack) Size() int {
	panic("unsupported")
}

func (u *UnsafeUint32Stack) Empty() bool {
	return u.tail == nil
}

func (u *UnsafeUint32Stack) Clear() {
	// fast clear
	u.tail = nil
}

func (u *UnsafeUint32Stack) SetTopExtra(extra uintptr) {
	if u.tail != nil {
		u.tail.extra = extra
	}
}

func (u *UnsafeUint32Stack) GetTopExtra() uintptr {
	if u.tail != nil {
		return u.tail.extra
	}
	return 0
}
