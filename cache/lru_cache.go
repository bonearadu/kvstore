package cache

import (
	"container/list"
	"sync"
)

type entry struct {
	key string
	val string
}

type LRUCache struct {
	store    *list.List
	elements map[string]*list.Element
	capacity int
	mu       sync.RWMutex
}

func (c *LRUCache) evict() {
	for len(c.elements) > c.capacity {
		lru := c.store.Back()
		e := c.store.Remove(lru)
		delete(c.elements, e.(*entry).key)
	}
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		elements: make(map[string]*list.Element),
		store:    list.New(),
		capacity: capacity,
		mu:       sync.RWMutex{},
	}
}

func (c *LRUCache) Read(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.elements[key]

	if ok {
		return val.Value.(*entry).val, ok
	}

	return "", false
}

func (c *LRUCache) Write(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.elements[key]
	if ok {
		c.store.Remove(val)
	}

	elem := c.store.PushFront(&entry{key, value})
	c.elements[key] = elem

	if len(c.elements) > c.capacity {
		c.evict()
	}
}

func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.elements[key]
	if ok {
		delete(c.elements, key)
		c.store.Remove(e)
	}
}
