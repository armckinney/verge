package providers

import (
	"sync"
	"time"
)

// CacheStats holds cache performance counters.
type CacheStats struct {
	Hits   int
	Misses int
	Size   int
}

// Cache is a generic TTL key-value cache.
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
	Stats() CacheStats
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

type memoryCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
	hits    int
	misses  int
}

// NewMemoryCache returns a new in-memory cache.
func NewMemoryCache() Cache {
	return &memoryCache{
		entries: make(map[string]cacheEntry),
	}
}

func (c *memoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		// Evict expired entry if present
		if ok {
			delete(c.entries, key)
		}
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.Lock()
	c.hits++
	c.mu.Unlock()
	return entry.value, true
}

func (c *memoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *memoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

func (c *memoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]cacheEntry)
}

func (c *memoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Count only non-expired entries
	size := 0
	now := time.Now()
	for _, e := range c.entries {
		if !now.After(e.expiresAt) {
			size++
		}
	}
	return CacheStats{
		Hits:   c.hits,
		Misses: c.misses,
		Size:   size,
	}
}
