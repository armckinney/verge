# T012: Implement `version parse` Command

**Phase**: 1 - Core Version Engine  
**Category**: CLI Commands  
**Complexity**: Low  
**Estimated Duration**: 2-3 hours

## Objective

Implement the `version parse` command to validate, parse, and display version information.

## Current State

- Stub command exists in CLI scaffolding

## Target State

- `verge version parse <version>` parses and normalizes input
- `--output` flag selects rendering format (semver, pep440, or ecosystem-specific)
- `--ecosystem` flag selects target ecosystem for rendering
- Output is available in text and JSON formats
- Error messages are clear for invalid input

## Acceptance Criteria

- [ ] `internal/cli/version_parse.go` implements full logic
- [ ] Command signature: `verge version parse <version> [flags]`
- [ ] `--output` flag with choices: semver, pep440, auto
- [ ] `--ecosystem` flag with choices: go, python, containers, terraform, github-actions, auto
- [ ] `--format` flag with choices: text, json
- [ ] Text output shows parsed components (major, minor, patch, stage, sequence, etc.)
- [ ] JSON output includes all parsed fields
- [ ] Invalid versions return clear error with suggestions
- [ ] Help text includes examples for each ecosystem

## Context

### Files to Update

- `internal/cli/version_parse.go` — implement parse command

### Command Behavior

```bash
# Basic parsing
$ verge version parse 1.2.3
Version:        1.2.3
Major:          1
Minor:          2
Patch:          3
Stage:          final
Sequence:       (none)
Scheme:         semver

# Prerelease
$ verge version parse v1.2.3-dev.4
Version:        v1.2.3-dev.4
Major:          1
Minor:          2
Patch:          3
Stage:          dev
Sequence:       4
Rendered (Go):  v1.2.3-dev.4
Rendered (Python): 1.2.3dev4
Scheme:         semver

# With JSON output
$ verge version parse 1.2.3 --format json
{
  "input": "1.2.3",
  "parsed": {
    "major": 1,
    "minor": 2,
    "patch": 3,
    "stage": "final",
    "sequence": null,
    "sequenceType": "unknown"
  },
  "schemeDetected": "semver",
  "rendered": {
    "go": "v1.2.3",
    "python": "1.2.3",
    "containers": "1.2.3",
    "terraform": "v1.2.3",
    "github-actions": "1.2.3"
  }
}
```

### Implementation Pattern

```go
var parseCmd = &cobra.Command{
    Use:   "parse <version>",
    Short: "Parse and validate a version string",
    Long:  `Parse a version string and display its components and rendered forms for all ecosystems.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVersionParse,
}

func runVersionParse(cmd *cobra.Command, args []string) error {
    input := args[0]
    
    // Parse
    parsed, err := parser.Parse(input)
    if err != nil {
        return fmt.Errorf("parse error: %w", err)
    }
    
    // Normalize
    normalized, err := normalizer.Normalize(parsed)
    if err != nil {
        return fmt.Errorf("normalize error: %w", err)
    }
    
    // Render for all ecosystems
    rendered := make(map[string]string)
    for _, eco := range []string{"go", "python", "containers", "terraform", "github-actions"} {
        out, err := renderer.Render(normalized, eco)
        if err == nil {
            rendered[eco] = out
        }
    }
    
    // Format and output
    formatter.FormatVersionParse(normalized, rendered)
    return nil
}
```

### Text Output Template

```
Input:          <original input>
Normalized:     <canonical form>
Scheme:         <detected>
Major:          <value>
Minor:          <value>
Patch:          <value>
Stage:          <stage name>
Sequence:       <value or (none)>
Sequence Type:  <type>

Rendered forms:
  Go:            <rendered>
  Python:        <rendered>
  Containers:    <rendered>
  Terraform:     <rendered>
  GitHub Actions:<rendered>
```

## Testing

- [ ] Unit test: Parse valid version strings
- [ ] Unit test: Reject invalid input with clear error
- [ ] Unit test: JSON output is valid JSON
- [ ] Unit test: All ecosystems render correctly
- [ ] Integration test: `verge version parse 1.2.3` succeeds
- [ ] Integration test: `verge version parse invalid --format json` returns error
- [ ] Integration test: Help text is available

## Related Tickets

- T003: CLI scaffolding
- T004: Output modes
- T007: Parser implementation (dependency)
- T008: Normalizer implementation (dependency)
- T011: Renderer implementation (dependency)

## Notes

- Show all rendered forms in text output for quick ecosystem reference
- Include sequence type in output (helps debug hash vs numeric)
