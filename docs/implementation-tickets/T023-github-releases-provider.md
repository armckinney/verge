# T023: Implement GitHub Releases Provider

**Phase**: 3 - Remote Providers  
**Category**: Provider  
**Complexity**: High  
**Estimated Duration**: 4-5 hours

## Objective

Implement the GitHub Releases provider to fetch versions from GitHub Releases API.

## Current State

- No GitHub provider exists

## Target State

- Fetches releases (and prereleases) from GitHub API
- Supports token-based authentication
- Efficient caching with conditional requests
- Respects rate limits

## Acceptance Criteria

- [ ] `internal/providers/github_releases.go` implements `VersionProvider`
- [ ] GetCurrent returns latest release
- [ ] GetLatest returns highest version (optionally including prereleases)
- [ ] List returns all releases sorted
- [ ] Token auth via `GITHUB_TOKEN` env var
- [ ] Caching with ETag support
- [ ] Error handling: 401 (auth), 403 (rate limit), 404 (not found), network errors
- [ ] Configuration: owner, repo, include drafts, prerelease handling
- [ ] Performance: Uses GraphQL or REST API efficiently

## Context

### Files to Create

- `internal/providers/github_releases.go`
- `internal/providers/github_releases_test.go`

### Configuration

```yaml
sources:
  github-releases:
    enabled: false
    owner: your-org
    repo: your-repo
    includePrerelease: true
    includeDrafts: false
```

### Implementation Notes

- Use REST API for simplicity (or GraphQL for efficiency later)
- Endpoints: `/repos/{owner}/{repo}/releases`
- Parse release tag_name as version string
- Use `X-Etag` for caching (conditional requests)
- Handle rate limiting: retry-after headers

## Testing

- [ ] Unit test: Mock API responses
- [ ] Integration test: Real API call with auth (optional)
- [ ] Error test: 401, 403, 404 handling
- [ ] Cache test: ETag-based caching works
- [ ] Prerelease test: Filtering works correctly

## Related Tickets

- T015: Provider interface
- T016: Git-tags provider (pattern)
- T025: Caching layer (shared)

## Notes

- Token secrets should never be logged
- Consider adding organization-level fallback
