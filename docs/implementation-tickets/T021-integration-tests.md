# T021: Integration Tests for Provider and Command Orchestration

**Phase**: 2 - Local Source Provider  
**Category**: Testing  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Establish integration tests that validate provider behavior, command orchestration, and multi-source precedence.

## Current State

- No integration tests exist

## Target State

- Full end-to-end tests with fake providers
- Multi-source precedence validation
- Command output validation
- Real git repository test case

## Acceptance Criteria

- [ ] `tests/integration/provider_test.go` with:
  - Fake provider implementations for testing
  - Multi-source provider chain tests
  - Source precedence validation
  - Conflict resolution tests
- [ ] `tests/integration/commands_test.go` with:
  - `version current` full workflow
  - `version latest` with constraints
  - `version bump` integration
  - Error handling tests
- [ ] Real git repository test case:
  - Create temporary git repo with tags
  - Test provider against real repo
  - Cleanup after test
- [ ] Output validation:
  - Text output matches expected format
  - JSON output is valid and complete
  - Error messages are actionable

## Context

### Files to Create

- `tests/integration/provider_test.go`
- `tests/integration/commands_test.go`
- `tests/fixtures/fake_provider.go` — mock provider for testing
- `tests/fixtures/test_repo.go` — temporary git repository setup

### Fake Provider Pattern

```go
type fakeProvider struct {
    versions []*VersionResult
    name     string
}

func (f *fakeProvider) GetCurrent(ctx context.Context, opts QueryOptions) (VersionResult, error) {
    // Return highest
}

func (f *fakeProvider) GetLatest(ctx context.Context, opts QueryOptions) (VersionResult, error) {
    // Apply filtering
}

func (f *fakeProvider) List(ctx context.Context, opts QueryOptions) ([]VersionResult, error) {
    return f.versions, nil
}
```

### Multi-Source Test Pattern

```go
func TestMultiSourcePrecedence(t *testing.T) {
    provider1 := fakeProvider{versions: []VersionResult{
        {Version: Version{1, 2, 3, Final, nil}, Source: "git-tags"},
    }}
    provider2 := fakeProvider{versions: []VersionResult{
        {Version: Version{1, 2, 4, Final, nil}, Source: "github-releases"},
    }}
    
    chain := NewProviderChain([]Provider{provider1, provider2})
    result, _ := chain.GetLatest(context.Background(), QueryOptions{})
    
    // Should select provider1 (higher precedence), but provider1 has 1.2.3
    // So should return 1.2.4 from provider2 if precedence allows
}
```

## Testing

- [ ] Integration test: Complete workflow from config to output
- [ ] Integration test: Multi-source provider chain
- [ ] Integration test: Provider precedence is respected
- [ ] Integration test: Error handling across provider chain
- [ ] Real repo test: Create temporary git repo and test provider
- [ ] Output test: Text and JSON formats are correct

## Related Tickets

- T015: Provider interface
- T016: Git-tags provider
- T018-T020: Commands

## Notes

- Use testify/assert for cleaner test code
- Create temporary directories for git repo tests; clean up after
- Test both success and failure paths
