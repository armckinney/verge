# Go CLI Specification: Version Management (SemVer + PEP 440)

## 1. Purpose

This document defines a practical plan to evolve this repository from an API template into a Go CLI application for version management.

The CLI will:

- parse and compare versions across SemVer and PEP 440-like forms
- determine the next bump target (major/minor/patch/prerelease)
- retrieve current and latest versions from pluggable sources
- provide stable command contracts so integrations can be replaced or extended later

## 2. Goals and Non-Goals

### Goals

- Support SemVer forms like `1.2.3`, `1.2.3-dev.4`, `1.2.3-rc.2`.
- Support PEP 440-like forms like `1.2.3dev4`, `1.2.3a1`, `1.2.3b2`, `1.2.3rc3`.
- Normalize both formats into one internal representation for ordering and bump logic.
- Expose commands for:
  - current version
  - bump type and next version
  - latest version (global and constrained)
- Support configurable sources (git tags, GitHub Releases, GHCR) with deterministic precedence.
- Use a project config file for repeatable behavior across local/CI/release workflows.

### Non-Goals (initial release)

- Full PEP 440 implementation of epochs (`1!`), post releases (`.post`), and local version segments (`+abc`) unless explicitly needed.
- Automatic publishing/release actions (this CLI will initially compute and fetch versions; publishing can be added later).

### Small-Team Scope Profile (Default)

This project is primarily for personal/small-team use, so v1 should optimize for reliability and low maintenance instead of broad distribution concerns.

Keep in v1:

- One canonical internal version model and deterministic comparison.
- Core commands: `version current`, `version bump`, `version latest`, `version parse`, `version compare`.
- One config format (`.verge.yaml`) with simple overrides.
- Provider set limited to practical sources (`git-tags`, then optional `github-releases`, then optional `ghcr`).
- Strong parser/comparator/bump tests and predictable JSON output.
- Explain mode (`--explain`) to show why a version was selected.

Defer unless needed:

- Plugin frameworks and external provider SDKs.
- Full cross-standard coverage beyond target SemVer + PEP 440-like subset.
- Advanced org policy engines and branch governance DSLs.
- Enterprise-grade telemetry and extensive schema migration systems.
- Deep supply-chain verification workflows.

## 3. Version Format Scope

### 3.1 Accepted input patterns

- SemVer core:
  - `MAJOR.MINOR.PATCH`
  - `MAJOR.MINOR.PATCH-STAGE.SEQ` (e.g., `1.2.3-dev.4`)
- PEP 440-like prerelease/dev:
  - `MAJOR.MINOR.PATCHdevN` (e.g., `1.2.3dev4`)
  - `MAJOR.MINOR.PATCHaN`, `MAJOR.MINOR.PATCHbN`, `MAJOR.MINOR.PATCHrcN`

### 3.2 Internal canonical model

Use one model for ordering and bumping:

```text
Version {
  major: int
  minor: int
  patch: int
  stage: enum(none, dev, alpha, beta, rc)
  sequence: int? // required when stage != none
  metadata: string? // optional, non-ordering
  original: string // original input for diagnostics
  scheme: enum(semver, pep440, auto)
}
```

### 3.3 Ordering rules (initial)

- Compare `major`, then `minor`, then `patch` numerically.
- Final (`stage=none`) is greater than prerelease/dev of same core.
- Stage ordering for same core and same sequence family:
  - `dev < alpha < beta < rc < final`
- For same stage, compare `sequence` numerically.

### 3.4 Rendering rules

Support output styles:

- `semver` output: `1.2.3-rc.4`
- `pep440` output: `1.2.3rc4`
- `auto` output: preserve preferred style from config or input source

## 3.5 Multi-Ecosystem Version Formats

The CLI must normalize and render versions for diverse ecosystems. Each has its own conventions for prefix, stage naming, and sequence representation.

### Ecosystem-specific patterns

**Containers**
- PR: `major.minor.patch-stage.<hash>` (e.g., `1.2.3-dev.a1b2c3d`)
- REL: `major.minor.patch` (e.g., `1.2.3`)
- latest: `latest` tag
- floating: semver-style tags (e.g., `1.2`, `1`)

