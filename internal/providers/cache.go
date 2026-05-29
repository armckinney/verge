package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"example.com/verge/internal/version"
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

// DiskCacheEntry represents a single serialized entry in the disk cache.
type DiskCacheEntry struct {
	Value     json.RawMessage `json:"value"`
	ExpiresAt time.Time       `json:"expires_at"`
}

// DiskCache represents a thread-safe persistent file-based JSON cache.
type DiskCache struct {
	mu       sync.RWMutex
	filePath string
	entries  map[string]DiskCacheEntry
}

// NewDiskCache returns a new persistent DiskCache, using the default user cache directory if path is empty.
func NewDiskCache(filePath string) (*DiskCache, error) {
	if filePath == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			cacheDir = os.TempDir()
		}
		filePath = filepath.Join(cacheDir, "verge", "cache.json")
	}

	c := &DiskCache{
		filePath: filePath,
		entries:  make(map[string]DiskCacheEntry),
	}

	if err := c.load(); err != nil {
		c.entries = make(map[string]DiskCacheEntry)
	}

	return c, nil
}

func (c *DiskCache) load() error {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var fileData struct {
		Entries map[string]DiskCacheEntry `json:"entries"`
	}
	if err := json.Unmarshal(data, &fileData); err != nil {
		return err
	}

	c.entries = fileData.Entries
	if c.entries == nil {
		c.entries = make(map[string]DiskCacheEntry)
	}
	return nil
}

func (c *DiskCache) save() error {
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fileData := struct {
		Entries map[string]DiskCacheEntry `json:"entries"`
	}{
		Entries: c.entries,
	}

	data, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := c.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, c.filePath)
}

// Get retrieves an entry from the disk cache and unmarshals it into target.
func (c *DiskCache) Get(key string, target interface{}) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return false, nil
	}

	if time.Now().After(entry.ExpiresAt) {
		return false, nil
	}

	if err := json.Unmarshal(entry.Value, target); err != nil {
		return false, fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	return true, nil
}

// Set serializes value and saves it under key with a TTL.
func (c *DiskCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	c.entries[key] = DiskCacheEntry{
		Value:     raw,
		ExpiresAt: time.Now().Add(ttl),
	}

	return c.save()
}

// Delete removes an entry from the disk cache.
func (c *DiskCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
	return c.save()
}

// Clear clears all entries in the disk cache.
func (c *DiskCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]DiskCacheEntry)
	return c.save()
}

// CachingProvider wraps a VersionProvider with transparent disk caching.
type CachingProvider struct {
	underlying VersionProvider
	cacheKey   string
	disabled   bool
	diskCache  *DiskCache
	ttl        time.Duration
}

// NewCachingProvider wraps a provider with caching.
func NewCachingProvider(underlying VersionProvider, cacheKey string, disabled bool, ttl time.Duration) VersionProvider {
	dc, _ := NewDiskCache("")
	return &CachingProvider{
		underlying: underlying,
		cacheKey:   cacheKey,
		disabled:   disabled,
		diskCache:  dc,
		ttl:        ttl,
	}
}

// NewCachingProviderWithCache wraps a provider with a specific DiskCache instance.
func NewCachingProviderWithCache(underlying VersionProvider, cacheKey string, disabled bool, ttl time.Duration, dc *DiskCache) VersionProvider {
	return &CachingProvider{
		underlying: underlying,
		cacheKey:   cacheKey,
		disabled:   disabled,
		diskCache:  dc,
		ttl:        ttl,
	}
}

// Name returns the name of the underlying provider.
func (cp *CachingProvider) Name() string {
	return cp.underlying.Name()
}

// GetLatest implements VersionProvider with caching.
func (cp *CachingProvider) GetLatest(versionType string) (*version.Version, error) {
	if cp.disabled || cp.diskCache == nil {
		return cp.underlying.GetLatest(versionType)
	}

	key := fmt.Sprintf("%s:%s:latest", cp.cacheKey, versionType)
	var cached version.Version
	found, err := cp.diskCache.Get(key, &cached)
	if err == nil && found {
		return &cached, nil
	}

	v, err := cp.underlying.GetLatest(versionType)
	if err != nil {
		return nil, err
	}

	_ = cp.diskCache.Set(key, v, cp.ttl)
	return v, nil
}

// GetLatestSpecific implements VersionProvider with caching.
func (cp *CachingProvider) GetLatestSpecific(versionType string, prefix string) (*version.Version, error) {
	if cp.disabled || cp.diskCache == nil {
		return cp.underlying.GetLatestSpecific(versionType, prefix)
	}

	key := fmt.Sprintf("%s:%s:%s:latest_specific", cp.cacheKey, versionType, prefix)
	var cached version.Version
	found, err := cp.diskCache.Get(key, &cached)
	if err == nil && found {
		return &cached, nil
	}

	v, err := cp.underlying.GetLatestSpecific(versionType, prefix)
	if err != nil {
		return nil, err
	}

	_ = cp.diskCache.Set(key, v, cp.ttl)
	return v, nil
}
