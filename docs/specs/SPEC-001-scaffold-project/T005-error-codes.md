# T005: Error Codes and Handling Strategy

**Phase**: 0 - Bootstrap CLI  
**Category**: Error Handling  
**Complexity**: Low  
**Estimated Duration**: 2-3 hours

## Objective

Define a consistent error code taxonomy and error handling patterns for the CLI.

## Current State

- No standardized error handling

## Target State

- All CLI errors use consistent exit codes
- Error messages are clear and actionable
- Errors include suggestions for common problems
- JSON output includes error details for machine parsing
- Error codes are documented for CI/automation

## Acceptance Criteria

- [ ] `internal/cli/errors.go` defines error types (ParseError, ProviderError, ConfigError, etc.)
- [ ] Exit codes are consistent across all commands:
  - `0`: Success
  - `1`: General error / unexpected failure
  - `2`: Command-line usage error (missing required flag, etc.)
  - `10`: Version comparison result (left < right)
  - `11`: Version comparison result (left > right)
  - `20`: Version not found in source
  - `21`: Invalid version format
  - `22`: Config file error
  - `30`: Network/provider error
- [ ] All error types implement `error` interface with consistent message format
- [ ] Error messages include suggestions (e.g., "Did you mean `--ecosystem`?")
- [ ] `--verbose` flag provides additional context (stack trace, raw API responses)
- [ ] JSON error output includes error code and message fields
- [ ] All commands handle errors gracefully without panicking
- [ ] `verge info` outputs version metadata (version, commit, date)

## Context

### Files to Create

- `internal/cli/errors.go` — error type definitions
- `internal/cli/error_handler.go` — error formatting and reporting
- `internal/version/errors.go` — version domain errors
- `internal/providers/errors.go` — provider errors

### Error Struct Design

```go
type CLIError struct {
    Code       int         // exit code
    Message    string      // user-facing message
    Details    string      // verbose details
    Suggestion string      // actionable suggestion
    Cause      error       // underlying error
}

type ParseError struct {
    Input  string
    Reason string
}

type ProviderError struct {
    Provider string
    Operation string // GetCurrent, GetLatest, List
    Message  string
    Retryable bool
}

type ConfigError struct {
    Field   string
    Message string
}
```

### Exit Code Reference

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | `verge current` found version |
| 1 | General error | Unexpected panic or error |
| 2 | Usage error | Missing required flag |
| 10 | Comparison: left < right | Used in CI gates |
| 11 | Comparison: left > right | Used in CI gates |
| 20 | Version not found | No tags in git repository |
| 21 | Invalid format | Unparseable version string |
| 22 | Config error | Invalid YAML or missing required field |
| 30 | Network error | GitHub API timeout or 403 Forbidden |

### Error Message Examples

```
# ParseError (exit 21)
Error: Invalid version format: "1.2.3.4.5"
Reason: Too many numeric components
Suggestion: Use semver format like "1.2.3" or "1.2.3-rc.1"

# ProviderError (exit 30)
Error: Failed to fetch from GitHub Releases
Reason: 403 Forbidden (rate limit exceeded)
Suggestion: Set $GITHUB_TOKEN environment variable for higher rate limits
Retryable: true

# ConfigError (exit 22)
Error: Invalid configuration
Field: sources.github-releases.owner
Reason: Required field missing
Suggestion: Set 'owner' in .verge.yaml or via --repo-owner flag
```

### JSON Error Output

```json
{
  "error": {
    "code": 21,
    "message": "Invalid version format",
    "details": "Version '1.2.3.4.5' has too many numeric components",
    "suggestion": "Use semver format like '1.2.3' or '1.2.3-rc.1'",
    "field": null
  }
}
```

### Design Notes

1. **Wrapping**: Use `fmt.Errorf("context: %w", err)` for error chains
2. **Suggestions**: Always provide next steps when possible
3. **Logging**: Log full error chain with `--verbose` but show user-friendly message by default
4. **Retryable Errors**: Distinguish between transient (network) and permanent (parse) errors
5. **Machine Parsing**: JSON output should allow automated error handling in CI

## Testing

- [ ] Unit test: Each error type formats correctly
- [ ] Unit test: Exit codes are correct for each error type
- [ ] Unit test: Error suggestions are clear and actionable
- [ ] Integration test: `verge parse invalid --format json` includes error code
- [ ] Integration test: `verge current` when no git tags exists returns exit code 20

## Related Tickets

- T003: CLI scaffolding (uses error handling)
- T004: Output modes (error output formatting)
- T012-T020: Individual commands use error handling

## Notes

- Document error codes in README for CI integration
- Ensure error messages do not leak sensitive information (tokens, paths)
- Test error handling with various failure scenarios (network, disk, permissions)
