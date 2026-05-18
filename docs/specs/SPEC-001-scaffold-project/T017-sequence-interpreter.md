# T017: Implement Sequence Interpreter (Type Detection)

**Phase**: 2 - Local Source Provider  
**Category**: Domain Logic  
**Complexity**: Medium  
**Estimated Duration**: 2-3 hours

## Objective

Implement the sequence type interpreter to detect whether a sequence is numeric, commit SHA, content hash, or build ID.

## Current State

- Parser identifies sequence exists but doesn't classify its type

## Target State

- Sequence interpreter determines type: numeric, commit-sha, content-hash, build-id
- Classification is deterministic and configurable
- Ordering is consistent across type boundaries

## Acceptance Criteria

- [ ] `internal/sequence/interpreter.go` defines detection logic
- [ ] Detects numeric sequences: `1`, `42`, `9999`
- [ ] Detects commit SHAs: 7+ hex characters (e.g., `a1b2c3d`)
- [ ] Detects content hashes: MD5 or SHA (32+ hex chars)
- [ ] Detects build IDs: `gh-12345`, `build-999`
- [ ] Configurable SHA length detection (default 7, min 6)
- [ ] Falls back to lexicographic for unknown types
- [ ] Numeric sequences always sort lower than hashes

## Context

### Files to Create

- `internal/sequence/interpreter.go`
- `internal/sequence/interpreter_test.go`

### Detection Rules

```go
type SequenceType string

const (
    TypeNumeric    SequenceType = "numeric"
    TypeCommitSHA  SequenceType = "commit-sha"
    TypeContentHash SequenceType = "content-hash"
    TypeBuildID    SequenceType = "build-id"
    TypeUnknown    SequenceType = "unknown"
)

func DetectType(seq interface{}, config SequenceConfig) SequenceType {
    s := fmt.Sprintf("%v", seq)
    
    // Check numeric
    if _, err := strconv.Atoi(s); err == nil {
        return TypeNumeric
    }
    
    // Check build ID patterns
    if strings.HasPrefix(s, config.GitHubBuildPattern) {
        return TypeBuildID
    }
    
    // Check hex string
    if isHex(s) {
        if len(s) >= config.HashLength {
            // Could be commit SHA or content hash
            if len(s) >= 32 {
                return TypeContentHash
            }
            return TypeCommitSHA
        }
    }
    
    return TypeUnknown
}

func isHex(s string) bool {
    for _, ch := range s {
        if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
            return false
        }
    }
    return true
}
```

### Configuration

From `.verge.yaml`:

```yaml
sequence:
  hashLength: 7                # minimum chars to detect as SHA
  allowContentHash: true
  ghBuildPattern: "gh-"
```

## Testing

- [ ] Unit test: Numeric detection (`1`, `42`, `0`)
- [ ] Unit test: Commit SHA detection (`a1b2c3d`, `abc1234def5678`)
- [ ] Unit test: Content hash detection (32+ chars)
- [ ] Unit test: Build ID detection (`gh-12345`, `build-999`)
- [ ] Unit test: Case-insensitive hex detection
- [ ] Unit test: Mixed alphanumeric fallback to unknown

## Related Tickets

- T016: Git-tags provider (uses interpreter)
- T009: Comparator (uses sequence type for ordering)

## Notes

- Keep interpretation consistent across all sources
- Document detection thresholds clearly
