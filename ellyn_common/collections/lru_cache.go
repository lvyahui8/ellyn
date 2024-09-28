package collections

import (
	"sync"
)

type cacheItem[K comparable] struct {
	key  K
	val  any
	prev *cacheItem[K]
	next *cacheItem[K]
}

type LRUCache[K comparable] struct {
	lock     sync.RWMutex
	capacity int
	head     *cacheItem[K]
	tail     *cacheItem[K]
	table    map[K]*cacheItem[K]
}

func NewLRUCache[K comparable](capacity int) *LRUCache[K] {
	return &LRUCache[K]{
		capacity: capacity,
		table:    make(map[K]*cacheItem[K]),
	}
}

func (cache *LRUCache[K]) Get(key K) (any, bool) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	item := cache.table[key]
	if item == nil {
		return nil, false
	}
	cache.moveToFirst(item)
	return item.val, true
}

func (cache *LRUCache[K]) Set(key K, value any) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	item := cache.table[key]
	if item == nil {
		if len(cache.table) >= cache.capacity {
			cache.removeItem(cache.tail)
		}
		item = &cacheItem[K]{
			val: value,
			key: key,
		}
		cache.table[item.key] = item
	}
	cache.moveToFirst(item)
}

func (cache *LRUCache[K]) Remove(key K) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	item := cache.table[key]
	cache.removeItem(item)
}

func (cache *LRUCache[K]) Values() (res []any) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	for cur := cache.head; cur != nil; {
		res = append(res, cur.val)
		cur = cur.next
	}
	return
}

func (cache *LRUCache[K]) moveToFirst(item *cacheItem[K]) {
	if cache.head != nil && cache.head.key == item.key {
		return
	}
	// 将元素从原链表摘出
	cache.removeItemFromLink(item)

	// 更新到链表头部
	if cache.head == nil {
		cache.tail = item
	} else {
		cache.head.prev = item
	}
	item.prev = nil
	item.next = cache.head
	cache.head = item
}

func (cache *LRUCache[K]) removeItemFromLink(item *cacheItem[K]) {
	if item == nil {
		return
	}
	// 将元素从原链表摘出
	if cache.tail != nil && cache.tail.key == item.key {
		cache.tail = item.prev
	}
	if item.prev != nil {
		item.prev.next = item.next
	}
	if item.next != nil {
		item.next.prev = item.prev
	}
}

func (cache *LRUCache[K]) removeItem(item *cacheItem[K]) {
	cache.removeItemFromLink(item)
	delete(cache.table, item.key)
}
