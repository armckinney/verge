# T011: Implement Version Renderer for All Ecosystems

**Phase**: 1 - Core Version Engine  
**Category**: Domain Logic  
**Complexity**: High  
**Estimated Duration**: 5-6 hours

## Objective

Implement version rendering for all five supported ecosystems (Go, Python, Containers, Terraform, GitHub Actions) with correct prefix and stage naming conventions.

## Current State

- No renderer exists

## Target State

- `Renderer` interface converts canonical Version to ecosystem-specific string
- All five ecosystems render correctly
- Prefix handling is automatic
- Stage names use ecosystem conventions
- Rendering is deterministic and reversible (can re-parse result)

## Acceptance Criteria

- [ ] `internal/version/render.go` implements `Renderer` interface
- [ ] Ecosystem renderers exist in `internal/ecosystems/`:
  - `go.go` — Go v-semver with v prefix
  - `python.go` — Python PEP440-like without dash
  - `containers.go` — Container format (plain or hashed)
  - `terraform.go` — Terraform v-semver with v prefix
  - `github_actions.go` — GitHub Actions semver
- [ ] Go rendering:
  - Final: `v1.2.3`
  - Prerelease: `v1.2.3-dev.4`, `v1.2.3-rc.2`
- [ ] Python rendering:
  - Final: `1.2.3`
  - Prerelease: `1.2.3dev4`, `1.2.3a1`, `1.2.3b2`, `1.2.3rc3`
- [ ] Container rendering:
  - Release: `1.2.3`
  - PR: `1.2.3-dev.a1b2c3d`
  - Floating: `1.2`, `1`
- [ ] Terraform rendering:
  - Final: `v1.2.3`
  - Prerelease: `v1.2.3-rc.2`
- [ ] GitHub Actions rendering:
  - Final: `1.2.3` or `v1.2.3`
  - Prerelease: `1.2.3-dev.4`
- [ ] All renderers handle null/zero sequences gracefully

## Context

### Files to Create

- `internal/version/render.go` — main renderer
- `internal/ecosystems/registry.go` — ecosystem registry
- `internal/ecosystems/types.go` — ecosystem interface
- `internal/ecosystems/go.go` — Go v-semver
- `internal/ecosystems/python.go` — Python PEP440
- `internal/ecosystems/containers.go` — Container formats
- `internal/ecosystems/terraform.go` — Terraform v-semver
- `internal/ecosystems/github_actions.go` — GitHub Actions
- `internal/version/render_test.go` — tests

### Ecosystem Specifications

**Go v-semver**
- Prefix: `v`
- Stage names: `dev`, `alpha`, `beta`, `rc`
- Separator: `-` (dash) before stage
- Example: `v1.2.3-dev.4`, `v1.2.3-rc.2`

**Python PEP440**
- Prefix: none
- Stage names: `dev`, `a`, `b`, `rc`
- Separator: none (stage immediately after patch)
- Example: `1.2.3dev4`, `1.2.3a1`, `1.2.3b2`, `1.2.3rc3`

**Containers**
- PR format: `major.minor.patch-stage.hash`
- Release format: `major.minor.patch`
- Latest: `latest`
- Floating: `major.minor`, `major`
- Example: `1.2.3-dev.a1b2c3d`, `1.2.3`

**Terraform v-semver**
- Prefix: `v`
- Stage names: `dev`, `alpha`, `beta`, `rc`
- Separator: `-` (dash) before stage
- Example: `v1.2.3-rc.2`, `v1.2.3`

**GitHub Actions**
- Prefix: optional `v`
- Stage names: `dev`, `alpha`, `beta`, `rc`
- Separator: `-` (dash) before stage
- Example: `1.2.3-dev.4`, `v1.2.3-rc.2`

### Renderer Interface Design

```go
type Renderer interface {
    // Render converts Version to ecosystem-specific string
    Render(v *Version) (string, error)
    
    // CanRender checks if Version can be rendered (e.g., hash sequences in some ecosystems)
    CanRender(v *Version) bool
}

type renderer struct {
    ecosystems map[string]Renderer
}

func (r *renderer) Render(v *Version, ecosystem string) (string, error) {
    eco, ok := r.ecosystems[ecosystem]
    if !ok {
        return "", fmt.Errorf("unknown ecosystem: %s", ecosystem)
    }
    return eco.Render(v)
}

// Each ecosystem implements Renderer
type goRenderer struct{}

func (g *goRenderer) Render(v *Version) (string, error) {
    if v.Stage == StageFinal {
        return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch), nil
    }
    stageName := stageToGo(v.Stage)
    return fmt.Sprintf("v%d.%d.%d-%s.%v", v.Major, v.Minor, v.Patch, stageName, v.Sequence), nil
}

type pythonRenderer struct{}

func (p *pythonRenderer) Render(v *Version) (string, error) {
    if v.Stage == StageFinal {
        return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch), nil
    }
    stageName := stageToPython(v.Stage)
    return fmt.Sprintf("%d.%d.%d%s%v", v.Major, v.Minor, v.Patch, stageName, v.Sequence), nil
}

// Stage name mappings
func stageToGo(s Stage) string {
    switch s {
    case StageDev:
        return "dev"
    case StageAlpha:
        return "alpha"
    case StageBeta:
        return "beta"
    case StageRC:
        return "rc"
    }
    return ""
}

func stageToPython(s Stage) string {
    switch s {
    case StageDev:
        return "dev"
    case StageAlpha:
        return "a"
    case StageBeta:
        return "b"
    case StageRC:
        return "rc"
    }
    return ""
}
```

### Test Corpus

```go
[]struct {
    version  string
    ecosystem string
    expected string
}{
    {"1.2.3", "go", "v1.2.3"},
    {"1.2.3-dev.4", "go", "v1.2.3-dev.4"},
    {"1.2.3-rc.2", "go", "v1.2.3-rc.2"},
    
    {"1.2.3", "python", "1.2.3"},
    {"1.2.3dev4", "python", "1.2.3dev4"},
    {"1.2.3a1", "python", "1.2.3a1"},
    {"1.2.3b2", "python", "1.2.3b2"},
    {"1.2.3rc3", "python", "1.2.3rc3"},
    
    {"1.2.3", "containers", "1.2.3"},
    {"1.2.3-dev.a1b2c3d", "containers", "1.2.3-dev.a1b2c3d"},
    
    {"1.2.3", "terraform", "v1.2.3"},
    {"1.2.3-rc.2", "terraform", "v1.2.3-rc.2"},
    
    {"1.2.3", "github-actions", "1.2.3"},
    {"1.2.3-dev.4", "github-actions", "1.2.3-dev.4"},
}
```

## Testing

- [ ] Unit test: All ecosystems render correctly for final releases
- [ ] Unit test: All ecosystems render correctly for prerelease versions
- [ ] Unit test: Stage name mapping is correct per ecosystem
- [ ] Unit test: Prefix handling is correct per ecosystem
- [ ] Unit test: 50+ render scenarios in table-driven format
- [ ] Integration test: Rendered string can be re-parsed to same canonical form

## Related Tickets

- T006: Version model
- T007: Parser (round-trip: parse → render → parse)
- T012: `version parse` command (uses renderer)
- T018: `version current` command (uses renderer)

## Notes

- Ensure rendering is deterministic (same input always produces same output)
- Test round-trip: parse input → render → parse again (should match original)
- Consider auto-detection of ecosystem in Phase 2
