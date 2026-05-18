# T006: Define Canonical Version Model and Interfaces

**Phase**: 1 - Core Version Engine  
**Category**: Domain Logic  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Define the canonical internal `Version` model and all supporting types/interfaces for version domain logic.

## Current State

- No version domain types exist

## Target State

- `Version` struct represents all version information internally
- All version operations (parse, normalize, compare, bump, render) work with this model
- Interface contracts are defined for extensibility
- All types include comprehensive validation and documentation

## Acceptance Criteria

- [ ] `internal/version/model.go` defines:
  - `Version` struct with Major, Minor, Patch, Stage, Sequence fields
  - `Stage` enum type (none, dev, alpha, beta, rc) with string conversion methods
  - `Scheme` enum type (semver, pep440, auto)
  - `SequenceType` enum type (numeric, commit-sha, content-hash, build-id)
- [ ] All types have `String()` methods for debugging
- [ ] `Version` type has validation logic to ensure consistency
- [ ] `internal/version/types.go` defines public interfaces:
  - `Parser` interface for version string parsing
  - `Normalizer` interface for canonicalization
  - `Comparator` interface for ordering
  - `Renderer` interface for ecosystem-specific formatting
- [ ] All public types and interfaces are documented with examples
- [ ] Type conversions (e.g., stage string → Stage enum) are well-defined

## Context

### Files to Create

- `internal/version/model.go` — core types
- `internal/version/types.go` — interfaces and contracts
- `internal/version/helpers.go` — utility functions for Stage/Scheme conversions

### Version Model Design

```go
package version

import "fmt"

type Stage int

const (
    StageFinal Stage = iota // release version
    StageDev
    StageAlpha
    StageBeta
    StageRC
)

func (s Stage) String() string { ... }
func StageFromString(s string) (Stage, error) { ... }

type Scheme string

const (
    SchemeSemVer Scheme = "semver"
    SchemePEP440 Scheme = "pep440"
    SchemeAuto   Scheme = "auto"
)

type SequenceType string

const (
    SeqTypeNumeric    SequenceType = "numeric"
    SeqTypeCommitSHA  SequenceType = "commit-sha"
    SeqTypeContentHash SequenceType = "content-hash"
    SeqTypeBuildID    SequenceType = "build-id"
    SeqTypeUnknown    SequenceType = "unknown"
)

// Version is the canonical internal representation
type Version struct {
    // Core version components
    Major int
    Minor int
    Patch int

    // Prerelease/development info
    Stage    Stage
    Sequence interface{} // int or string

    // Metadata
    SequenceType SequenceType
    Original     string
    Scheme       Scheme
}

func (v *Version) String() string { ... }
func (v *Version) Validate() error { ... }
func (v *Version) IsPrerelease() bool { ... }
func (v *Version) Core() string { ... } // "1.2.3"
```

### Interface Contracts

```go
// Parser converts a version string into a Version struct
type Parser interface {
    Parse(input string) (*Version, error)
}

// Normalizer canonicalizes a Version for comparison
type Normalizer interface {
    Normalize(v *Version) (*Version, error)
}

// Comparator determines ordering between two versions
type Comparator interface {
    Compare(left, right *Version) int // -1, 0, or 1
}

// Renderer converts Version to ecosystem-specific format
type Renderer interface {
    Render(v *Version, ecosystem string) (string, error)
}

// Bumper computes next version given bump kind
type Bumper interface {
    Bump(v *Version, kind BumpKind, stage Stage) (*Version, error)
}

type BumpKind string

const (
    BumpMajor     BumpKind = "major"
    BumpMinor     BumpKind = "minor"
    BumpPatch     BumpKind = "patch"
    BumpPrerelease BumpKind = "prerelease"
    BumpFinal     BumpKind = "final"
)
```

### Design Notes

1. **Sequence Flexibility**: Sequence can be int or string (hash) for flexibility
2. **Validation**: Each Version should validate its state (e.g., Stage==StageFinal requires Sequence==nil)
3. **Immutability**: Treat Version structs as logically immutable after creation
4. **String Representations**: Provide multiple string views (canonical, debug, display)
5. **Type Safety**: Use enums (iota) for Stage/Scheme to prevent invalid values

## Testing

- [ ] Unit test: All Stage values convert to/from strings correctly
- [ ] Unit test: Version.String() produces meaningful output
- [ ] Unit test: Version.Validate() catches invalid combinations
- [ ] Unit test: Version.IsPrerelease() returns correct bool
- [ ] Unit test: StageFromString() handles case-insensitive input

## Related Tickets

- T007: Parser implementation
- T008: Normalizer implementation
- T009: Comparator implementation
- T011: Renderer implementation
- T010: Bumper implementation

## Notes

- Keep types pure and side-effect-free; all I/O happens in parser/renderer/provider layers
- Document numeric sequence limits (e.g., uint32 vs int64)
