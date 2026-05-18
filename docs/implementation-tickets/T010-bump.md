# T010: Implement Version Bump Logic

**Phase**: 1 - Core Version Engine  
**Category**: Domain Logic  
**Complexity**: High  
**Estimated Duration**: 4-5 hours

## Objective

Implement version bump logic for major, minor, patch, prerelease, and final release transitions.

## Current State

- No bump logic exists

## Target State

- `Bumper` interface computes next version for given bump kind
- All bump kinds are supported (major, minor, patch, prerelease, final)
- Bump logic handles prerelease sequences correctly
- Clear error messages for invalid bumps

## Acceptance Criteria

- [ ] `internal/version/bump.go` implements `Bumper` interface
- [ ] Bump kind `major`:
  - `1.2.3` → `2.0.0`
  - `1.2.3-rc.1` → `2.0.0`
  - `2.0.0-dev.1` → `2.0.0` (if final) or `3.0.0` (if major)
- [ ] Bump kind `minor`:
  - `1.2.3` → `1.3.0`
  - `1.2.3-rc.1` → `1.3.0`
- [ ] Bump kind `patch`:
  - `1.2.3` → `1.2.4`
  - `1.2.3-rc.1` → `1.2.3` (finalize) or `1.2.4` (bump)
- [ ] Bump kind `prerelease`:
  - `1.2.3` with stage `dev` → `1.2.4-dev.1`
  - `1.2.3-dev.1` with stage `dev` → `1.2.3-dev.2`
  - `1.2.3-dev.1` with stage `rc` → `1.2.3-rc.1` (stage change resets sequence)
- [ ] Bump kind `final`:
  - `1.2.3-rc.1` → `1.2.3`
  - `1.2.3-dev.1` → `1.2.3-dev.1` error (invalid transition)
- [ ] Sequence handling:
  - Numeric sequences increment: `4` → `5`
  - Hash sequences cannot bump prerelease (error)
- [ ] Clear errors for invalid bumps

## Context

### Files to Create

- `internal/version/bump.go` — bumper implementation
- `internal/version/bump_test.go` — table-driven tests

### Bump Algorithm

```go
type bumper struct{}

type BumpKind string

const (
    BumpMajor      BumpKind = "major"
    BumpMinor      BumpKind = "minor"
    BumpPatch      BumpKind = "patch"
    BumpPrerelease BumpKind = "prerelease"
    BumpFinal      BumpKind = "final"
)

func (b *bumper) Bump(v *Version, kind BumpKind, stage Stage) (*Version, error) {
    switch kind {
    case BumpMajor:
        return b.bumpMajor(v), nil
    case BumpMinor:
        return b.bumpMinor(v), nil
    case BumpPatch:
        return b.bumpPatch(v), nil
    case BumpPrerelease:
        return b.bumpPrerelease(v, stage)
    case BumpFinal:
        return b.bumpFinal(v)
    default:
        return nil, fmt.Errorf("unknown bump kind: %s", kind)
    }
}

func (b *bumper) bumpMajor(v *Version) *Version {
    return &Version{
        Major:    v.Major + 1,
        Minor:    0,
        Patch:    0,
        Stage:    StageFinal,
        Sequence: nil,
    }
}

func (b *bumper) bumpMinor(v *Version) *Version {
    return &Version{
        Major:    v.Major,
        Minor:    v.Minor + 1,
        Patch:    0,
        Stage:    StageFinal,
        Sequence: nil,
    }
}

func (b *bumper) bumpPatch(v *Version) *Version {
    return &Version{
        Major:    v.Major,
        Minor:    v.Minor,
        Patch:    v.Patch + 1,
        Stage:    StageFinal,
        Sequence: nil,
    }
}

func (b *bumper) bumpPrerelease(v *Version, stage Stage) (*Version, error) {
    // Only numeric sequences can bump prerelease
    seq, ok := v.Sequence.(int)
    if v.Sequence != nil && !ok {
        return nil, fmt.Errorf("cannot bump prerelease with hash sequence: %v", v.Sequence)
    }
    
    // If stage changes, reset sequence; otherwise increment
    if v.Stage != stage || v.Stage == StageFinal {
        // Transition to new stage
        return &Version{
            Major:        v.Major,
            Minor:        v.Minor,
            Patch:        v.Patch,
            Stage:        stage,
            Sequence:     1,
            SequenceType: SeqTypeNumeric,
        }, nil
    }
    
    // Same stage, increment sequence
    return &Version{
        Major:        v.Major,
        Minor:        v.Minor,
        Patch:        v.Patch,
        Stage:        stage,
        Sequence:     seq + 1,
        SequenceType: SeqTypeNumeric,
    }, nil
}

func (b *bumper) bumpFinal(v *Version) (*Version, error) {
    if v.Stage == StageFinal {
        return nil, fmt.Errorf("already a final release")
    }
    return &Version{
        Major:    v.Major,
        Minor:    v.Minor,
        Patch:    v.Patch,
        Stage:    StageFinal,
        Sequence: nil,
    }, nil
}
```

### Test Corpus

```go
[]struct {
    input    string
    kind     BumpKind
    stage    Stage
    expected string
    error    bool
}{
    // Major bumps
    {"1.2.3", BumpMajor, StageFinal, "2.0.0", false},
    {"1.2.3-rc.1", BumpMajor, StageFinal, "2.0.0", false},
    
    // Minor bumps
    {"1.2.3", BumpMinor, StageFinal, "1.3.0", false},
    
    // Patch bumps
    {"1.2.3", BumpPatch, StageFinal, "1.2.4", false},
    
    // Prerelease bumps
    {"1.2.3", BumpPrerelease, StageDev, "1.2.4-dev.1", false},
    {"1.2.3-dev.1", BumpPrerelease, StageDev, "1.2.3-dev.2", false},
    {"1.2.3-dev.1", BumpPrerelease, StageRC, "1.2.3-rc.1", false},
    
    // Final bumps
    {"1.2.3-rc.1", BumpFinal, StageFinal, "1.2.3", false},
    {"1.2.3", BumpFinal, StageFinal, "", true}, // error: already final
    
    // Hash sequence errors
    {"1.2.3-dev.a1b2c3d", BumpPrerelease, StageDev, "", true},
}
```

## Testing

- [ ] Unit test: 40+ bump scenarios in table-driven format
- [ ] Unit test: All bump kinds work correctly
- [ ] Unit test: Prerelease sequence increments correctly
- [ ] Unit test: Stage transitions reset sequence
- [ ] Unit test: Hash sequences reject prerelease bump
- [ ] Unit test: Invalid bumps return clear errors

## Related Tickets

- T006: Version model
- T010: This ticket
- T020: `version bump` command (uses bumper)

## Notes

- Bump creates a new Version; original is not mutated
- Document bump strategy for each kind clearly
- Consider "auto" bump detection in Phase 4
