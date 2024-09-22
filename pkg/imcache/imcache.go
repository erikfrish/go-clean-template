package imcache

import (
	"sync"
)

type imCache[KeyT comparable, ValueT any] struct {
	mu sync.RWMutex
	m  map[KeyT]ValueT
}

func New[KeyT comparable, ValueT any]() *imCache[KeyT, ValueT] {
	return &imCache[KeyT, ValueT]{
		sync.RWMutex{},
		make(map[KeyT]ValueT),
	}
}

func (c *imCache[KeyT, ValueT]) Set(key KeyT, value ValueT) {
	c.mu.Lock()
	c.m[key] = value
	c.mu.Unlock()
}

func (c *imCache[KeyT, ValueT]) Swap(key KeyT, val ValueT) (ValueT, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	old, ok := c.m[key]
	c.m[key] = val
	return old, ok
}

func (c *imCache[KeyT, ValueT]) Get(key KeyT) (ValueT, bool) {
	c.mu.RLock()
	value, ok := c.m[key]
	c.mu.RUnlock()
	return value, ok
}

func (c *imCache[KeyT, ValueT]) Delete(key KeyT) {
	c.mu.Lock()
	delete(c.m, key)
	c.mu.Unlock()
}

func (c *imCache[KeyT, ValueT]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.m)
}

func (c *imCache[KeyT, ValueT]) Clear() {
	c.mu.Lock()
	c.m = make(map[KeyT]ValueT)
	c.mu.Unlock()
}

func (c *imCache[KeyT, ValueT]) Keys() []KeyT {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]KeyT, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	return keys
}

func (c *imCache[KeyT, ValueT]) Values() []ValueT {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]ValueT, 0, len(c.m))
	for _, v := range c.m {
		values = append(values, v)
	}
	return values
}

func (c *imCache[KeyT, ValueT]) GetMap() map[KeyT]ValueT {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items := make(map[KeyT]ValueT, len(c.m))
	for k, v := range c.m {
		items[k] = v
	}
	return items
}

func (c *imCache[KeyT, ValueT]) SetMap(items map[KeyT]ValueT) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m = items
}
