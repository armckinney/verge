# T003: CLI Scaffolding and Subcommand Structure

**Phase**: 0 - Bootstrap CLI  
**Category**: CLI Framework  
**Complexity**: Low  
**Estimated Duration**: 2-3 hours

## Objective

Establish the Cobra command structure for the CLI with root command and subcommands for the version domain.

## Current State

- No CLI command structure exists

## Target State

- Root command `verctl` with help text
- `verctl version` subcommand (top-level category)
- `verctl version parse`, `verctl version compare`, `verctl version current`, `verctl version latest`, `verctl version bump` subcommands
- All commands use consistent flag patterns (--format, --ecosystem, --explain, etc.)
- Help text is clear and includes ecosystem examples

## Acceptance Criteria

- [ ] `verctl --help` displays all subcommands
- [ ] `verctl version --help` displays all version subcommands
- [ ] `verctl version parse --help` shows relevant flags (--output, --ecosystem)
- [ ] Each command has a `Short` and `Long` description in Cobra
- [ ] `--verbose` flag is available on root and propagates to subcommands
- [ ] `--ecosystem` flag is available on commands that support it
- [ ] `--format` flag supports `text` and `json` values
- [ ] `--explain` flag is available on source-fetching commands
- [ ] All commands can be executed without error (even if they fail gracefully for missing data)
- [ ] Cobra command tree is documented in `internal/cli/` with clear package structure

## Context

### Files to Create

- `internal/cli/root.go` — root command definition
- `internal/cli/version.go` — version command group
- `internal/cli/version_parse.go` — parse subcommand
- `internal/cli/version_compare.go` — compare subcommand
- `internal/cli/version_current.go` — current subcommand (stub for now)
- `internal/cli/version_latest.go` — latest subcommand (stub for now)
- `internal/cli/version_bump.go` — bump subcommand (stub for now)

### Files to Update

- `cmd/verctl/main.go` — call `cli.Execute(rootCmd)` from Cobra root

### Command Flag Design

```
verctl [--verbose] version
  parse <version> [--output semver|pep440|auto] [--ecosystem go|python|containers|...]
  compare <version1> <version2> [--format text|json]
  current [--source git-tags|github-releases|ghcr] [--ecosystem ...] [--format text|json] [--explain]
  latest [--source ...] [--ecosystem ...] [--constraint "^1.2.3"] [--core 1.2.3] [--stage dev] [--channel rel|pr] [--format text|json] [--explain]
  bump --from <version> --kind major|minor|patch|prerelease|final [--stage dev] [--ecosystem ...] [--format text|json]
```

### Design Notes

1. **Cobra Structure**: Each command is a separate file with `&cobra.Command{}` definition
2. **Global Flags**: Implement a `GlobalFlags` struct in `cli/flags.go` to hold common flags
3. **Config Access**: Commands should receive config via context or dependency injection
4. **Error Handling**: Commands should use consistent error reporting (see T005)
5. **Stubs**: Phase 0 commands are stubs; they'll be implemented in later phases

## Testing

- [ ] Unit test: Cobra command tree can be parsed
- [ ] Integration test: `verctl version parse 1.2.3 --format json` executes (stub output)
- [ ] Integration test: `verctl version compare 1.2.3 1.2.4` executes (stub output)
- [ ] Integration test: All commands return exit code 0 (for now; actual logic in Phase 1+)

## Related Tickets

- T001: Project structure
- T002: Config loading
- T004: Output modes
- T005: Error codes
- T012: Implement `version parse` command (real implementation)
- T013: Implement `version compare` command (real implementation)
- T018: Implement `version current` command (real implementation)
- T019: Implement `version latest` command (real implementation)
- T020: Implement `version bump` command (real implementation)

## Notes

- Use `cobra.Command.PreRunE` to load config and validate flags
- Keep each subcommand in a separate file for maintainability
- Add contextual help examples for each command (e.g., `verctl version parse v1.2.3-rc.1`)
