package gls

import (
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"sync/atomic"
)

var m = collections.NewConcurrentMap(2048, func(key interface{}) int {
	return int(key.(uint64))
})

var localId int64 = 0

type routineLocal struct {
	id int64
}

func NewRoutineLocal() *routineLocal {
	return &routineLocal{
		id: atomic.AddInt64(&localId, 1),
	}
}

func (rl *routineLocal) Set(val interface{}) {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	var table map[int64]interface{}
	if !ok {
		table = make(map[int64]interface{})
		m.Store(goId, table)
	} else {
		table = tableVal.(map[int64]interface{})
	}
	table[rl.id] = val
}

func (rl *routineLocal) Get() (interface{}, bool) {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	if !ok {
		return nil, false
	}
	table := tableVal.(map[int64]interface{})
	res, ok := table[rl.id]
	return res, ok
}

func (rl *routineLocal) Clear() {
	goId := GetGoId()
	tableVal, ok := m.Load(goId)
	if !ok {
		return
	}
	table := tableVal.(map[int64]interface{})
	delete(table, rl.id)
}