**Devcontainer Features**
- semver: `major.minor.patch` or `major.minor.patch-stage.seq`
- latest: `latest` tag

**Python (PyPI)**
- PR: `major.minor.patchstageseq` (e.g., `1.2.3dev4`, `1.2.3a1`, `1.2.3b2`, `1.2.3rc3`)
- REL: `major.minor.patch` (e.g., `1.2.3`)

**Go (v-semver)**
- PR: `vmajor.minor.patch-stage.sequence` (e.g., `v1.2.3-dev.4`)
- REL: `vmajor.minor.patch` (e.g., `v1.2.3`)

**Terraform (v-semver)**
- PR: `vmajor.minor.patch-stage.sequence` (e.g., `v1.2.3-rc.2`)
- REL: `vmajor.minor.patch` (e.g., `v1.2.3`)

**GitHub Actions / RWF**
- semver: `major.minor.patch` or `vmajor.minor.patch`
- Prerelease: `major.minor.patch-stage.sequence`

### Component extraction and interpretation

All formats normalize to the canonical model by extracting:

- **prefix**: optional leading character (`v`, empty, or other ecosystem-specific)
- **major**: integer, required
- **minor**: integer, required
- **patch**: integer, required
- **stage**: keyword mapping (`dev` → dev, `a` → alpha, `b` → beta, `rc` → rc, or empty → final)
- **sequence**: numeric or hash, interpreted based on source and config

### Sequence interpretation strategies

Sequence can be:

- **Numeric**: simple integer (e.g., `1`, `42`, `1001`).
- **Commit SHA**: Git commit hash, short form (7-40 chars, e.g., `a1b2c3d`).
- **Content hash**: Hash of file contents or build context (e.g., MD5 or SHA256 short form).
- **Build number**: CI build ID or run number (e.g., `gh-run-12345`).

Interpretation rules (per config and context):

- If sequence is purely numeric, treat as numeric.
- If sequence is alphanumeric (7+ chars, hex-like), treat as short SHA or content hash.
- If sequence matches `gh-` prefix, treat as GitHub build identifier.
- Otherwise, compare lexicographically (fallback).

For ordering purposes:

- Numeric sequences: compare as integers.
- Hash sequences: compare lexicographically (hashes are unique per commit/build).
- Mixed: numeric < hash (final releases sort higher than prerelease builds).

## 4. CLI Command Surface (v1)

Use a root command like `verge` (placeholder name).

### 4.1 `version current`

Purpose: get current project version.

Examples:

```bash
verge version current
verge version current --source git-tags
verge version current --format json
verge version current --ecosystem go
verge version current --ecosystem containers
```

Behavior:

- Detects ecosystem from config or `--ecosystem` flag; falls back to `auto` detection.
- Resolves sources by config precedence unless `--source` is provided.
- Returns one selected version with provenance, rendered in target ecosystem format.
- If no version is found, exits non-zero with actionable error.

Output fields (json):

- `version` (in target ecosystem format)
- `normalizedVersion` (canonical internal form)
- `schemeDetected` (detected format pattern)
- `ecosystem` (target output ecosystem)
- `source` (where version came from)
- `raw` (original unparsed string)
- `sequenceType` (numeric, commit-sha, content-hash, build-id, or unknown)

### 4.2 `version bump`

Purpose: determine bump kind and/or generate next version.

Examples:

```bash
verge version bump --from 1.2.3 --kind minor
verge version bump --from 1.2.3-dev.4 --kind prerelease --stage dev
verge version bump --from 1.2.3rc2 --kind final
```

Supported bump kinds (v1):

- `major`
- `minor`
- `patch`
- `prerelease`
- `final`
- `auto` (optional heuristic from commit messages; see 7.2)

Behavior examples:

- `1.2.3` + `minor` => `1.3.0`
- `1.2.3-dev.4` + `prerelease --stage dev` => `1.2.3-dev.5`
- `1.2.3rc2` + `final` => `1.2.3`

### 4.3 `version latest`

Purpose: fetch latest available version from one or more sources.

Examples:

```bash
verge version latest
verge version latest --source github-releases
verge version latest --constraint "^1.2.3"
verge version latest --core 1.2.3 --stage dev
verge version latest --ecosystem python
verge version latest --ecosystem containers --channel rel
```

