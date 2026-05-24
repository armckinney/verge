# T030: Implement Changelog-Friendly JSON Outputs

**Phase**: 4 - CI/Release Enhancements  
**Category**: Release Features  
**Complexity**: Medium  
**Estimated Duration**: 2-3 hours

## Objective

Implement enhanced JSON output for changelog and release note generation.

## Current State

- Basic JSON output exists

## Target State

- JSON includes metadata for changelog generation
- Version history and provenance information
- Diff information (commits between versions)
- Machine-readable release notes

## Acceptance Criteria

- [ ] `internal/cli/output_changelog.go` formats changelog-friendly output
- [ ] JSON includes:
  - Previous version (from provider)
  - Current version
  - Bump type (major/minor/patch)
  - List of commits between versions (if available)
  - Changelog-friendly metadata
- [ ] Integration: `version bump` and `version current` support `--changelog` flag
- [ ] Output is compatible with changelog tools (standard JSON structure)
- [ ] Configuration can control metadata included

## Context

### Files to Create

- `internal/cli/output_changelog.go`
- `internal/cli/output_changelog_test.go`

### Changelog JSON Schema

```json
{
  "version": {
    "from": "1.2.3",
    "to": "1.3.0",
    "bumpType": "minor"
  },
  "metadata": {
    "timestamp": "2025-01-15T10:30:00Z",
    "source": "git-tags",
    "commits": [
      {
        "hash": "abc1234",
        "message": "feat: add new feature",
        "type": "feat"
      },
      {
        "hash": "def5678",
        "message": "fix: correct bug",
        "type": "fix"
      }
    ]
  }
}
```

### Usage Example

```bash
# Generate changelog data
$ verge bump --from 1.2.3 --kind minor --changelog --format json > changelog.json

# Use with changelog tools
$ conventional-changelog -i CHANGELOG.md < changelog.json
```

## Testing

- [ ] Unit test: JSON schema validation
- [ ] Unit test: Changelog metadata is complete
- [ ] Integration test: Output can be used by changelog tools

## Related Tickets

- T004: Output modes (base output)
- T020: `version bump` command (integration)

## Notes

- Keep JSON schema extensible for future enhancements
- Consider compatibility with existing changelog tools
- Document schema for consumers
