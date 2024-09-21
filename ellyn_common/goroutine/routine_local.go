package goroutine

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
)

var m = collections.NewNumberKeyConcurrentMap[uint64, map[any]any](4096)

type RoutineLocal[T any] struct {
	// - 不能用空结构体，空结构体地址是一样，多个local会变成同一个key
	_ byte
}

func (rl *RoutineLocal[T]) Set(val T) {
	goId := GetGoId()
	table, ok := m.Load(goId)
	if !ok {
		table = make(map[any]any)
		m.Store(goId, table)
	}
	table[rl] = val
}

func (rl *RoutineLocal[T]) Get() (res T, exist bool) {
	goId := GetGoId()
	table, ok := m.Load(goId)
	if !ok {
		exist = false
		return
	}
	obj, exist := table[rl]
	if exist {
		res = obj.(T)
	}
	return
}

func (rl *RoutineLocal[T]) Clear() {
	goId := GetGoId()
	table, ok := m.Load(goId)
	if !ok {
		return
	}
	delete(table, rl)
	if len(table) == 0 {
		m.Delete(goId) // 释放map占用的内存
	}
}
