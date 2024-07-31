package gls

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
)

var m = collections.NewConcurrentMap(2048, func(key interface{}) int {
	return int(key.(uint64))
})

type RoutineLocal struct {
}

func (rl *RoutineLocal) Set(val interface{}) {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	var table map[interface{}]interface{}
	if !ok {
		table = make(map[interface{}]interface{})
		m.Store(goId, table)
	} else {
		table = tableVal.(map[interface{}]interface{})
	}
	table[rl] = val
}

func (rl *RoutineLocal) Get() (interface{}, bool) {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	if !ok {
		return nil, false
	}
	table := tableVal.(map[interface{}]interface{})
	res, ok := table[rl]
	return res, ok
}

func (rl *RoutineLocal) Clear() {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	if !ok {
		return
	}
	table := tableVal.(map[interface{}]interface{})
	delete(table, rl)
}
