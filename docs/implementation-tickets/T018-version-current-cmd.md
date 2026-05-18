# T018: Implement `version current` Command

**Phase**: 2 - Local Source Provider  
**Category**: CLI Commands  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Implement the `version current` command to fetch the current project version from configured sources.

## Current State

- Stub command exists

## Target State

- Fetches current version from git-tags provider
- Returns single version with provenance
- Supports ecosystem-specific output
- Supports explain mode to show selection logic

## Acceptance Criteria

- [ ] `internal/cli/version_current.go` implements full logic
- [ ] Command signature: `verctl version current [flags]`
- [ ] `--source` flag to override config (git-tags, etc.)
- [ ] `--ecosystem` flag for output format
- [ ] `--format` flag (text, json)
- [ ] `--explain` flag shows selection process
- [ ] Returns error with exit code 20 if no version found
- [ ] Text output shows version, source, ecosystem
- [ ] JSON output includes all metadata

## Context

### Files to Update

- `internal/cli/version_current.go`

### Command Behavior

```bash
# Default: use config, output to text
$ verctl version current
Version:   1.2.3
Ecosystem: go
Source:    git-tags
Rendered:  v1.2.3

# With explain
$ verctl version current --explain
Candidates from git-tags:
  v1.2.3 (final)
  v1.2.3-rc.1 (prerelease)
  v1.2.2

Selected: v1.2.3 (highest final version)
Version:   1.2.3
Ecosystem: go
Source:    git-tags
Rendered:  v1.2.3

# JSON output
$ verctl version current --format json
{
  "version": "v1.2.3",
  "normalized": "1.2.3",
  "ecosystem": "go",
  "source": "git-tags",
  "timestamp": "2025-01-15T10:30:00Z",
  "sequenceType": "numeric"
}

# Override source
$ verctl version current --source git-tags
```

### Implementation Pattern

```go
var currentCmd = &cobra.Command{
    Use:   "current",
    Short: "Get current project version",
    RunE:  runVersionCurrent,
}

func runVersionCurrent(cmd *cobra.Command, args []string) error {
    cfg := getConfig()
    explain := cmd.Flag("explain").Changed
    
    // Get provider
    provider := getProvider(cfg) // git-tags by default
    
    // Fetch current version
    result, err := provider.GetCurrent(context.Background(), QueryOptions{
        IncludePrerelease: cfg.Sources.GitTags.IncludePrerelease,
    })
    if err != nil {
        return fmt.Errorf("exit code 20: no version found: %w", err)
    }
    
    // Render to target ecosystem
    rendered, _ := renderer.Render(result.Version, cfg.Ecosystem)
    
    // Output explain trace if requested
    if explain {
        showExplain(result)
    }
    
    // Format and output
    formatter.FormatVersionCurrent(result, rendered)
    return nil
}
```

## Testing

- [ ] Unit test: Fetch version from mock provider
- [ ] Integration test: `verctl version current` returns correct version
- [ ] Integration test: `--explain` shows candidate list
- [ ] Integration test: JSON output is valid
- [ ] Error test: No git tags → exit code 20

## Related Tickets

- T003: CLI scaffolding
- T004: Output modes
- T015: Provider interface
- T016: Git-tags provider
- T020: Explain mode

## Notes

- Explain mode should show all candidates considered and why one was selected
- Consider adding --source override for testing
