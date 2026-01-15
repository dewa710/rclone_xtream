package xtream

import (
	"sync"
	"time"
)

type cacheItem struct {
	value any
	exp   time.Time
}

type Cache struct {
	mu    sync.RWMutex
	ttl   time.Duration
	items map[string]cacheItem
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl:   ttl,
		items: map[string]cacheItem{},
	}
}

func (c *Cache) Get(k string) (any, bool) {
	c.mu.RLock()
	i, ok := c.items[k]
	c.mu.RUnlock()
	if !ok || time.Now().After(i.exp) {
		return nil, false
	}
	return i.value, true
}

func (c *Cache) Set(k string, v any) {
	c.mu.Lock()
	c.items[k] = cacheItem{v, time.Now().Add(c.ttl)}
	c.mu.Unlock()
}
