# T015: Define Provider Interface and Contract Tests

**Phase**: 2 - Local Source Provider  
**Category**: Provider Architecture  
**Complexity**: Medium  
**Estimated Duration**: 2-3 hours

## Objective

Define the `VersionProvider` interface contract and establish test patterns for provider implementations.

## Current State

- No provider interface or contract tests exist

## Target State

- `VersionProvider` interface is fully defined
- Contract tests ensure provider implementations are correct
- Error handling is consistent across providers
- Caching hints are part of the interface

## Acceptance Criteria

- [ ] `internal/providers/provider.go` defines:
  - `VersionProvider` interface with GetCurrent, GetLatest, List methods
  - `VersionResult` struct with version, source, timestamp, raw data
  - `QueryOptions` struct for filtering (core, stage, limit, etc.)
- [ ] Contract tests validate:
  - GetCurrent returns exactly one version
  - GetLatest returns exactly one version (highest)
  - List returns all versions in precedence order
  - Errors are properly categorized
- [ ] Provider implementations can opt-in to caching
- [ ] All providers log meaningful diagnostic information

## Context

### Files to Create

- `internal/providers/provider.go` — interface definitions
- `internal/providers/contract_test.go` — contract tests

### Interface Design

```go
type VersionProvider interface {
    // Name returns provider identifier (e.g., "git-tags", "github-releases")
    Name() string
    
    // GetCurrent returns the current/latest version
    GetCurrent(ctx context.Context, opts QueryOptions) (VersionResult, error)
    
    // GetLatest returns the highest available version matching constraints
    GetLatest(ctx context.Context, opts QueryOptions) (VersionResult, error)
    
    // List returns all available versions (unfiltered, for direct access)
    List(ctx context.Context, opts QueryOptions) ([]VersionResult, error)
}

type VersionResult struct {
    Version        *version.Version
    RawTag         string
    Source         string // provider name
    Ecosystem      string // detected ecosystem
    Timestamp      time.Time
    URL            string // reference link (optional)
    Digest         string // for containers (optional)
}

type QueryOptions struct {
    // Filtering
    Core          string // e.g., "1.2.3" to get versions of 1.2.3.*
    Stage         version.Stage
    IncludePrerelease bool
    
    // Pagination
    Limit         int
    Offset        int
    
    // Source-specific
    Raw           map[string]interface{}
}
```

### Contract Test Pattern

```go
// Contract test factory
func runProviderContractTests(t *testing.T, provider VersionProvider) {
    t.Run("GetCurrent", func(t *testing.T) {
        result, err := provider.GetCurrent(context.Background(), QueryOptions{})
        if err != nil {
            t.Fatalf("GetCurrent error: %v", err)
        }
        if result.Version == nil {
            t.Fatal("GetCurrent returned nil version")
        }
    })
    
    t.Run("GetLatest", func(t *testing.T) {
        result, err := provider.GetLatest(context.Background(), QueryOptions{})
        if err != nil {
            t.Fatalf("GetLatest error: %v", err)
        }
        if result.Version == nil {
            t.Fatal("GetLatest returned nil version")
        }
    })
    
    t.Run("List", func(t *testing.T) {
        results, err := provider.List(context.Background(), QueryOptions{Limit: 10})
        if err != nil {
            t.Fatalf("List error: %v", err)
        }
        if len(results) == 0 {
            t.Fatal("List returned no results")
        }
    })
}
```

## Testing

- [ ] Unit test: Interface methods are callable
- [ ] Contract test: All provider implementations pass contract
- [ ] Error handling: Test network errors, timeouts, authentication failures

## Related Tickets

- T016: Git-tags provider (implements interface)
- T023: GitHub Releases provider (implements interface)
- T024: GHCR provider (implements interface)

## Notes

- Contract tests should be reusable for all provider implementations
- Consider adding optional interfaces for caching/performance hints