Behavior:

- Detects ecosystem from config or `--ecosystem` flag; falls back to `auto` detection.
- Without constraints: highest version across configured source set.
- With `--core 1.2.3 --stage dev`: latest matching prerelease (e.g., `1.2.3-dev.42` for Go, `1.2.3dev42` for Python).
- With `--channel rel` or `--channel pr`: filter to release or prerelease versions only (ecosystem-aware).
- If multiple sources are enabled, tie-break by config precedence and timestamp when needed.
- For container ecosystem: can request `rel`, `pr`, or `floating` variants.

### 4.4 `version parse` (recommended)

Purpose: validate and normalize a version string.

Example:

```bash
verge version parse 1.2.3rc1 --output semver
```

### 4.5 `version compare` (recommended)

Purpose: compare two versions for CI/policy gates.

Example:

```bash
verge version compare 1.2.3 1.2.4
```

Result:

- exit code semantics (`0` equal, `10` left<right, `11` left>right) or explicit json field.

## 5. Source Integration Interfaces

Define provider contracts early so implementations can change safely.

```go
type VersionProvider interface {
    Name() string
    GetCurrent(ctx context.Context, opts QueryOptions) (VersionResult, error)
    GetLatest(ctx context.Context, opts QueryOptions) (VersionResult, error)
    List(ctx context.Context, opts QueryOptions) ([]VersionResult, error)
}
```

Recommended initial providers:

- `git-tags`
  - local or remote tags
  - configurable tag prefix (for example `v`)
  - ecosystem-aware parsing (Go tags vs plain semver)
- `github-releases`
  - release and prerelease support
  - optional include/exclude drafts
  - maps GitHub release names/tags to normalized versions
- `ghcr`
  - image tag listing
  - optional repository/image selector
  - interprets container image tags (PR hashes, REL versions, floating tags)
- `pypi`
  - Python package index integration
  - interprets PEP 440-like prerelease notation
- `terraform-registry`
  - Terraform module/provider version listing
  - v-semver format support

## 6. Configuration File

### 6.1 Format recommendation

Use YAML for initial CLI adoption (`.verge.yaml`) because it is common in CI and already present in the Go ecosystem. TOML support can be added later.

### 6.2 Example config

```yaml
version: 1

# Default ecosystem and format for this project
ecosystem: go           # go | python | containers | terraform | github-actions | auto

format:
  input: auto                    # auto | semver | pep440 | eco-specific
  output: auto                   # auto | semver | pep440 | eco-specific
  tagPrefix: v                   # prefix for git tags and releases
  sequenceInterpreter: auto      # auto | numeric | commit-sha | content-hash | build-id

sources:
  precedence:
    - git-tags
    - github-releases
    - ghcr
    - pypi

  git-tags:
    enabled: true
    fetch: false
    includePrerelease: true
    ecosystemParsing: go         # parse tags as Go v-semver

  github-releases:
    enabled: false
    owner: your-org
    repo: your-repo
    includePrerelease: true
    includeDrafts: false

  ghcr:
    enabled: false
    image: ghcr.io/your-org/your-image
    includePrerelease: true
    channelFilter: null          # filter to 'rel', 'pr', 'floating' or null for all

  pypi:
    enabled: false
    packageName: your-package
    includePrerelease: true

  terraform-registry:
    enabled: false
    module: your-org/your-module/aws

sequence:
  # Interpretation rules for sequence component
  hashLength: 7                  # short SHA length for commit hashes
  allowContentHash: true         # allow MD5/SHA256 file hashes as sequences
  ghBuildPattern: "gh-"         # prefix for GitHub Actions build IDs

rules:
  prereleaseStage: dev
  allowMajorZeroBreaking: true
  defaultBump: patch
  # Ecosystem-specific rules can be added per deployment context

autoBump:
  conventionalCommits: true
  breakingTokens:
    - "BREAKING CHANGE"
    - "!:"
```

### 6.3 Config precedence

1. CLI flags
2. environment variables
3. project config file
4. built-in defaults

## 7. SDLC Use Cases Beyond Basic Commands

### 7.1 Release pipeline use cases

