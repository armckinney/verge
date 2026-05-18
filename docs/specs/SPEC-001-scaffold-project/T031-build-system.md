# T031: Cross-Platform Build System Setup

**Phase**: 0 - Bootstrap CLI  
**Category**: Build & Distribution  
**Complexity**: Medium  
**Estimated Duration**: 2-3 hours

## Objective

Set up cross-platform build infrastructure using goreleaser to produce binaries for all 5 target platforms.

## Current State

- No cross-platform build system exists
- Manual builds would require platform-specific configuration

## Target State

- Automated builds for all platforms via goreleaser
- Single configuration file `.goreleaser.yaml`
- Reproducible, signed builds
- Release artifacts ready for GitHub Releases distribution

## Acceptance Criteria

- [ ] `.goreleaser.yaml` configuration file created
- [ ] Build targets configured for all 5 platforms:
  - `app-darwin-arm64` (macOS ARM64)
  - `app-linux-amd64` (Linux x86-64)
  - `app-linux-arm64` (Linux ARM64)
  - `app-windows-amd64.exe` (Windows x86-64)
  - `app-windows-arm64.exe` (Windows ARM64)
- [ ] `verge version -V` works and reports version from build metadata
- [ ] Build metadata includes Git commit hash and build timestamp
- [ ] SHA256 checksums generated for all binaries
- [ ] Local build works: `goreleaser build --single-target`
- [ ] Release build produces all artifacts: `goreleaser release --snapshot`
- [ ] Makefile targets for building: `make build`, `make build-snapshot`
- [ ] CI workflow (GitHub Actions) configured to build and publish on git tags

## Context

### Files to Create/Update

- `.goreleaser.yaml` — goreleaser configuration
- `Makefile` — add build targets
- `.github/workflows/release.yml` — GitHub Actions workflow

### Sample `.goreleaser.yaml` Structure

```yaml
version: 2

project_name: verge

builds:
  - id: verge
    main: ./cmd/verge
    binary: verge
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
      - linux
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.Version={{.Version}}
      - -X main.Commit={{.Commit}}
      - -X main.Date={{.Date}}

archives:
  - format: tar.gz
    name_template: >-
      {{- .ProjectName }}_{{- .Version }}_
      {{- if eq .Os "darwin" }}darwin{{ else }}{{ .Os }}{{ end }}_
      {{- .Arch }}
    wrap_in_directory: true
    files:
      - README.md
      - LICENSE

checksum:
  name_template: checksums.txt
  algorithm: sha256

release:
  github:
    owner: your-org
    repo: your-repo
  name_template: "{{.ProjectName}} v{{.Version}}"
```

### Makefile Targets

```makefile
build:
	goreleaser build --single-target

build-snapshot:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist
```

### GitHub Actions Workflow

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      - uses: goreleaser/goreleaser-action@v4
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Testing

- [ ] Build locally on macOS/Linux/Windows (or Docker)
- [ ] Verify binary names match expected format
- [ ] Verify version information with `-V` flag on each platform binary
- [ ] Verify checksums match (sha256sum -c checksums.txt)
- [ ] Test snapshot release build produces all artifacts

## Related Tickets

- T001: Project structure (includes build system setup)
- T005: Error codes (version output format)

## Notes

- CGO_ENABLED=0 ensures static linking (no runtime dependencies)
- Test binaries on actual platforms if possible (use Docker for Linux variants)
- GitHub token is required for release creation (stored in GitHub Actions secrets)
- Consider code signing and notarization for macOS in v2 (Apple notarization)
