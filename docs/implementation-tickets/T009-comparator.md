# T009: Implement Version Comparator with Deterministic Ordering

**Phase**: 1 - Core Version Engine  
**Category**: Domain Logic  
**Complexity**: High  
**Estimated Duration**: 4-5 hours

## Objective

Implement robust version comparison with deterministic ordering across SemVer, PEP440, ecosystem-specific, and mixed sequence types.

## Current State

- No comparator exists

## Target State

- `Comparator` interface implements version ordering
- Ordering follows semantic versioning rules with PEP440 compatibility
- Hash sequences are compared lexicographically
- Numeric sequences are compared numerically
- Final releases sort higher than prereleases
- Consistent tie-breaking for multi-source scenarios

## Acceptance Criteria

- [ ] `internal/version/compare.go` implements `Comparator` interface
- [ ] `Compare(left, right *Version) int` returns -1, 0, or 1
- [ ] Ordering rules are followed (major, minor, patch, stage, sequence)
- [ ] Stage ordering is enforced: dev < alpha < beta < rc < final
- [ ] Numeric sequences compared numerically: `1 < 2 < 10`
- [ ] Hash sequences compared lexicographically: `"a1b2c3d" < "b2c3d4e"`
- [ ] Final versions sort higher than prerelease of same core: `1.2.3 > 1.2.3-rc.99`
- [ ] Mixed sequence types ordered correctly (numeric < hash by default, or configurable)
- [ ] Edge cases handled: zero sequences, very large numbers, case sensitivity
- [ ] Comprehensive table-driven tests cover all combinations

## Context

### Files to Create

- `internal/version/compare.go` — comparator implementation
- `internal/version/compare_test.go` — table-driven tests

### Comparison Algorithm

```go
type comparator struct{}

func (c *comparator) Compare(left, right *Version) int {
    // 1. Compare core (major, minor, patch)
    if left.Major != right.Major {
        return cmpInt(left.Major, right.Major)
    }
    if left.Minor != right.Minor {
        return cmpInt(left.Minor, right.Minor)
    }
    if left.Patch != right.Patch {
        return cmpInt(left.Patch, right.Patch)
    }
    
    // 2. Final release sorts higher than prerelease
    if (left.Stage == StageFinal) != (right.Stage == StageFinal) {
        if left.Stage == StageFinal {
            return 1
        }
        return -1
    }
    
    // 3. Compare stage (if both prerelease)
    if cmp := c.compareStage(left.Stage, right.Stage); cmp != 0 {
        return cmp
    }
    
    // 4. Compare sequence
    return c.compareSequence(left.Sequence, right.Sequence)
}

func cmpInt(a, b int) int {
    if a < b {
        return -1
    } else if a > b {
        return 1
    }
    return 0
}

func (c *comparator) compareStage(left, right Stage) int {
    stageOrder := map[Stage]int{
        StageDev:   0,
        StageAlpha: 1,
        StageBeta:  2,
        StageRC:    3,
    }
    lv := stageOrder[left]
    rv := stageOrder[right]
    return cmpInt(lv, rv)
}

func (c *comparator) compareSequence(left, right interface{}) int {
    // Handle numeric sequences
    if li, lok := toInt(left); lok {
        if ri, rok := toInt(right); rok {
            return cmpInt(li, ri)
        }
        // Numeric < hash (numeric releases sort lower)
        return -1
    }
    if ri, rok := toInt(right); rok {
        // Hash > numeric
        return 1
    }
    
    // Both are strings (hashes or build IDs)
    ls, _ := toString(left)
    rs, _ := toString(right)
    return strings.Compare(ls, rs)
}

func toInt(v interface{}) (int, bool) {
    switch i := v.(type) {
    case int:
        return i, true
    case int64:
        return int(i), true
    }
    return 0, false
}

func toString(v interface{}) (string, bool) {
    if s, ok := v.(string); ok {
        return s, true
    }
    return "", false
}
```

### Test Corpus (Table-Driven)

Comprehensive test cases covering:

```go
[]struct {
    left     string
    right    string
    expected int // -1, 0, 1
}{
    // Basic core version comparison
    {"1.0.0", "2.0.0", -1},
    {"2.0.0", "1.0.0", 1},
    {"1.0.0", "1.0.0", 0},
    {"1.2.3", "1.2.4", -1},
    {"1.2.3", "1.3.0", -1},
    {"1.2.3", "2.0.0", -1},
    
    // Prerelease vs release
    {"1.2.3", "1.2.3-rc.1", 1},
    {"1.2.3-rc.1", "1.2.3", -1},
    
    // Stage ordering
    {"1.2.3-dev.1", "1.2.3-alpha.1", -1},
    {"1.2.3-alpha.1", "1.2.3-beta.1", -1},
    {"1.2.3-beta.1", "1.2.3-rc.1", -1},
    {"1.2.3-rc.1", "1.2.3", -1},
    
    // Sequence comparison (numeric)
    {"1.2.3-dev.1", "1.2.3-dev.2", -1},
    {"1.2.3-dev.99", "1.2.3-dev.100", -1},
    
    // Sequence comparison (hash)
    {"1.2.3-dev.a1b2c3d", "1.2.3-dev.b2c3d4e", -1},
    {"1.2.3-dev.z9z9z9z", "1.2.3-dev.a0a0a0a", 1},
    
    // Mixed sequence types
    {"1.2.3-dev.1", "1.2.3-dev.a1b2c3d", -1},
    {"1.2.3-dev.a1b2c3d", "1.2.3-dev.1", 1},
    
    // Edge cases
    {"0.0.0", "0.0.1", -1},
    {"0.1.0", "0.2.0", -1},
    {"1.0.0-rc.1", "1.0.0-rc.01", 0}, // leading zeros normalized
}
```

### Design Notes

1. **Stage Ordering**: dev < alpha < beta < rc (standard prerelease hierarchy)
2. **Release Precedence**: Release always > prerelease of same core
3. **Numeric vs Hash**: Numeric sorts lower (for consistent ordering across build systems)
4. **Lexicographic Hashing**: String hashes compared lexicographically (natural SHA order)
5. **Transitivity**: Comparison must be transitive (if A < B and B < C, then A < C)

## Testing

- [ ] Unit test: 100+ comparison pairs in table-driven format
- [ ] Unit test: Stage ordering is strict
- [ ] Unit test: Release > prerelease
- [ ] Unit test: Numeric sequence comparison
- [ ] Unit test: Hash sequence comparison
- [ ] Unit test: Mixed sequence types
- [ ] Unit test: Transitivity property holds
- [ ] Unit test: Edge cases (zero, very large numbers, leading zeros)

## Related Tickets

- T006: Version model
- T007: Parser
- T008: Normalizer
- T009: This ticket
- T013: `version compare` command (uses comparator)

## Notes

- Ensure comparison is deterministic and reproducible
- Document stage ordering precedence clearly
- Test with real-world version strings from golden corpus
