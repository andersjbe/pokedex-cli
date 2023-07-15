package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mux *sync.Mutex
}

func NewCache(interval time.Duration, mux *sync.Mutex) Cache {
	c := Cache {
		cache: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	newEntry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	// c.mux.Lock()
	c.cache[key] = newEntry
	// c.mux.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	// c.mux.Lock()
	entry, ok := c.cache[key]
	// c.mux.Unlock()
	if (!ok) {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reap(interval time.Duration) {
	// c. mux.Lock()
	for key, entry := range c.cache {
		if time.Now().Sub(entry.createdAt) > interval {
			delete(c.cache, key)
		}
	}
	// c.mux.Unlock()
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}
