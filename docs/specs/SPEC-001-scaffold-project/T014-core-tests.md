# T014: Comprehensive Unit Tests and Golden Test Corpus

**Phase**: 1 - Core Version Engine  
**Category**: Testing  
**Complexity**: High  
**Estimated Duration**: 4-5 hours

## Objective

Establish comprehensive test coverage for all Phase 1 components (parser, normalizer, comparator, bumper, renderer) with a golden test corpus of real-world version strings.

## Current State

- Individual ticket test cases exist
- No comprehensive test suite exists

## Target State

- All Phase 1 code has >95% test coverage
- Golden test corpus includes 100+ real-world version strings
- Table-driven tests cover all documented behaviors
- Test corpus includes edge cases, errors, and ecosystem-specific patterns
- CI gates enforce coverage requirements

## Acceptance Criteria

- [ ] Test files created:
  - `internal/version/parse_test.go` with 50+ test cases
  - `internal/version/normalize_test.go` with 30+ test cases
  - `internal/version/compare_test.go` with 100+ test cases
  - `internal/version/bump_test.go` with 40+ test cases
  - `internal/version/render_test.go` with 50+ test cases
- [ ] `tests/fixtures/golden_versions.go` contains:
  - Real Go repository tags (kubernetes, docker, etc.)
  - Real Python PyPI versions
  - Real container image tags
  - Real Terraform module versions
  - Real GitHub release versions
- [ ] Coverage metrics:
  - `internal/version/*` >95% coverage
  - `internal/ecosystems/*` >90% coverage
- [ ] All test tables include:
  - Input
  - Expected output or behavior
  - Error case (if applicable)
  - Description or comment
- [ ] Benchmark tests for performance-critical functions (parser, comparator)
- [ ] `make test` runs all tests and generates coverage report

## Context

### Golden Test Corpus Structure

```go
// tests/fixtures/golden_versions.go
package fixtures

var GoVersions = []Version{
    // From kubernetes/kubernetes
    {Input: "v1.28.0", Expected: "1.28.0", Ecosystem: "go"},
    {Input: "v1.28.0-rc.1", Expected: "1.28.0-rc.1", Ecosystem: "go"},
    {Input: "v1.27.8", Expected: "1.27.8", Ecosystem: "go"},
    // ... more real examples
}

var PythonVersions = []Version{
    // From pip/setuptools
    {Input: "1.0.0", Expected: "1.0.0", Ecosystem: "python"},
    {Input: "1.0.0dev1", Expected: "1.0.0dev1", Ecosystem: "python"},
    {Input: "1.0.0a1", Expected: "1.0.0a1", Ecosystem: "python"},
    {Input: "2.0.0rc1", Expected: "2.0.0rc1", Ecosystem: "python"},
    // ... more real examples
}

var ContainerVersions = []Version{
    // From Docker Hub
    {Input: "1.2.3", Expected: "1.2.3", Ecosystem: "containers"},
    {Input: "1.2.3-dev.abc1234", Expected: "1.2.3-dev.abc1234", Ecosystem: "containers"},
    {Input: "latest", Expected: "latest", Ecosystem: "containers"},
    // ... more real examples
}

// And similar for TerraformVersions, GitHubVersions
```

### Test File Template

```go
// internal/version/parse_test.go
package version

import "testing"

func TestParseValidSemVer(t *testing.T) {
    tests := []struct {
        input    string
        expected *Version
        name     string
    }{
        {
            name:  "simple final release",
            input: "1.2.3",
            expected: &Version{Major: 1, Minor: 2, Patch: 3, Stage: StageFinal},
        },
        {
            name:  "with dev prerelease",
            input: "1.2.3-dev.4",
            expected: &Version{Major: 1, Minor: 2, Patch: 3, Stage: StageDev, Sequence: 4},
        },
        // ... more cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := defaultParser.Parse(tt.input)
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if result.Major != tt.expected.Major {
                t.Errorf("major: got %d, want %d", result.Major, tt.expected.Major)
            }
            // ... more assertions
        })
    }
}

func TestParseInvalid(t *testing.T) {
    tests := []struct {
        input string
        name  string
    }{
        {"1.2.3.4", "too many components"},
        {"1.x.3", "non-numeric minor"},
        {"1.2.3-xxx.4", "invalid stage"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := defaultParser.Parse(tt.input)
            if err == nil {
                t.Errorf("expected error, got nil")
            }
        })
    }
}

func BenchmarkParse(b *testing.B) {
    input := "v1.2.3-rc.4"
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        defaultParser.Parse(input)
    }
}
```

### Coverage Report Script

Add to Makefile:

```makefile
coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
	@go tool cover -func=coverage.out | tail -1
```

## Testing

- [ ] Unit test: All parser test cases pass
- [ ] Unit test: All normalizer test cases pass
- [ ] Unit test: All comparator test cases pass (especially transitivity)
- [ ] Unit test: All bumper test cases pass
- [ ] Unit test: All renderer test cases pass
- [ ] Benchmark: Parser performance is acceptable (< 1ms per parse)
- [ ] Benchmark: Comparator performance is acceptable (< 1µs per compare)
- [ ] Coverage: >95% coverage on version logic
- [ ] Coverage: >90% coverage on ecosystems logic
- [ ] Golden corpus: All real-world versions parse and render correctly

## Related Tickets

- T007: Parser (tests)
- T008: Normalizer (tests)
- T009: Comparator (tests)
- T010: Bumper (tests)
- T011: Renderer (tests)
- T012: Parse command (integration tests)
- T013: Compare command (integration tests)

## Test Corpus Sources

Collect real versions from:
- GitHub: kubernetes/kubernetes, moby/moby, hashicorp/terraform, etc.
- PyPI: setuptools, pip, requests, django, etc.
- Docker Hub: golang, node, python, etc.
- Terraform Registry: aws, google, azure providers
- GitHub Releases: major projects

## Notes

- Keep test tables readable with clear naming
- Include edge cases: leading zeros, very large numbers, zero sequences
- Test round-trip: parse → normalize → render → parse (should match)
- Benchmark regressions should trigger CI failure
- Document test corpus sources in comments
