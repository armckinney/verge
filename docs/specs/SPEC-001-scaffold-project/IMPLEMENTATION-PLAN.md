# Semantic Versioning CLI: Implementation Plan

## Overview

This document provides a high-level roadmap for implementing the version management CLI specified in [semver-cli-spec.md](semver-cli-spec.md).

The implementation is organized into **4 phases** with incrementally increasing scope, designed for iterative delivery and confidence building.

## Key Principles

1. **Strong Core First**: Build the version parser, normalizer, comparator, and bump logic with exhaustive test coverage before touching providers.
2. **Local Then Remote**: Start with `git-tags` as the only provider, then add remote sources (GitHub, GHCR, PyPI).
3. **Multi-Ecosystem Ready**: Version and ecosystem abstractions are baked in from Phase 0, but only Go v-semver is fully tested in Phase 1; others scale in Phase 2.
4. **Explain-Driven Development**: Add `--explain` mode early to surface reasoning about version selection and bump logic.
5. **Small-Team Scope**: All features are pragmatic for a single team; defer enterprise-grade features (plugins, full attestation, complex policies).- **Cross-Platform Ready**: Build and distribution system must support 5 target platforms (macOS ARM64, Linux x86-64/ARM64, Windows x86-64/ARM64) from Phase 0.
## Phase 0: Bootstrap CLI (Foundational Setup)

**Goal**: Project structure, config system, and command scaffolding.

**Duration estimate**: 2-3 days

**Deliverables**:
- Migrate from API template to CLI structure
- Config file loading and validation
- Root CLI and subcommand scaffolding
- Uniform output modes (text / JSON)
- Error code taxonomy

**Key tickets**: T001–T005, T031 (build system)

**Exit criteria**:
- `verge version` displays help
- `verge parse --help` shows available options
- `.verge.yaml` loads without errors
- JSON and text output modes work for all commands
- Build system is set up and can produce binaries for all 5 target platforms (darwin-arm64, linux-amd64, linux-arm64, windows-amd64, windows-arm64)

---

## Phase 1: Core Version Engine (Domain Logic)

**Goal**: Robust version parsing, normalization, comparison, and bump logic with multi-ecosystem support.

**Duration estimate**: 4-5 days

**Deliverables**:
- Canonical version model
- Version parser for SemVer, PEP 440, and all five ecosystems
- Version normalizer
- Version comparator (deterministic ordering)
- Version bump logic
- Version renderer per ecosystem
- `version parse` command
- `version compare` command
- Exhaustive unit tests and golden test corpus

**Key tickets**: T006–T014

**Dependencies**: Phase 0 completed

**Exit criteria**:
- Parser handles all SemVer and PEP 440-like input patterns
- Comparator handles numeric/hash sequences and mixed orderings
- All ecosystem renderers produce correct output
- 95%+ test coverage on version logic
- Golden test corpus passes all assertions

---

## Phase 2: Local Source Provider & Core Commands

**Goal**: Local version retrieval and core commands with explain mode.

**Duration estimate**: 3-4 days

**Deliverables**:
- Provider interface and contract definition
- `git-tags` provider implementation
- Sequence interpreter (numeric vs hash detection)
- `version current` command (git-tags only)
- `version latest` command (git-tags only)
- `version bump` command
- `--explain` mode for selection transparency
- Integration tests for provider + commands
- Real-world git repository test cases

**Key tickets**: T015–T022

**Dependencies**: Phase 1 completed

**Exit criteria**:
- `verge current --source git-tags` returns correct version
- `verge latest --source git-tags --explain` shows candidate filtering
- `verge bump --from 1.2.3 --kind minor` computes correct bump
- All git-tag real-world test cases pass
- Provider contract tests pass

---

## Phase 3: Remote Providers & Caching

**Goal**: GitHub Releases, GHCR, PyPI support with lightweight caching.

**Duration estimate**: 3-4 days

**Deliverables**:
- `github-releases` provider with auth token support
- `ghcr` provider with image tag parsing
- `pypi` provider (optional, deferred if time-limited)
- Lightweight cache layer with TTL and fingerprinting
- Retry + timeout logic
- Multi-source precedence resolution
- Integration tests across all providers

**Key tickets**: T023–T027

**Dependencies**: Phase 2 completed

**Exit criteria**:
- `verge latest --source github-releases` returns correct version
- `verge latest --source ghcr` returns correct container tag
- Caching reduces redundant API calls
- Rate-limit errors are handled gracefully
- Multi-source conflict resolution works (precedence order respected)

---

## Phase 4: CI/Release Enhancements

**Goal**: Conventional commits parsing, policy checks, and release-ready outputs.

**Duration estimate**: 2-3 days

**Deliverables**:
- Conventional commits parser for auto bump detection
- Policy checks (optional, can be minimal)
- Changelog-friendly JSON outputs
- Integration tests for auto bump

**Key tickets**: T028–T030

**Dependencies**: Phase 3 completed

**Exit criteria**:
- `verge bump --auto` reads commit history and suggests correct bump
- Policy checks (if implemented) validate version against branch rules
- JSON output includes changelog metadata
- Real-world repo workflows pass

---

## Milestone Timeline