- Validate that candidate version is greater than current published version.
- Resolve next prerelease number for a target core (`1.2.3-dev.N`).
- Ensure branch policies (for example `main` allows final only, `develop` allows dev prereleases).

### 7.2 Automatic bump detection

- Determine bump from commit history:
  - `feat` -> minor
  - `fix` -> patch
  - breaking change marker -> major
- Allow override by explicit flag.

### 7.3 Changelog and artifact traceability

- Emit machine-readable output (`json`) for release notes tooling.
- Include source metadata (tag/release id/digest) to trace why a version was selected.

### 7.4 Consumer/runtime checks

- Compare running app version against latest allowed channel.
- Decide if update is available for stable or prerelease channel.

### 7.5 Explainability and local workflow support

- Add `--explain` mode to show candidate versions, filtering, precedence, and final selection.
- Add optional lightweight cache with TTL for remote sources to reduce API calls.
- Ensure git-tags only workflows work without network access.

## 8. Proposed Code Architecture (Target)

Suggested structure after migration toward CLI:

```text
cmd/verge/main.go
internal/cli/
  root.go
  version_current.go
  version_bump.go
  version_latest.go
  version_parse.go
  version_compare.go
internal/version/
  parse.go
  normalize.go
  compare.go
  bump.go
  render.go
internal/ecosystems/
  registry.go
  types.go
  go.go              # Go v-semver ecosystem
  python.go          # Python PEP440
  containers.go      # Container formats
  terraform.go       # Terraform v-semver
  github_actions.go  # GitHub Actions semver
internal/providers/
  provider.go
  git_tags.go
  github_releases.go
  ghcr.go
  pypi.go
  terraform_registry.go
internal/config/
  load.go
  schema.go
internal/sequence/
  interpreter.go     # numeric, SHA, content-hash, build-id detection
```

Notes:

- `internal/version` contains source-agnostic domain logic (parse, normalize, compare, bump, render).
- `internal/ecosystems` handles ecosystem-specific format rules, prefix handling, and rendering.
- `internal/providers` fetches raw versions and source metadata; delegates parsing to ecosystems.
- `internal/sequence` handles numeric, hash, and build-id interpretation for ordering.
- `internal/cli` orchestrates config + ecosystems + providers + domain logic.

## 9. Testing Strategy

### 9.1 Unit tests

- Parser tests for all accepted forms (SemVer, PEP 440, per-ecosystem formats).
- Comparator table tests for edge ordering cases, including hash/numeric sequences.
- Bump tests for each kind/stage transition.
- Ecosystem-specific rendering tests (Go, Python, containers, Terraform, GitHub Actions).
- Sequence interpretation tests (numeric, SHA, content-hash, build-id detection).

### 9.2 Integration tests

- Fake provider implementations to validate source precedence and latest selection.
- Golden CLI output tests (text and json) for all supported ecosystems.
- Multi-source conflict resolution (same version from git-tags vs github-releases).

### 9.3 Contract tests

- Define provider behavior requirements and run against each implemented provider.
- Ecosystem parsing contract: each provider must correctly map raw tags to normalized versions.
- Keep this lightweight in v1 (no plugin compatibility matrix yet).

### 9.4 Real-world corpus

- Collect golden version strings from actual projects (Go repos, PyPI packages, container images, Terraform modules).
- Use as regression suite for parser and ordering logic.

## 10. Error Handling and UX

- Consistent error codes and messages.
- Actionable suggestions for malformed versions.
- `--format json` for machine use, human-readable by default.
- `--verbose` for source resolution diagnostics.

## 11. Security and Operational Concerns

- Handle GitHub tokens via environment variables (no token in config file by default).
- Respect API rate limits and return clear retry/backoff messages.
- Timeouts and context cancellation for network providers.
- Keep a minimal security baseline for small-team use: no secrets in config, clear token lookup order, and optional cache directory controls.

## 12. Implementation Roadmap

### Phase 0: Bootstrap CLI

- Add root CLI command and subcommand scaffolding.
- Add config loading and validation.
- Add uniform output modes (`text` and `json`) and baseline error codes.

### Phase 1: Core Version Engine

- Implement parse/normalize/compare/render.
- Implement bump logic and exhaustive unit tests.
- Add `version parse` and `version compare` first for confidence in core behavior.

