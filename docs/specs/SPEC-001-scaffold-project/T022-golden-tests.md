# T022: Golden Tests and Real-World Version Validation

**Phase**: 2 - Local Source Provider  
**Category**: Testing  
**Complexity**: Medium  
**Estimated Duration**: 2-3 hours

## Objective

Establish golden tests using real-world version strings from major projects to ensure robustness across diverse formats.

## Current State

- No golden tests exist for real-world versions

## Target State

- 100+ real-world version strings are tested
- Golden test corpus covers all supported ecosystems
- Parser, normalizer, comparator, renderer are all tested end-to-end
- Test failures indicate regressions immediately

## Acceptance Criteria

- [ ] `tests/fixtures/golden_corpus.go` contains:
  - Real Go repository versions (20+)
  - Real Python package versions (20+)
  - Real container image versions (20+)
  - Real Terraform versions (15+)
  - Real GitHub release versions (15+)
- [ ] `tests/golden_test.go` with table-driven tests:
  - Parse each version string
  - Verify parsing succeeds
  - Render to all ecosystems
  - Verify rendering is valid
  - Compare versions for sanity
- [ ] Coverage includes:
  - Normal releases
  - Prerelease versions
  - Development versions
  - Hash-based sequences
  - Edge cases (leading zeros, very large numbers)
- [ ] Test execution on CI produces golden artifacts for comparison

## Context

### Files to Create

- `tests/fixtures/golden_corpus.go` — real-world versions
- `tests/golden_test.go` — validation tests

### Golden Corpus Template

```go
package fixtures

type GoldenVersion struct {
    Input       string
    Ecosystem   string
    Description string
    ExpectedErr bool
}

var GoldenVersions = []GoldenVersion{
    // Go
    {Input: "v1.28.0", Ecosystem: "go", Description: "kubernetes/kubernetes"},
    {Input: "v1.28.0-rc.1", Ecosystem: "go", Description: "kubernetes/kubernetes prerelease"},
    {Input: "v2.40.1", Ecosystem: "go", Description: "docker/cli"},
    {Input: "v1.0.0-dev.abc1234", Ecosystem: "go", Description: "custom dev version"},
    
    // Python
    {Input: "1.0.0", Ecosystem: "python", Description: "pip/setuptools"},
    {Input: "1.0.0dev1", Ecosystem: "python", Description: "pip dev version"},
    {Input: "2.0.0a1", Ecosystem: "python", Description: "django alpha"},
    {Input: "3.11.0rc1", Ecosystem: "python", Description: "python rc"},
    
    // Containers
    {Input: "1.2.3", Ecosystem: "containers", Description: "docker hub"},
    {Input: "1.2.3-dev.a1b2c3d", Ecosystem: "containers", Description: "PR version"},
    {Input: "latest", Ecosystem: "containers", Description: "docker latest tag"},
    {Input: "1.2", Ecosystem: "containers", Description: "floating semver"},
    
    // Terraform
    {Input: "v4.5.0", Ecosystem: "terraform", Description: "aws provider"},
    {Input: "v4.5.0-rc1", Ecosystem: "terraform", Description: "aws rc"},
    
    // GitHub Actions
    {Input: "1.0.0", Ecosystem: "github-actions", Description: "action release"},
    {Input: "v1.0.0", Ecosystem: "github-actions", Description: "action with v prefix"},
}

// Invalid versions to test rejection
var InvalidVersions = []struct{
    Input string
    Reason string
}{
    {"1.2.3.4", "too many components"},
    {"1.x.3", "non-numeric minor"},
    {"1.2.3-xxx.4", "invalid stage"},
}
```

### Golden Test Implementation

```go
func TestGoldenVersions(t *testing.T) {
    for _, golden := range fixtures.GoldenVersions {
        t.Run(golden.Description, func(t *testing.T) {
            // Parse
            parsed, err := version.DefaultParser.Parse(golden.Input)
            if err != nil {
                t.Fatalf("parse error: %v", err)
            }
            
            // Normalize
            normalized, err := version.DefaultNormalizer.Normalize(parsed)
            if err != nil {
                t.Fatalf("normalize error: %v", err)
            }
            
            // Render
            rendered, err := version.DefaultRenderer.Render(normalized, golden.Ecosystem)
            if err != nil {
                t.Fatalf("render error: %v", err)
            }
            
            // Re-parse rendered version
            reparsed, err := version.DefaultParser.Parse(rendered)
            if err != nil {
                t.Fatalf("reparse error: %v", err)
            }
            
            // Verify round-trip
            if version.DefaultComparator.Compare(parsed, reparsed) != 0 {
                t.Errorf("round-trip failed: %s -> %s -> %s", golden.Input, rendered, reparsed)
            }
        })
    }
}

func TestInvalidVersions(t *testing.T) {
    for _, invalid := range fixtures.InvalidVersions {
        t.Run(invalid.Reason, func(t *testing.T) {
            _, err := version.DefaultParser.Parse(invalid.Input)
            if err == nil {
                t.Errorf("expected error for %q, got none", invalid.Input)
            }
        })
    }
}
```

### Golden Corpus Sources

Collect real versions from:
- **Go**: github.com/kubernetes/kubernetes, golang/go, docker/moby
- **Python**: pypi.org (setuptools, pip, requests, django)
- **Containers**: docker hub, gcr.io, docker.io
- **Terraform**: registry.terraform.io (aws, google, azure providers)
- **GitHub**: Major project releases (note versioning patterns)

## Testing

- [ ] Unit test: All golden versions parse successfully
- [ ] Unit test: All golden versions normalize correctly
- [ ] Unit test: All golden versions render to correct ecosystem format
- [ ] Unit test: Round-trip (parse → render → parse) preserves semantics
- [ ] Unit test: Invalid versions are rejected
- [ ] Performance: All 100+ versions parse in <100ms total

## Related Tickets

- T007: Parser (tested)
- T008: Normalizer (tested)
- T011: Renderer (tested)
- T014: Core tests (builds on this)

## Notes

- Add new golden versions as edge cases are discovered
- Document source for each golden version (link to actual project)
- Consider using git tags from live repos in CI for ongoing validation
