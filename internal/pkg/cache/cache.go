package cache

import (
	"sync"
	"time"
)

type item struct {
	value     any
	expiresAt time.Time
}

// Cache is a simple TTL cache replacing go-zero collection.Cache.
type Cache struct {
	mu    sync.RWMutex
	ttl   time.Duration
	limit int
	items map[string]item
}

func New(ttl time.Duration, limit int) (*Cache, error) {
	return &Cache{
		ttl:   ttl,
		limit: limit,
		items: make(map[string]item),
	}, nil
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(it.expiresAt) {
		if ok {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}
		return nil, false
	}

	return it.value, true
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.limit > 0 && len(c.items) >= c.limit {
		for k := range c.items {
			delete(c.items, k)
			break
		}
	}

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}
