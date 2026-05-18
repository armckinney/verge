# T024: Implement GHCR Provider

**Phase**: 3 - Remote Providers  
**Category**: Provider  
**Complexity**: High  
**Estimated Duration**: 4-5 hours

## Objective

Implement the GitHub Container Registry (GHCR) provider to fetch versions from container image tags.

## Current State

- No GHCR provider exists

## Target State

- Fetches tags from GHCR image
- Parses container versioning patterns (PR hashes, REL versions, floating)
- Supports channel filtering (rel, pr, floating)
- Respects rate limits

## Acceptance Criteria

- [ ] `internal/providers/ghcr.go` implements `VersionProvider`
- [ ] GetCurrent returns latest image tag
- [ ] GetLatest returns highest version (respecting channel filter)
- [ ] List returns all tags sorted
- [ ] Channel filtering: `rel` (final only), `pr` (prerelease only), `floating` (1.2, 1)
- [ ] Caching with efficient tag listing
- [ ] Configuration: image path, authentication (optional)
- [ ] Pattern detection: recognizes PR hashes, release versions

## Context

### Files to Create

- `internal/providers/ghcr.go`
- `internal/providers/ghcr_test.go`

### Configuration

```yaml
sources:
  ghcr:
    enabled: false
    image: ghcr.io/your-org/your-image
    includePrerelease: true
    channelFilter: null  # rel | pr | floating | null
```

### Pattern Recognition

- Release: `1.2.3` (matches final version pattern)
- PR: `1.2.3-dev.a1b2c3d` (matches prerelease with hash)
- Floating: `1.2`, `1` (major.minor or major only)
- Special: `latest` (mapped to current)

## Testing

- [ ] Unit test: Parse various container tag formats
- [ ] Unit test: Channel filtering works
- [ ] Integration test: Real GHCR API call (mock)
- [ ] Error test: 401, 404, rate limit handling

## Related Tickets

- T015: Provider interface
- T023: GitHub Releases (pattern)

## Notes

- Consider lazy loading of large tag lists
- Support both public and private images (auth via token)
