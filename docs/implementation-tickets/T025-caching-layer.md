# T025: Implement Lightweight Caching Layer

**Phase**: 3 - Remote Providers  
**Category**: Performance  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Implement a lightweight caching layer with TTL and conditional request support to reduce API calls.

## Current State

- Individual providers implement basic caching

## Target State

- Centralized cache with pluggable backends (in-memory, optional file)
- TTL-based expiration
- Conditional request support (ETag, Last-Modified)
- Cache eviction policies

## Acceptance Criteria

- [ ] `internal/providers/cache.go` defines cache interface
- [ ] In-memory cache implementation with TTL
- [ ] Optional file-based cache for persistence
- [ ] ETag support for HTTP conditional requests
- [ ] Cache key based on provider + query options
- [ ] Configurable TTL per provider
- [ ] Cache statistics (hits, misses, size)
- [ ] No secrets stored in cache

## Context

### Files to Create

- `internal/providers/cache.go` — cache interface and implementation
- `internal/providers/cache_test.go` — tests

### Cache Interface Design

```go
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
    Clear()
    Stats() CacheStats
}

type CacheStats struct {
    Hits   int
    Misses int
    Size   int
}

// Example usage in provider
func (g *githubProvider) GetLatest(ctx context.Context, opts QueryOptions) (VersionResult, error) {
    cacheKey := fmt.Sprintf("gh:latest:%s", opts.Hash())
    
    if cached, ok := cache.Get(cacheKey); ok {
        return cached.(VersionResult), nil
    }
    
    result, err := g.fetchLatest(ctx, opts)
    if err == nil {
        cache.Set(cacheKey, result, g.cacheTTL)
    }
    return result, err
}
```

## Testing

- [ ] Unit test: Cache hit/miss behavior
- [ ] Unit test: TTL expiration
- [ ] Unit test: ETag support
- [ ] Unit test: Cache statistics
- [ ] Performance test: Cache reduces API calls

## Related Tickets

- T016: Git-tags (uses caching)
- T023: GitHub Releases (uses caching)
- T024: GHCR (uses caching)

## Notes

- Keep cache implementation simple and fast
- Consider memory limits for in-memory cache
