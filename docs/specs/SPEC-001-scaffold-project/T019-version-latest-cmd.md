# T019: Implement `version latest` Command

**Phase**: 2 - Local Source Provider  
**Category**: CLI Commands  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Implement the `version latest` command to fetch the latest available version with optional filtering.

## Current State

- Stub command exists

## Target State

- Fetches latest version from providers
- Supports version constraints (`--core`, `--stage`)
- Supports explain mode
- Handles no results gracefully

## Acceptance Criteria

- [ ] `internal/cli/version_latest.go` implements full logic
- [ ] Command signature: `verge latest [flags]`
- [ ] `--source` flag (git-tags, etc.)
- [ ] `--constraint` flag for version ranges (v1, future enhancement)
- [ ] `--core` flag to filter by core version (e.g., `1.2.3`)
- [ ] `--stage` flag to filter by stage (dev, alpha, beta, rc, final)
- [ ] `--ecosystem` flag for output format
- [ ] `--format` flag (text, json)
- [ ] `--explain` flag shows filtering and selection
- [ ] Error handling: no matching versions

## Context

### Files to Update

- `internal/cli/version_latest.go`

### Command Behavior

```bash
# Highest version
$ verge latest
Version:   2.0.0
Source:    git-tags

# Latest of specific core
$ verge latest --core 1.2.3 --stage dev
Version:   1.2.3-dev.5
Source:    git-tags

# With explain
$ verge latest --explain
Candidates from git-tags:
  v2.0.0 (final)
  v1.2.3 (final)
  v1.2.3-dev.5 (prerelease)
  v1.2.3-dev.4 (prerelease)
  ...

Selected: v2.0.0 (highest version)
Version:   2.0.0
Source:    git-tags
```

## Testing

- [ ] Unit test: Get latest from mock provider
- [ ] Integration test: `verge latest` returns highest
- [ ] Integration test: `--core 1.2.3 --stage dev` filters correctly
- [ ] Integration test: `--explain` shows candidates
- [ ] Error test: No matching versions

## Related Tickets

- T003: CLI scaffolding
- T016: Git-tags provider
- T020: Explain mode

## Notes

- Filtering should be applied after fetching all versions
- Explain mode should show how filtering was applied