### Phase 2: Local Source Provider

- Implement `git-tags` provider.
- Ship `version current`, `version latest`, `version bump` with local-only behavior.
- Add `--explain` selection trace output.

### Phase 3: Remote Providers

- Add `github-releases` provider.
- Add `ghcr` provider.
- Add lightweight response caching and bounded retries/timeouts.

### Phase 4: CI/Release Enhancements

- Add auto bump from commit history.
- Add changelog-friendly json outputs and optional policy checks.
- Keep advanced policy engine out of scope unless requirements grow.

## 13. Acceptance Criteria (v1)

- CLI can parse and compare SemVer, PEP 440-like, and ecosystem-specific version formats.
- CLI can compute next version for major/minor/patch/prerelease/final.
- CLI can fetch latest and current versions from `git-tags` with ecosystem-aware parsing (at least Go and plain semver).
- CLI can render output in target ecosystem format (Go v-semver, Python PEP 440, containers, Terraform, GitHub Actions).
- CLI supports ecosystem selection (`--ecosystem` flag or config).
- CLI supports config file + flag overrides + environment variables.
- CLI supports human-readable and JSON outputs, with explain mode for selection debugging.
- Unit and integration tests cover parser, ordering, bumping, source precedence, and ecosystem-specific rendering.
- Sequence interpretation (numeric vs hash) works correctly across formats and sources.
- Optional providers (github-releases, ghcr, pypi) parse and normalize correctly when enabled.

## 14. Explicitly Out of Scope for Small-Team v1

- Full PEP 440 feature parity (epoch, post, local labels) beyond PR/REL patterns.
- Dynamic third-party provider plugin architecture.
- Advanced policy DSL and org-wide governance enforcement.
- Signature verification and full supply-chain attestation workflows.
- Support for ecosystems beyond Go, Python, containers, Terraform, GitHub Actions (others can be added incrementally).

## 15. Build and Distribution

### 15.1 Target Platforms (v1)

The CLI shall be built and distributed for the following platforms:

**macOS**
- `app-darwin-arm64` (Apple Silicon, M1/M2/M3)

**Linux**
- `app-linux-amd64` (x86-64 architecture)
- `app-linux-arm64` (ARM 64-bit architecture)

**Windows**
- `app-windows-amd64.exe` (x86-64 architecture)
- `app-windows-arm64.exe` (ARM 64-bit architecture)

### 15.2 Build Process

- Use `goreleaser` for automated multi-platform builds.
- Builds are triggered on git tags (e.g., `v1.0.0`).
- Artifacts are stored with platform and architecture in filename:
  - Format: `verge-<version>-<platform>-<arch>` or `.exe` suffix for Windows.
- Support reproducible builds (same source = same binary hash).

### 15.3 Distribution Channels

**GitHub Releases** (required for v1)
- Upload all platform binaries as release artifacts.
- Include checksums (SHA256) for integrity verification.
- Sign binaries (optional for v1, recommended for v2).

**Package Managers** (future)
- Homebrew (macOS) — `brew install verge`
- Apt/Snap (Linux) — defer to v2
- Chocolatey (Windows) — defer to v2

### 15.4 Installation Methods (v1)

**Manual download**
- Users download binary from GitHub Releases.
- Binary is executable immediately (no installation step).
- Optional: checksum verification with `sha256sum -c checksums.txt`.

**Curl/Shell script** (optional for v1)
- Provide `install.sh` for Unix-like systems.
- Script downloads correct binary based on `uname` (OS and arch detection).

### 15.5 Version Information

- CLI shall report version with `-V` or `--version` flag.
- Output format: `verge version X.Y.Z` (can be parsed by downstream tools).
- Build metadata (Git commit, build timestamp) can be embedded in binary for debugging.

## 16. Open Questions

- Should full PEP 440 (epoch/post/local) be in v1 or v2?
- Is `v` prefix mandatory/optional for tag parsing?
- What should be the default output style (`semver` vs `auto`)?
- Do we need strict compatibility with Python packaging comparison semantics, or only the subset above?
- Should binaries be code-signed (e.g., Apple notarization)? Defer to v2.
