package collections

import "container/list"

type Stack interface {
	Push(val any)
	Pop() any
	Top() any
	Size() int
}

type UnsafeStack struct {
	elements *list.List
}

func NewUnsafeStack() *UnsafeStack {
	return &UnsafeStack{
		elements: list.New(),
	}
}

func (u UnsafeStack) Push(val any) {
	u.elements.PushBack(val)
}

func (u UnsafeStack) Pop() any {
	v := u.elements.Back()
	if v != nil {
		u.elements.Remove(v)
		return v.Value
	}
	return nil
}

func (u UnsafeStack) Top() any {
	v := u.elements.Back()
	if v != nil {
		return v.Value
	} else {
		return nil
	}
}

func (u UnsafeStack) Size() int {
	return u.elements.Len()
}

// UnsafeCompressedStack 非并发安全的Stack
type UnsafeCompressedStack struct {
	elements *list.List
	count    int
}

type Frame interface {
	Equals(value any) bool
	Init()
}

type stackElement struct {
	val   any
	max   int
	count int
}

func NewUnsafeCompressedStack() *UnsafeCompressedStack {
	return &UnsafeCompressedStack{
		elements: list.New(),
	}
}

func (s *UnsafeCompressedStack) Push(val any) {
	s.count++
	back := s.elements.Back()
	f, _ := val.(Frame)

	if back != nil {
		ele := back.Value.(*stackElement)
		var eq bool
		if f != nil {
			eq = f.Equals(ele.val)
		} else {
			eq = ele.val == val
		}
		if eq {
			ele.max++
			ele.count++
			return
		}
	}
	if f != nil {
		f.Init()
	}
	s.elements.PushBack(&stackElement{val: val, max: 1, count: 1})
}

func (s *UnsafeCompressedStack) Pop() any {
	e := s.elements.Back()
	if e != nil {
		ele := e.Value.(*stackElement)
		if ele.count == 1 {
			s.elements.Remove(e)
		} else {
			ele.count--
		}
		s.count--
		return ele.val
	} else {
		return nil
	}
}

func (s *UnsafeCompressedStack) Top() any {
	e := s.elements.Back()
	if e != nil {
		return e.Value.(*stackElement).val
	}
	return nil
}

func (s *UnsafeCompressedStack) Size() int {
	return s.count
}
