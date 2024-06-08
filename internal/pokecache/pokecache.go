package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	counts map[string]CacheEntry
	mu     *sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		counts: make(map[string]CacheEntry),
		mu:     &sync.Mutex{},
	}

	go cache.reapLoop(interval)

	return cache
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.counts[key] = entry
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.counts[key]

	if !ok {
		return nil, false
	}

	return v.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.counts {
			if time.Since(entry.createdAt) > interval {
				delete(c.counts, key)
			}
		}
		c.mu.Unlock()
	}
}
