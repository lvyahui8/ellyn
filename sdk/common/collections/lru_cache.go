package collections

import (
	"github.com/lvyahui8/ellyn/sdk/common/definitions"
	"sync"
)

type cacheItem[K comparable, V definitions.Recyclable] struct {
	key  K
	val  V
	prev *cacheItem[K, V]
	next *cacheItem[K, V]
}

type LRUCache[K comparable, V definitions.Recyclable] struct {
	lock     sync.RWMutex
	capacity int
	head     *cacheItem[K, V]
	tail     *cacheItem[K, V]
	table    map[K]*cacheItem[K, V]
}

func NewLRUCache[K comparable, V definitions.Recyclable](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		table:    make(map[K]*cacheItem[K, V]),
	}
}

func (cache *LRUCache[K, V]) Get(key K) (V, bool) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	return cache.get(key)
}

func (cache *LRUCache[K, V]) get(key K) (v V, ok bool) {
	item := cache.table[key]
	if item == nil {
		return
	}
	cache.moveToFirst(item)
	return item.val, true
}

func (cache *LRUCache[K, V]) GetWithDefault(key K, createDefault func() V) V {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	val, exist := cache.get(key)
	if !exist {
		val = createDefault()
		cache.set(key, val)
	}
	return val
}

func (cache *LRUCache[K, V]) Set(key K, value V) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.set(key, value)
}

func (cache *LRUCache[K, V]) set(key K, value V) {
	item := cache.table[key]
	if item == nil {
		if len(cache.table) >= cache.capacity {
			// 空间满了，清理最久未使用的元素
			cache.removeItem(cache.tail)
		}
		item = &cacheItem[K, V]{
			val: value,
			key: key,
		}
		cache.table[item.key] = item
	}
	cache.moveToFirst(item)
}

func (cache *LRUCache[K, V]) Remove(key K) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	item := cache.table[key]
	cache.removeItem(item)
}

func (cache *LRUCache[K, V]) Values() (res []V) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	for cur := cache.head; cur != nil; {
		res = append(res, cur.val)
		cur = cur.next
	}
	return
}

func (cache *LRUCache[K, V]) moveToFirst(item *cacheItem[K, V]) {
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

func (cache *LRUCache[K, V]) removeItemFromLink(item *cacheItem[K, V]) {
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

func (cache *LRUCache[K, V]) removeItem(item *cacheItem[K, V]) {
	cache.removeItemFromLink(item)
	item.val.Recycle()
	delete(cache.table, item.key)
}
