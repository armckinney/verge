# T013: Implement `version compare` Command

**Phase**: 1 - Core Version Engine  
**Category**: CLI Commands  
**Complexity**: Low  
**Estimated Duration**: 2-3 hours

## Objective

Implement the `version compare` command to compare two versions for CI/policy gates.

## Current State

- Stub command exists in CLI scaffolding

## Target State

- `verge compare <version1> <version2>` returns comparison result
- Exit codes are machine-friendly for CI integration
- Text and JSON output modes supported
- Output includes human-readable and machine-readable comparison

## Acceptance Criteria

- [ ] `internal/cli/version_compare.go` implements full logic
- [ ] Command signature: `verge compare <version1> <version2> [flags]`
- [ ] `--format` flag with choices: text, json
- [ ] Exit codes:
  - `0` if equal
  - `10` if left < right
  - `11` if left > right
  - `1` on error (invalid version)
- [ ] Text output shows comparison result (left < right, left == right, left > right)
- [ ] JSON output includes comparison field and exit code
- [ ] Clear error for invalid versions with suggestions

## Context

### Files to Update

- `internal/cli/version_compare.go` — implement compare command

### Command Behavior

```bash
# Equal versions
$ verge compare 1.2.3 1.2.3
Comparison: 1.2.3 == 1.2.3
$ echo $?
0

# Left < Right
$ verge compare 1.2.3 1.2.4
Comparison: 1.2.3 < 1.2.4
$ echo $?
10

# Left > Right
$ verge compare 1.2.4 1.2.3
Comparison: 1.2.4 > 1.2.3
$ echo $?
11

# JSON output
$ verge compare 1.2.3 1.2.4 --format json
{
  "left": {
    "version": "1.2.3",
    "normalized": "1.2.3"
  },
  "right": {
    "version": "1.2.4",
    "normalized": "1.2.4"
  },
  "comparison": "left<right",
  "exitCode": 10
}
```

### Implementation Pattern

```go
var compareCmd = &cobra.Command{
    Use:   "compare <version1> <version2>",
    Short: "Compare two versions",
    Long:  `Compare two versions and return exit code for CI integration.`,
    Args:  cobra.ExactArgs(2),
    RunE:  runVersionCompare,
}

func runVersionCompare(cmd *cobra.Command, args []string) error {
    left, err := parser.Parse(args[0])
    if err != nil {
        return fmt.Errorf("parse error on left: %w", err)
    }
    
    right, err := parser.Parse(args[1])
    if err != nil {
        return fmt.Errorf("parse error on right: %w", err)
    }
    
    leftNorm, _ := normalizer.Normalize(left)
    rightNorm, _ := normalizer.Normalize(right)
    
    result := comparator.Compare(leftNorm, rightNorm)
    
    // Format output
    formatter.FormatComparison(leftNorm, rightNorm, result)
    
    // Set exit code
    switch result {
    case -1:
        os.Exit(10)
    case 1:
        os.Exit(11)
    }
    return nil
}
```

### Usage Examples

Use cases for CI/CD:

```bash
# Check if update is available
NEW_VERSION=$(verge latest)
CURRENT_VERSION=$(verge current)
if verge compare "$CURRENT_VERSION" "$NEW_VERSION" > /dev/null 2>&1; then
    # Versions are equal
    echo "Already at latest"
else
    if [ $? -eq 10 ]; then
        echo "Update available: $CURRENT_VERSION -> $NEW_VERSION"
    fi
fi

# Validate release version
if verge compare "$RELEASE_VERSION" "$PREVIOUS_VERSION"; then
    echo "Release version must be higher than previous"
    exit 1
fi
```

## Testing

- [ ] Unit test: Equal versions return 0
- [ ] Unit test: Left < Right returns exit code 10
- [ ] Unit test: Left > Right returns exit code 11
- [ ] Unit test: Invalid versions return error
- [ ] Unit test: JSON output is valid and includes exit code
- [ ] Integration test: `verge compare 1.2.3 1.2.4 && echo "equal"` works
- [ ] Integration test: Exit codes are correct for CI scripts

## Related Tickets

- T003: CLI scaffolding
- T004: Output modes
- T005: Error codes
- T009: Comparator implementation (dependency)

## Notes

- Exit codes 10/11 are chosen to not conflict with standard Unix codes
- Comparison should work across all formats (SemVer, PEP440, etc.)
- Consider adding `--explain` flag in Phase 2 for debugging
