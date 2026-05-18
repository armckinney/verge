# T001: Migrate Project Structure from API Template to CLI

**Phase**: 0 - Bootstrap CLI  
**Category**: Project Setup  
**Complexity**: Medium  
**Estimated Duration**: 4-6 hours

## Objective

Migrate the repository from an API server template structure into a CLI application structure, establishing the foundation for all subsequent development.

## Current State

- Project is currently structured as an HTTP API server
- `cmd/api/main.go` is the entry point
- `internal/` contains server, database, handlers, middleware, models, repository
- `tests/` contains database setup and environment scripts

## Target State

New project structure:

```
cmd/verge/
  main.go              # CLI entry point (replaces cmd/api/main.go)
internal/
  cli/                 # NEW: command definitions and orchestration
    root.go
    version.go
  version/             # NEW: domain logic for version management
    model.go
    parse.go
    normalize.go
    compare.go
    bump.go
    render.go
  ecosystems/          # NEW: ecosystem-specific logic
    registry.go
    types.go
    go.go
    python.go
    containers.go
    terraform.go
    github_actions.go
  providers/           # NEW: version source integrations
    provider.go
    git_tags.go
  sequence/            # NEW: sequence type detection
    interpreter.go
  config/              # EXISTING but refactored
    load.go
    schema.go
    (remove database configs)
tests/
  fixtures/            # NEW: test data for versioning
    versions.go
    (existing db setup can be removed or archived)
.verge.yaml           # NEW: default config file template
go.mod                 # Update module path reference
```

## Acceptance Criteria

- [ ] `cmd/verge/main.go` exists and compiles
- [ ] All old API-related code in `internal/` is removed or moved to an archive branch
- [ ] `internal/cli/root.go` imports properly and defines a Cobra root command
- [ ] `internal/version/model.go` defines the canonical `Version` struct
- [ ] `internal/config/schema.go` defines the CLI config schema (YAML)
- [ ] `go.mod` is updated with current dependencies (remove HTTP server libs if possible)
- [ ] `make build` produces a `verge` binary (not `main`)
- [ ] `make run` is updated to run the CLI (or removed in favor of `make build`)
- [ ] Old database connection logic and environment variable loading for API is removed
- [ ] `.verge.yaml` template file exists with basic structure

## Context

### Files to Create

- `cmd/verge/main.go`
- `internal/cli/root.go`
- `internal/cli/version.go`
- `internal/version/model.go`
- `internal/config/schema.go`
- `.verge.yaml` (template)

### Files to Remove

- `cmd/api/` (entire directory)
- `internal/server/` (HTTP server logic)
- `internal/database/` (database connections)
- `internal/handlers/` (HTTP handlers)
- `internal/middleware/` (HTTP middleware)
- `internal/models/user.go` (domain model, can be archived)
- `internal/repository/` (database queries)

### Files to Update

- `go.mod` — remove HTTP and database dependencies if no longer needed
- `Makefile` — update build target from `main` to `verge`; add `build-snapshot` target for goreleaser
- `README.md` — update Getting Started section for CLI workflow
- `.gitignore` — add `/dist/` directory for goreleaser artifacts
- `tests/` — archive existing database tests or move to separate branch

### Build System

During this migration, establish the foundation for cross-platform builds:
- Initialize goreleaser configuration (detailed in T031)
- Update Makefile to support build targets compatible with goreleaser
- Add `.goreleaser.yaml` to version control
- Define build variables (version, commit, date) that can be passed to Go at compile time

### Design Notes

1. **Cobra Framework**: Use `github.com/spf13/cobra` for command structure (already in go.mod)
2. **Main Entry Point**: Minimal `cmd/verge/main.go` that just calls `cli.Execute()`
3. **Version Model**: The canonical `Version` struct will be the core abstraction:
   ```go
   type Version struct {
       Major    int
       Minor    int
       Patch    int
       Stage    Stage       // enum: none, dev, alpha, beta, rc
       Sequence interface{} // int or string (hash)
       Original string      // raw input
       Scheme   Scheme      // enum: semver, pep440, auto
   }
   ```

4. **Config Schema**: Start simple; ecosystem detection and provider precedence are the key settings.

## Testing

- [ ] Unit test: `root.go` can instantiate the root command
- [ ] Integration: `verge --version` returns version string
- [ ] Build check: `go build ./cmd/verge` succeeds with no errors

## Related Tickets

- T002: Config loading and schema validation
- T003: CLI scaffolding (subcommands)
- T004: Output modes (text/JSON)
- T005: Error codes
- T031: Cross-platform build system (goreleaser setup)

## Notes

- This ticket removes old code; ensure any IP or important logic is archived before deletion.
- The Makefile target should still allow `make test` to work for any unit tests we keep.
- Database setup scripts in `tests/db/` can be archived or removed; they're no longer needed for CLI.
