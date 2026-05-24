# T004: Uniform Output Modes (Text and JSON)

**Phase**: 0 - Bootstrap CLI  
**Category**: CLI Output  
**Complexity**: Low  
**Estimated Duration**: 2-3 hours

## Objective

Establish consistent output formatting for all CLI commands in both human-readable (text) and machine-readable (JSON) formats.

## Current State

- No output formatter exists

## Target State

- All commands support `--format text` (default) and `--format json`
- Text output is human-friendly and easy to read in terminals
- JSON output is fully valid and includes all relevant metadata
- Output is consistent across all commands
- JSON schema is documented for machine consumers

## Acceptance Criteria

- [ ] `cli/output.go` defines `Formatter` interface with `FormatVersion()` and `FormatList()`
- [ ] Text formatter produces readable output (e.g., `Version: 1.2.3, Source: git-tags`)
- [ ] JSON formatter produces valid JSON with schema (struct tags document fields)
- [ ] `--format json` can be piped to `jq` without errors
- [ ] Text output includes timestamp or other useful metadata
- [ ] Error output uses consistent format (separate from success output)
- [ ] `--format text` is default when not specified
- [ ] Invalid `--format` value produces clear error
- [ ] Output modes are tested with mock data

## Context

### Files to Create

- `internal/cli/output.go` — defines formatters
- `internal/cli/formats/text.go` — text formatter implementation
- `internal/cli/formats/json.go` — JSON formatter implementation
- `tests/fixtures/output_test.go` — golden tests for output

### Output Structures

All commands should return a normalized output struct:

```go
// For single version results
type VersionOutput struct {
    Version            string `json:"version"`
    NormalizedVersion  string `json:"normalizedVersion"`
    SchemeDetected     string `json:"schemeDetected"`
    Ecosystem          string `json:"ecosystem"`
    Source             string `json:"source"`
    Raw                string `json:"raw"`
    SequenceType       string `json:"sequenceType"` // numeric, commit-sha, etc
    Timestamp          string `json:"timestamp"`
}

// For comparison results
type ComparisonOutput struct {
    Left       string `json:"left"`
    Right      string `json:"right"`
    Comparison string `json:"comparison"` // "equal", "left<right", "left>right"
    ExitCode   int    `json:"exitCode"`
}

// For bump results
type BumpOutput struct {
    From       string `json:"from"`
    BumpKind   string `json:"bumpKind"`
    To         string `json:"to"`
    Ecosystem  string `json:"ecosystem"`
}
```

### Text Output Examples

```
# version current (text)
Version:        1.2.3
Scheme:         semver
Ecosystem:      go
Source:         git-tags
Normalized:     1.2.3
Raw:            v1.2.3

# version parse (text)
Input:          1.2.3-dev.4
Parsed:         1.2.3-dev.4
Rendered (Go):  v1.2.3-dev.4
Rendered (Python): 1.2.3dev4
Scheme:         semver
Stage:          dev
Sequence:       4

# version compare (text)
Left:           1.2.3
Right:          1.2.4
Result:         left < right
```

### JSON Output Examples

```json
{
  "version": "v1.2.3",
  "normalizedVersion": "1.2.3",
  "schemeDetected": "semver",
  "ecosystem": "go",
  "source": "git-tags",
  "raw": "v1.2.3",
  "sequenceType": "numeric",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

### Design Notes

1. **Formatter Interface**: Define as an interface so other formats can be added later
2. **Timestamp**: Include ISO 8601 for reproducibility
3. **Error Output**: Errors go to stderr; success to stdout
4. **Structured Logging**: Use consistent field names across all outputs
5. **Escape Special Characters**: JSON output properly escapes strings

## Testing

- [ ] Unit test: Text formatter produces readable output
- [ ] Unit test: JSON formatter produces valid JSON
- [ ] Unit test: JSON can be unmarshalled back to struct
- [ ] Integration test: `verge parse 1.2.3 --format json | jq .version` works
- [ ] Integration test: `verge parse 1.2.3 --format text` produces human-friendly output
- [ ] Golden test: Compare output against known good examples

## Related Tickets

- T003: CLI scaffolding
- T005: Error codes and handling
- T012-T020: Individual command implementations use this formatter

## Notes

- Keep text output under 80 characters where possible for terminal display
- Include all relevant metadata in JSON for debugging
- Consider color output for text mode (optional enhancement for Phase 2+)
