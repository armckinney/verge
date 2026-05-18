# T008: Implement Version Normalizer and Canonicalization

**Phase**: 1 - Core Version Engine  
**Category**: Domain Logic  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Implement version normalization to ensure consistent ordering and comparison across different input formats.

## Current State

- No normalizer exists

## Target State

- `Normalizer` interface converts any parsed Version to canonical form
- Normalization is idempotent
- Numeric sequence is always an int
- Hash sequences are standardized (lowercase, fixed length for SHAs)
- Stage names are canonical

## Acceptance Criteria

- [ ] `internal/version/normalize.go` implements `Normalizer` interface
- [ ] All parsed versions are converted to canonical form:
  - Sequence is converted to appropriate type (int, string)
  - Stage is canonical enum value (never string)
  - Scheme is set to detected pattern
- [ ] Numeric sequences are parsed to int:
  - `"4"` → `4` (int)
  - `"0"` → `0` (int)
- [ ] Hash sequences are standardized:
  - Lowercase hexadecimal (if SHA)
  - Preserved length for content hashes
- [ ] Normalization is idempotent:
  - `Normalize(Normalize(v))` == `Normalize(v)`
- [ ] Build IDs are preserved as-is
- [ ] Errors are clear if sequence cannot be interpreted

## Context

### Files to Create

- `internal/version/normalize.go` — normalizer implementation
- `internal/version/normalize_test.go` — tests

### Normalization Rules

```go
type normalizer struct{}

func (n *normalizer) Normalize(v *Version) (*Version, error) {
    // Copy version
    normalized := *v
    
    // Ensure sequence is correct type
    if normalized.Sequence != nil {
        seqInt, seqStr := parseSequence(normalized.Sequence)
        if seqInt != nil {
            normalized.Sequence = *seqInt
            normalized.SequenceType = SeqTypeNumeric
        } else if seqStr != nil {
            normalized.Sequence = standardizeSequenceString(*seqStr, normalized.SequenceType)
        }
    }
    
    // Validate consistency
    if err := normalized.Validate(); err != nil {
        return nil, err
    }
    
    return &normalized, nil
}

// Helper to parse sequence from any type
func parseSequence(seq interface{}) (*int, *string) {
    switch v := seq.(type) {
    case int:
        return &v, nil
    case int64:
        i := int(v)
        return &i, nil
    case string:
        if i, err := strconv.Atoi(v); err == nil {
            return &i, nil
        }
        // It's a hash or build ID
        return nil, &v
    }
    return nil, nil
}

// Standardize string sequences (lowercase for SHA, preserve build IDs)
func standardizeSequenceString(s string, seqType SequenceType) string {
    if seqType == SeqTypeCommitSHA {
        return strings.ToLower(s)
    }
    // Preserve build IDs and other formats as-is
    return s
}
```

### Validation Rules

After normalization, Version must satisfy:
- Major >= 0, Minor >= 0, Patch >= 0
- If Stage != StageFinal, then Sequence must be non-nil and non-zero
- If Stage == StageFinal, then Sequence must be nil
- Stage is valid enum value (not string)
- SequenceType matches actual Sequence type

## Testing

- [ ] Unit test: Numeric string sequence converts to int
- [ ] Unit test: Numeric int sequence preserved
- [ ] Unit test: Hash sequence stays string and lowercased
- [ ] Unit test: Build ID preserved as-is
- [ ] Unit test: Normalization is idempotent
- [ ] Unit test: Invalid sequences return error
- [ ] Unit test: Final release has no sequence

## Related Tickets

- T006: Version model
- T007: Parser (produces input to normalizer)
- T009: Comparator (works with normalized versions)

## Notes

- Normalization happens after parsing but before comparison/rendering
- Consider normalizing `Sequence` to int64 for very large sequence numbers
