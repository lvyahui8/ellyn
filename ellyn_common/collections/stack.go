package collections

import "container/list"

type Stack interface {
	Push(val interface{})
	Pop() interface{}
	Top() interface{}
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

func (u UnsafeStack) Push(val interface{}) {
	u.elements.PushBack(val)
}

func (u UnsafeStack) Pop() interface{} {
	v := u.elements.Back()
	if v != nil {
		u.elements.Remove(v)
		return v.Value
	}
	return nil
}

func (u UnsafeStack) Top() interface{} {
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

type stackElement struct {
	val   interface{}
	max   int
	count int
}

func NewUnsafeCompressedStack() *UnsafeCompressedStack {
	return &UnsafeCompressedStack{
		elements: list.New(),
	}
}

func (s *UnsafeCompressedStack) Push(val interface{}) {
	s.count++
	back := s.elements.Back()
	if back != nil {
		ele := back.Value.(*stackElement)
		if ele.val == val {
			ele.max++
			ele.count++
			return
		}
	}
	s.elements.PushBack(&stackElement{val: val, max: 1, count: 1})
}

func (s *UnsafeCompressedStack) Pop() interface{} {
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

func (s *UnsafeCompressedStack) Top() interface{} {
	e := s.elements.Back()
	if e != nil {
		return e.Value.(*stackElement).val
	}
	return nil
}

func (s *UnsafeCompressedStack) Size() int {
	return s.count
}
