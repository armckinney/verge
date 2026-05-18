# T016: Implement Git-Tags Provider

**Phase**: 2 - Local Source Provider  
**Category**: Provider  
**Complexity**: High  
**Estimated Duration**: 4-5 hours

## Objective

Implement the `git-tags` provider to fetch versions from local or remote git repositories.

## Current State

- No git provider exists

## Target State

- Git provider fetches tags from local/remote repositories
- Ecosystem-aware parsing (Go, Python, containers, etc.)
- Tag prefix filtering (e.g., `v` prefix)
- Prerelease filtering
- Efficient caching of tag list

## Acceptance Criteria

- [ ] `internal/providers/git_tags.go` implements `VersionProvider`
- [ ] GetCurrent returns the highest git tag (parsed as version)
- [ ] GetLatest returns highest version matching constraints
- [ ] List returns all tags sorted by version
- [ ] Tag prefix filtering: `v1.2.3` matches when prefix is `v`
- [ ] Supports local git repository (`.git` detection)
- [ ] Supports remote repository (optional, URL-based)
- [ ] Prerelease filtering works correctly
- [ ] Error handling: no git repo, no tags, invalid tags
- [ ] Performance: caches tag list to avoid repeated git calls
- [ ] Respects `fetch: false` config (local tags only)

## Context

### Files to Create

- `internal/providers/git_tags.go` — git provider
- `internal/providers/git_tags_test.go` — tests

### Implementation Strategy

```go
type gitTagsProvider struct {
    repoPath      string
    tagPrefix     string
    fetchRemote   bool
    cache         []string
    cacheTime     time.Time
    cacheTTL      time.Duration
}

func (g *gitTagsProvider) GetCurrent(ctx context.Context, opts QueryOptions) (VersionResult, error) {
    tags, err := g.getTags(ctx)
    if err != nil {
        return VersionResult{}, err
    }
    
    // Filter and sort
    versions := g.parseAndSort(tags, opts)
    if len(versions) == 0 {
        return VersionResult{}, errors.New("no matching versions found")
    }
    
    return versions[0], nil
}

func (g *gitTagsProvider) getTags(ctx context.Context) ([]string, error) {
    // Check cache
    if time.Since(g.cacheTime) < g.cacheTTL && len(g.cache) > 0 {
        return g.cache, nil
    }
    
    // Run: git tag -l
    cmd := exec.CommandContext(ctx, "git", "tag", "-l")
    cmd.Dir = g.repoPath
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("git tag failed: %w", err)
    }
    
    tags := strings.Split(strings.TrimSpace(string(output)), "\n")
    g.cache = tags
    g.cacheTime = time.Now()
    
    return tags, nil
}

func (g *gitTagsProvider) parseAndSort(tags []string, opts QueryOptions) []VersionResult {
    var results []VersionResult
    for _, tag := range tags {
        // Strip prefix
        versionStr := strings.TrimPrefix(tag, g.tagPrefix)
        
        // Parse
        parsed, err := version.DefaultParser.Parse(versionStr)
        if err != nil {
            continue // Skip unparseable tags
        }
        
        // Filter by stage
        if opts.IncludePrerelease == false && parsed.Stage != version.StageFinal {
            continue
        }
        
        results = append(results, VersionResult{
            Version:   parsed,
            RawTag:    tag,
            Source:    g.Name(),
            Timestamp: time.Now(),
        })
    }
    
    // Sort by version (descending)
    sort.Slice(results, func(i, j int) bool {
        return version.DefaultComparator.Compare(results[i].Version, results[j].Version) > 0
    })
    
    return results
}
```

### Configuration

From `.verge.yaml`:

```yaml
sources:
  git-tags:
    enabled: true
    fetch: false              # local tags only
    includePrerelease: true
    ecosystemParsing: go
```

### Test Strategy

- Test with real git repository (fibonacci of tags)
- Test with local tags only (`fetch: false`)
- Test with invalid tags mixed in
- Test error cases: no git repo, no tags
- Test caching: verify cache is used and TTL respected

## Testing

- [ ] Unit test: Parse various git tag formats
- [ ] Unit test: Filter by prerelease
- [ ] Integration test: GetCurrent with real git repo
- [ ] Integration test: GetLatest with constraints
- [ ] Integration test: Caching works and TTL is respected
- [ ] Error test: No git repo available

## Related Tickets

- T015: Provider interface
- T017: Sequence interpreter (used by provider)
- T018: version current command (uses provider)
- T019: version latest command (uses provider)

## Notes

- Handle large tag lists efficiently (cache, pagination)
- Support both v-prefixed and unprefixed tags
- Consider lazy evaluation for large repos
