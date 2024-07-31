package gls

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
)

var m = collections.NewConcurrentMap(2048, func(key any) int {
	return int(key.(uint64))
})

type RoutineLocal[T any] struct {
}

func (rl *RoutineLocal[T]) Set(val T) {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	var table map[any]any
	if !ok {
		table = make(map[any]any)
		m.Store(goId, table)
	} else {
		table = tableVal.(map[any]any)
	}
	table[rl] = val
}

func (rl *RoutineLocal[T]) Get() (res T, exist bool) {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	if !ok {
		exist = false
		return
	}
	table := tableVal.(map[any]any)
	obj, exist := table[rl]
	if exist {
		res = obj.(T)
	}
	return
}

func (rl *RoutineLocal[T]) Clear() {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	if !ok {
		return
	}
	table := tableVal.(map[any]any)
	delete(table, rl)
}
