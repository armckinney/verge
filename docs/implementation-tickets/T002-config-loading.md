# T002: Config Loading and Schema Validation

**Phase**: 0 - Bootstrap CLI  
**Category**: Configuration  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Implement the configuration system that loads `.verge.yaml`, validates against schema, and exposes configuration to all CLI commands.

## Current State

- No CLI config system exists
- Old API used `.env` files for environment variables

## Target State

- `.verge.yaml` is loaded from project root or config path
- Configuration is validated against a JSON schema at startup
- Config is available as a singleton or dependency-injected struct to all commands
- Environment variables can override config values
- CLI flags have highest precedence
- Default config (hardcoded) is used as fallback

## Acceptance Criteria

- [ ] `config.Schema` struct definition covers all config fields from semver-cli-spec section 6
- [ ] `config.Load()` loads `.verge.yaml` from current directory or `$VERCTL_CONFIG` path
- [ ] `config.Validate()` checks required fields and logs meaningful errors
- [ ] Environment variables (`VERCTL_*`) can override config file values
- [ ] `config.Default()` returns safe defaults (git-tags enabled, semver output, etc.)
- [ ] Config is accessible throughout CLI via a global or context-injected variable
- [ ] `verge version --config /path/to/custom.yaml` works
- [ ] YAML parsing handles missing optional fields gracefully
- [ ] Invalid YAML returns clear error message with line number

## Context

### Files to Create

- `internal/config/schema.go` — defines `Config` struct with YAML tags
- `internal/config/load.go` — implements loading logic, validation, and precedence
- `internal/config/defaults.go` — default config values

### Files to Update

- `cmd/verge/main.go` — load config at startup
- `internal/cli/root.go` — expose config to subcommands

### Config Schema Structure

```yaml
version: 1

# Target output ecosystem
ecosystem: go

format:
  input: auto
  output: auto
  tagPrefix: v
  sequenceInterpreter: auto

sources:
  precedence:
    - git-tags
    - github-releases
    - ghcr
  git-tags:
    enabled: true
    fetch: false
    includePrerelease: true
    ecosystemParsing: go
  github-releases:
    enabled: false
    owner: ""
    repo: ""
    includePrerelease: true
    includeDrafts: false
  ghcr:
    enabled: false
    image: ""
    includePrerelease: true
    channelFilter: null
  pypi:
    enabled: false
    packageName: ""
    includePrerelease: true
  terraform-registry:
    enabled: false
    module: ""

sequence:
  hashLength: 7
  allowContentHash: true
  ghBuildPattern: "gh-"

rules:
  prereleaseStage: dev
  allowMajorZeroBreaking: true
  defaultBump: patch

autoBump:
  conventionalCommits: true
  breakingTokens:
    - "BREAKING CHANGE"
    - "!:"
```

### Design Notes

1. **YAML Parsing**: Use `gopkg.in/yaml.v3` (already in go.mod)
2. **Precedence Order**: CLI flags > env vars > config file > defaults
3. **Partial Config**: Config file can specify only overrides; missing fields use defaults
4. **Token Secrets**: GitHub tokens must come from env vars (`GITHUB_TOKEN`), never from config file
5. **Error Handling**: Validation errors should list all issues, not just first

## Testing

- [ ] Unit test: Load `.verge.yaml` with all fields
- [ ] Unit test: Load config with missing optional fields → uses defaults
- [ ] Unit test: Load invalid YAML → clear error message
- [ ] Unit test: Override config value via env var `VERCTL_ECOSYSTEM=python`
- [ ] Unit test: Load non-existent config file → uses defaults (or error if required)
- [ ] Integration test: `verge version current` loads config successfully

## Related Tickets

- T001: Project structure
- T003: CLI scaffolding (consumes config)
- T004: Output modes (configured here)
- T005: Error codes (validation errors)

## Notes

- Ensure all GitHub tokens are loaded only from environment, never printed in debug output
- Config file should be optional for simple workflows (e.g., git-tags only)
- Consider allowing config to be specified via `--config` CLI flag
