# T007: Implement Version Parser (SemVer, PEP440, Ecosystems)

**Phase**: 1 - Core Version Engine  
**Category**: Domain Logic  
**Complexity**: High  
**Estimated Duration**: 5-6 hours

## Objective

Implement robust parsing of version strings in SemVer, PEP440, and ecosystem-specific formats into the canonical Version model.

## Current State

- No parser exists

## Target State

- `version.Parser` interface is implemented
- Parser handles SemVer: `1.2.3`, `1.2.3-dev.4`, `1.2.3-rc.2`
- Parser handles PEP440-like: `1.2.3dev4`, `1.2.3a1`, `1.2.3b2`, `1.2.3rc3`
- Parser handles prefixes: `v1.2.3`, `1.2.3` (v-prefix optional)
- Parser detects sequence type (numeric, hash, build-id)
- Parser validates input and returns clear errors

## Acceptance Criteria

- [ ] `internal/version/parse.go` implements `Parser` interface
- [ ] Parses SemVer forms:
  - `1.2.3` → Version{1, 2, 3, StageFinal, nil, ...}
  - `1.2.3-dev.4` → Version{1, 2, 3, StageDev, 4, ...}
  - `1.2.3-rc.2` → Version{1, 2, 3, StageRC, 2, ...}
- [ ] Parses PEP440-like forms:
  - `1.2.3dev4` → Version{1, 2, 3, StageDev, 4, ...}
  - `1.2.3a1` → Version{1, 2, 3, StageAlpha, 1, ...}
  - `1.2.3b2` → Version{1, 2, 3, StageBeta, 2, ...}
  - `1.2.3rc3` → Version{1, 2, 3, StageRC, 3, ...}
- [ ] Handles v-prefix:
  - `v1.2.3` → Version{1, 2, 3, StageFinal, nil, ...}
  - Prefix is preserved in Original field
- [ ] Detects sequence type:
  - `4` → SeqTypeNumeric
  - `a1b2c3d` (7+ hex chars) → SeqTypeCommitSHA
  - `gh-12345` → SeqTypeBuildID
- [ ] Rejects invalid versions:
  - Non-numeric major/minor/patch
  - More than 3 numeric components (e.g., `1.2.3.4`)
  - Invalid stage names (e.g., `1.2.3-alpha.4` when expecting `a` or `alpha`)
- [ ] Error messages include suggestions
- [ ] Parser is case-insensitive for stage names

## Context

### Files to Create

- `internal/version/parse.go` — main parser implementation
- `internal/version/parse_test.go` — comprehensive tests

### Parsing Strategy

Use regex patterns for initial tokenization, then semantic validation:

```go
// Patterns to match (in order)
// 1. Optional prefix (v or empty)
// 2. Major.Minor.Patch (required)
// 3. Optional -Stage.Sequence or StageSequence
// 4. Optional metadata (ignored for now)

type parser struct{}

func (p *parser) Parse(input string) (*Version, error) {
    normalized := strings.TrimSpace(input)
    
    // Try SemVer format first (most common)
    if v, err := p.parseSemVer(normalized); err == nil {
        return v, nil
    }
    
    // Try PEP440 format
    if v, err := p.parsePEP440(normalized); err == nil {
        return v, nil
    }
    
    return nil, fmt.Errorf("parse error: %q matches no known version format", input)
}
```

### Test Corpus

Create comprehensive test table covering:
- All SemVer patterns (final, dev, alpha, beta, rc)
- All PEP440 patterns
- Prefixes (v, empty)
- Numeric and hash sequences
- Build IDs
- Edge cases (leading zeros, very large numbers)
- Errors (invalid stage, too many components, non-numeric core)

Example test cases:

```
"1.2.3" → Valid SemVer
"v1.2.3" → Valid SemVer with v prefix
"1.2.3-dev.4" → Valid SemVer with dev prerelease
"1.2.3-dev.a1b2c3d" → Valid SemVer with hash sequence
"1.2.3dev4" → Valid PEP440
"1.2.3a1" → Valid PEP440 alpha
"1.2.3rc1" → Valid PEP440 rc
"1.2.3-rc.2" → Valid SemVer rc
"v1.2.3-rc.2" → Valid SemVer with prefix
"1.2.3-ghaction.gh-12345" → GitHub Actions build ID
"1.2.3.4" → Invalid (too many components)
"1.2.x" → Invalid (non-numeric patch)
"1.2.3-xxx.4" → Invalid (unknown stage)
```

### Design Notes

1. **Greedy Matching**: Try most specific patterns first
2. **Error Recovery**: Provide suggestions for common mistakes
3. **Case Insensitivity**: Accept `RC`, `rc`, `Rc` for stage names
4. **Sequence Detection**: Inspect sequence to determine type (numeric vs hash)
5. **Metadata**: Ignore metadata (anything after `+` in SemVer) for now

## Testing

- [ ] Unit test: Table-driven tests for 50+ version strings
- [ ] Unit test: All SemVer patterns parse correctly
- [ ] Unit test: All PEP440 patterns parse correctly
- [ ] Unit test: Sequence types are detected correctly
- [ ] Unit test: Invalid versions return clear errors with suggestions
- [ ] Unit test: v-prefix is handled correctly
- [ ] Integration test: `verge parse 1.2.3` outputs correct normalized form

## Related Tickets

- T006: Version model
- T008: Normalizer (works with parsed versions)
- T009: Comparator (works with parsed versions)
- T011: Renderer (produces output from parsed versions)
- T014: Unit tests (cover parser exhaustively)

## Notes

- Regex should be compiled once and reused (performance)
- Consider using `regexp` package for pattern matching
- Document parsing precedence and ambiguities in comments
