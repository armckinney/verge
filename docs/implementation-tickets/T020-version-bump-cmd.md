# T020: Implement `version bump` Command

**Phase**: 2 - Local Source Provider  
**Category**: CLI Commands  
**Complexity**: Low  
**Estimated Duration**: 2-3 hours

## Objective

Implement the `version bump` command to compute the next version.

## Current State

- Stub command exists

## Target State

- Bumps from a given version to next version
- Supports all bump kinds (major, minor, patch, prerelease, final)
- Supports ecosystem-specific output
- Clear output of source and destination

## Acceptance Criteria

- [ ] `internal/cli/version_bump.go` implements full logic
- [ ] Command signature: `verge version bump --from <version> --kind <kind> [flags]`
- [ ] `--from` flag (required) — current version
- [ ] `--kind` flag (required) — bump type
- [ ] `--stage` flag (optional, for prerelease) — target stage
- [ ] `--ecosystem` flag for output format
- [ ] `--format` flag (text, json)
- [ ] Returns new version and exit code 0
- [ ] Clear error for invalid bumps

## Context

### Files to Update

- `internal/cli/version_bump.go`

### Command Behavior

```bash
$ verge version bump --from 1.2.3 --kind minor
From:     1.2.3
Bump:     minor
To:       1.3.0
Ecosystem: go
Rendered: v1.3.0

$ verge version bump --from 1.2.3 --kind prerelease --stage dev
From:     1.2.3
Bump:     prerelease
Stage:    dev
To:       1.2.4-dev.1

$ verge version bump --from 1.2.3-rc.1 --kind final
From:     1.2.3-rc.1
Bump:     final
To:       1.2.3
```

## Testing

- [ ] Unit test: Bump from various versions
- [ ] Integration test: All bump kinds work
- [ ] Error test: Invalid bumps return clear error

## Related Tickets

- T010: Bumper implementation (dependency)

## Notes

- Use bumper from T010
- Output should show before/after clearly