| Phase | Tickets | Estimated Duration | Key Deliverable |
|-------|---------|---------------------|-----------------|
| 0     | T001–T005, T031 | 2–3 days | Bootstrapped CLI with config and build system |
| 1     | T006–T014 | 4–5 days | Robust version engine + commands |
| 2     | T015–T022 | 3–4 days | Local git provider + explain mode |
| 3     | T023–T027 | 3–4 days | Remote providers + caching |
| 4     | T028–T030 | 2–3 days | CI integrations + polish |
| **Total** | **31 tickets** | **14–19 days** | **v1 Release** |

---

## Ticket Categories by Domain

### Bootstrap & Build
- T001: Project structure migration
- T002: Config loading and schema
- T003: CLI scaffolding
- T004: Output modes
- T005: Error codes
- T031: Cross-platform build system

### Core Version Domain
- T006: Version model definition
- T007: Parser implementation
- T008: Normalizer implementation
- T009: Comparator implementation
- T010: Bump logic implementation
- T011: Renderer implementation
- T012: `version parse` command
- T013: `version compare` command
- T014: Core tests and corpus

### Local Provider
- T015: Provider interface
- T016: Git-tags provider
- T017: Sequence interpreter
- T018: `version current` command
- T019: `version latest` command
- T020: Explain mode
- T021: Integration tests
- T022: Golden tests

### Remote Providers
- T023: GitHub Releases provider
- T024: GHCR provider
- T025: Caching layer
- T026: Retries and timeouts
- T027: Provider integration tests

### Release Features
- T028: Conventional commits parser
- T029: Policy checks
- T030: Changelog outputs

---

## Build and Distribution

### Cross-Platform Targets (v1)

The CLI is built and distributed for 5 platforms:
- **macOS ARM64** (`app-darwin-arm64`) — Apple Silicon M1/M2/M3
- **Linux x86-64** (`app-linux-amd64`) — Standard desktop/server Linux
- **Linux ARM64** (`app-linux-arm64`) — ARM-based Linux (Raspberry Pi, AWS Graviton)
- **Windows x86-64** (`app-windows-amd64.exe`) — Standard Windows desktop/server
- **Windows ARM64** (`app-windows-arm64.exe`) — ARM-based Windows

### Build System

**Tooling**: goreleaser for automated cross-platform builds
**Configuration**: `.goreleaser.yaml` (created in T031)
**Trigger**: Automated builds on git tags (e.g., `v1.0.0`)
**Artifacts**: Binaries, checksums (SHA256), release notes

### Distribution Channels (v1)

**GitHub Releases** (required)
- All platform binaries uploaded as release assets
- Checksums file for integrity verification
- Installation via direct download

**Optional installation methods** (v1+)
- Shell script download: `curl | bash` style installation
- Homebrew for macOS (v2+)
- Package managers for Linux (v2+, defer to ecosystem maintainers)

### Version Information

- CLI reports build metadata via `verge info`
- Output includes version, commit, and date fields
- Build metadata (Git commit, timestamp) embedded and visible with `--verbose`
- Version matches git tag for reproducibility

**Related tickets**: T031 (build system setup)

---

### Testing Strategy
- **Unit**: Parser, normalizer, comparator, bump logic → Phase 1
- **Integration**: Provider + command orchestration → Phase 2
- **Contract**: Provider behavior requirements → Phase 2, 3
- **Golden**: Real-world version strings → Phase 2, embedded in later phases
- **End-to-End**: Full CLI workflows → Phase 3+

### Documentation
- Inline code comments for complex logic (especially sequence interpretation)
- README updates for new commands and ecosystem support
- Config schema documentation (inline in config/schema.go)
- Command help text with ecosystem examples

### Quality Gates
- All phases must maintain >90% test coverage on production code
- No security issues in token/credential handling
- Graceful degradation for network failures
- Clear, actionable error messages for all failure modes

---

## How to Use This Plan

1. **Review [semver-cli-spec.md](semver-cli-spec.md)** for full requirements and rationale.
2. **Start with Phase 0**: Run tickets T001–T005 sequentially to establish foundations.
3. **Proceed to Phase 1**: T006–T014 focus on correctness; this phase is the safety net for later phases.
4. **Parallelize where possible**: Phase 1 tests (T014) and Phase 2 integration (T021) can run concurrently with implementation.
5. **Adjust scope as needed**: If remote providers are not needed immediately, Phase 3 can be deferred; core CLI is complete after Phase 2.

---

## Success Criteria (v1)

- ✅ CLI builds and runs
- ✅ Cross-platform binaries produced for all 5 targets (darwin-arm64, linux-amd64, linux-arm64, windows-amd64, windows-arm64)
- ✅ All 31 tickets completed with passing tests
- ✅ Core commands (`parse`, `compare`, `current`, `latest`, `bump`) work with git-tags
- ✅ Multi-ecosystem support verified for Go, Python, containers, Terraform, GitHub Actions
- ✅ Explain mode provides clear reasoning
- ✅ Config file + flag overrides work correctly
- ✅ JSON and text outputs are correct and stable
- ✅ Real-world test corpus (50+ version strings) passes
- ✅ No blocker bugs in Phase 2+ gate criteria

---

## Related Documents

- [semver-cli-spec.md](semver-cli-spec.md) — Detailed feature specification
- [implementation-tickets/](implementation-tickets/) — Individual ticket descriptions with context and acceptance criteria
