# T029: Implement Optional Policy Checks

**Phase**: 4 - CI/Release Enhancements  
**Category**: Release Features  
**Complexity**: Medium  
**Estimated Duration**: 2-3 hours

## Objective

Implement lightweight policy checks for release validation (optional, can be minimal).

## Current State

- No policy checks exist

## Target State

- Validate version against simple policies
- Support branch-specific rules (optional)
- Integration with version bump / current commands
- Clear error messages for policy violations

## Acceptance Criteria

- [ ] `internal/version/policy.go` defines policy checks
- [ ] Validate version >= previous version (for releases)
- [ ] Optional: enforce specific stage on branches (main=final only, develop=any)
- [ ] Fail fast with clear error messages
- [ ] Configuration via `.verctl.yaml` (optional)
- [ ] Integration: `version bump` can validate result against policy
- [ ] All checks are optional (can be disabled)

## Context

### Files to Create

- `internal/version/policy.go`
- `internal/version/policy_test.go`

### Configuration (Optional)

```yaml
rules:
  # Optional policy checks
  branchPolicies:
    main: final        # only final releases on main
    develop: any       # any version on develop
  minorVersion: true   # enforce minor.patch increment
```

### Implementation Pattern

```go
type PolicyChecker struct {
    rules *config.Rules
}

func (p *PolicyChecker) Validate(version *Version, currentVersion *Version) error {
    // Check monotonic increase
    if p.rules.MonotonicIncrease && currentVersion != nil {
        if comparator.Compare(version, currentVersion) <= 0 {
            return fmt.Errorf("new version must be > current version")
        }
    }
    
    // Branch-specific checks (if branch info available)
    // ...
    
    return nil
}
```

## Testing

- [ ] Unit test: Validate monotonic increase
- [ ] Unit test: Skip checks when disabled
- [ ] Error test: Clear error messages

## Related Tickets

- T020: `version bump` command (uses policy)

## Notes

- Keep policies simple for v1; complex DSL deferred to v2
- All checks are optional and can be disabled
