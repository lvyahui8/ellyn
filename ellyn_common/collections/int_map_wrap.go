package collections

import (
	"math"
)

// IntMapWrap 用于int map并发只读的情况，通过bitmap减少map无效读取的损耗
type IntMapWrap[T any] struct {
	data   map[int]T
	flags  *BitMap
	maxKey int
	minKey int
}

func NewIntMapWrap[T any](data map[int]T) *IntMapWrap[T] {
	m := &IntMapWrap[T]{}
	m.data = data
	maxKey := math.MinInt
	minKey := math.MaxInt
	for key := range data {
		if key > maxKey {
			maxKey = key
		}
		if key < minKey {
			minKey = key
		}
	}
	m.maxKey = maxKey
	m.minKey = minKey
	m.flags = NewBitMap(uint(maxKey - minKey + 1))
	for key := range data {
		m.flags.Set(uint(key - minKey))
	}
	return m
}

func (m *IntMapWrap[T]) Get(key int) (val T, exist bool) {
	if key > m.maxKey || key < m.minKey {
		return
	}
	if !m.flags.Get(uint(key - m.minKey)) {
		exist = false
		return
	}

	val, exist = m.data[key]
	return
}
