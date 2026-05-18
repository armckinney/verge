# T028: Implement Conventional Commits Parser for Auto Bump Detection

**Phase**: 4 - CI/Release Enhancements  
**Category**: Release Features  
**Complexity**: High  
**Estimated Duration**: 3-4 hours

## Objective

Implement automatic version bump detection from conventional commits history.

## Current State

- No conventional commits support

## Target State

- Parser reads commit messages and determines bump type
- Supports `feat:`, `fix:`, `BREAKING CHANGE:` markers
- Integrates with `version bump --auto` flag
- Respects configuration for custom breaking tokens

## Acceptance Criteria

- [ ] `internal/version/conventional.go` implements parser
- [ ] Detects `feat:` → `minor` bump
- [ ] Detects `fix:` → `patch` bump
- [ ] Detects breaking change → `major` bump
- [ ] Customizable breaking tokens via config
- [ ] `--auto` flag for `version bump` command
- [ ] Integration: read git commit history (latest tag → HEAD)
- [ ] Configuration from `.verctl.yaml`
- [ ] Handles edge cases (no commits, no tags, etc.)

## Context

### Files to Create

- `internal/version/conventional.go`
- `internal/version/conventional_test.go`

### Configuration

```yaml
autoBump:
  conventionalCommits: true
  breakingTokens:
    - "BREAKING CHANGE"
    - "!:"
```

### Implementation Pattern

```go
func DetectBumpFromHistory() (BumpKind, error) {
    tags := getGitTags()
    if len(tags) == 0 {
        return "", errors.New("no git tags found")
    }
    
    latestTag := tags[0]
    commits, _ := getCommitsSince(latestTag)
    
    hasMajor := false
    hasMinor := false
    hasPatch := false
    
    for _, commit := range commits {
        if isBreakingChange(commit) {
            hasMajor = true
        } else if strings.HasPrefix(commit.Type, "feat") {
            hasMinor = true
        } else if strings.HasPrefix(commit.Type, "fix") {
            hasPatch = true
        }
    }
    
    if hasMajor {
        return BumpMajor, nil
    } else if hasMinor {
        return BumpMinor, nil
    } else if hasPatch {
        return BumpPatch, nil
    }
    return "", errors.New("no relevant commits found")
}
```

## Testing

- [ ] Unit test: Detect feat → minor
- [ ] Unit test: Detect fix → patch
- [ ] Unit test: Detect breaking → major
- [ ] Unit test: Custom breaking tokens
- [ ] Integration test: Real git history
- [ ] Error test: No tags, no commits

## Related Tickets

- T020: `version bump` command (--auto flag)

## Notes

- Use git commands or git2go library
- Cache commit history for performance
- Consider scope prefixes (feat(scope): message)
