# Verge Architecture Overview

## Project Summary

**Verge** is a semantic versioning CLI tool designed to parse, compare, bump, and query versions across multiple ecosystems (Go, Python, Terraform, Containers, GitHub Actions). It provides a unified interface for version management with support for multiple version format schemes (semver, v-semver, pep440) and multiple sources (Git tags, GitHub releases, container registries).

---

## Core Concepts

### Version Model

A `Version` in Verge is a structured representation of a semantic version with the following components:

- **Major, Minor, Patch**: Core semantic version numbers
- **Stage**: Release stage (e.g., alpha, beta, rc, final)
- **Sequence**: Build metadata (commit SHA, git height, content hash, etc.)
- **SequenceType**: Classification of the sequence (numeric, commit SHA, content hash, build ID, etc.)
- **Original**: Raw string representation

The version model is ecosystem-agnostic, allowing versions from different sources to be compared and normalized.

### Ecosystem Abstraction

Verge abstracts version format differences across ecosystems through registry patterns. Each ecosystem (Go, Python, Terraform) defines:

- **Format schemes**: How versions are parsed and rendered for that ecosystem
- **Renderers**: Functions that format a `Version` into an ecosystem-specific string
- **Parsing rules**: How to convert ecosystem-specific version strings to the generic `Version` model

### Version Sources (Providers)

Verge queries versions from multiple sources:

- **Git Tags**: Local or remote Git repository tags
- **GitHub Releases**: GitHub API for official releases
- **GHCR**: GitHub Container Registry for container images
- **Extensible**: New providers can be added following the `VersionProvider` interface

### Configuration-Driven Behavior

Verge uses YAML-based configuration (`.verge.yaml`) to define:

- Which ecosystems and formats to use
- Which sources to query and their precedence
- Sequence interpretation rules (e.g., commit SHA length)
- Bump policies and rules
- Auto-bump configurations

---

## Architecture Layers

### 1. CLI Layer (`internal/cli/`)

**Responsibility**: Command-line interface and user interaction

- **Root command** (`root.go`): Sets up global flags, registers subcommands
- **Version command** (`version.go`): Parent command for version operations
- **Subcommands**:
  - `version parse`: Parse and normalize a version string
  - `version compare`: Compare two versions
  - `version bump`: Bump a version according to rules
  - `version current`: Get current version from configured sources
  - `version latest`: Get latest version from configured sources
  - `version info`: Display version information
- **Output handlers** (`output.go`, `output_changelog.go`): Format results as text or JSON
- **Error handling** (`errors.go`): Structured error codes and messages

### 2. Configuration Layer (`internal/config/`)

**Responsibility**: Load and validate configuration

- **Schema** (`schema.go`): Configuration structure and supported options
- **Loader** (`load.go`): Parse `.verge.yaml` files and apply defaults
- **Defaults** (`defaults.go`): Default configuration values for different ecosystems

### 3. Version Processing Layer (`internal/version/`)

**Responsibility**: Core version logic

**Key components**:

- **Parser** (`parse.go`): Convert string to `Version` struct using format rules
- **Normalizer** (`normalize.go`): Canonicalize versions (e.g., convert strings to integers where applicable)
- **Comparator** (`compare.go`): Compare two `Version` objects using semantic versioning rules
- **Bumper** (`bump.go`): Increment version based on bump kind (major, minor, patch, prerelease, final)
- **Renderer** (`render.go`): Convert `Version` to formatted string for output
- **Policy checker** (`policy.go`): Validate versions against configured policies
- **Conventional commits** (`conventional.go`): Parse commit messages to determine bump requirements
- **Type definitions** (`types.go`): Interfaces and constants for version operations

### 4. Ecosystem Abstraction Layer (`internal/ecosystems/`)

**Responsibility**: Handle ecosystem-specific formatting

- **Registry** (`registry.go`): Central lookup for ecosystems and their renderers
- **Formats** (`formats.go`): Define supported format schemes (semver, v-semver, pep440) and ecosystem aliases (go, python, terraform, etc.)
- **Types** (`types.go`): Define `EcosystemRenderer` interface

### 5. Provider Layer (`internal/providers/`)

**Responsibility**: Fetch versions from external sources

- **Provider interface** (`provider.go`): `VersionProvider` contract that all sources implement
- **Git Tags** (`git_tags.go`): Query versions from local/remote Git tags
- **GitHub Releases** (`github_releases.go`): Query GitHub API for releases
- **GHCR** (`ghcr.go`): Query GitHub Container Registry
- **Cache** (`cache.go`, `cache_test.go`): In-memory caching to avoid redundant queries
- **Retry** (`retry.go`, `retry_test.go`): Exponential backoff retry logic for failed requests

### 6. Sequence Interpretation Layer (`internal/sequence/`)

**Responsibility**: Classify and handle version metadata sequences

- **Interpreter** (`interpreter.go`): Detect and classify sequences (commit SHA, content hash, build ID, numeric)

### 7. Testing Layer (`tests/`)

**Responsibility**: Comprehensive test coverage

- **Golden tests** (`golden_test.go`): Test parsing and rendering against a corpus of known versions
- **Integration tests** (`integration/`): End-to-end testing of commands and provider functionality
- **Test fixtures** (`fixtures/`): Test data and expected outputs

---

## Data Flow Diagrams

### Parse Command Flow

```
User Input (version string)
    ↓
[CLI] Parse Command
    ↓
[Config] Load .verge.yaml
    ↓
[Ecosystem] Select parser for format
    ↓
[Version] Parse string to Version struct
    ↓
[Version] Normalize Version
    ↓
[Version] Render to output format
    ↓
[CLI] Format and display output (text/JSON)
```

### Bump Command Flow

```
User Input (version, bump kind)
    ↓
[CLI] Version Bump Command
    ↓
[Config] Load configuration
    ↓
[Version] Parse current version
    ↓
[Version] Bump according to kind
    ↓
[Version] Apply policy checks
    ↓
[Version] Render bumped version
    ↓
[CLI] Output result
```

### Current/Latest Command Flow

```
[CLI] Current/Latest Command
    ↓
[Config] Load configuration
    ↓
[Providers] Initialize enabled providers (Git Tags, GitHub, GHCR)
    ↓
[Providers] Fetch versions from all sources
    ↓
[Providers] Cache results for future queries
    ↓
[Version] Parse raw results to Version objects
    ↓
[Version] Normalize and sort versions
    ↓
[Version] Select current or latest version
    ↓
[CLI] Format and output result
```

---

## Key Interfaces

### Version Processing

```go
// Parser converts strings to Version objects
type Parser interface {
    Parse(input string) (*Version, error)
}

// Comparator ranks versions
type Comparator interface {
    Compare(a, b *Version) int  // -1, 0, or 1
}

// Bumper increments versions
type Bumper interface {
    Bump(v *Version, kind BumpKind, stage Stage) (*Version, error)
}

// Renderer converts to strings
type Renderer interface {
    Render(v *Version) string
}
```

### Ecosystem Support

```go
// EcosystemRenderer formats versions for specific ecosystems
type EcosystemRenderer interface {
    Name() string
    Render(major, minor, patch int, stage string, 
            sequence interface{}, isPrerelease bool) string
}
```

### Version Sources

```go
// VersionProvider fetches versions from external sources
type VersionProvider interface {
    Name() string
    Fetch(opts QueryOptions) ([]*VersionResult, error)
}
```

---

## Extension Points

### Adding a New Ecosystem

1. Create a renderer in `internal/ecosystems/`
2. Register it in the ecosystem registry
3. Define format schemes for parsing/rendering
4. Add test cases for version parsing/rendering

### Adding a New Version Source

1. Implement the `VersionProvider` interface in `internal/providers/`
2. Add configuration schema to `internal/config/schema.go`
3. Update provider initialization logic
4. Add integration tests

### Custom Policies

1. Extend `internal/version/policy.go` with new validation rules
2. Update configuration schema to expose policy options
3. Apply policies in bump logic

---

## Configuration Structure

The `.verge.yaml` file controls Verge's behavior:

```yaml
version: 1
ecosystem: go                 # Default format for parsing/rendering

format:
  input: semver               # Input format scheme
  output: semver              # Output format scheme
  tagPrefix: "v"              # Prefix for version tags

sources:
  precedence: [git-tags, github-releases]  # Which source to prioritize
  
  git-tags:
    enabled: true
    fetch: true              # Fetch from remote
    includePrerelease: true
    
  github-releases:
    enabled: false
    owner: armckinney
    repo: verge
    
  ghcr:
    enabled: false
    image: ghcr.io/org/image

sequence:
  hashLength: 7              # Abbreviated commit SHA length

autoBump:
  enabled: false             # Auto-bump based on conventional commits
```

---

## Error Handling

Verge uses structured error codes for different failure scenarios:

- **Parse errors**: Invalid version format
- **Config errors**: Missing or invalid configuration
- **Provider errors**: Failed to fetch from source (with retry logic)
- **Policy errors**: Version violates configured policies
- **Comparison errors**: Cannot compare versions from different formats

All errors are mapped to exit codes for CI/CD integration.

---

## Dependencies

**Primary**:
- `github.com/spf13/cobra`: CLI framework
- `gopkg.in/yaml.v3`: YAML configuration parsing

**No external dependencies for core version logic**, allowing lightweight distribution and deployment.

---

## Performance Considerations

- **Caching**: In-memory cache for provider results to avoid redundant API calls
- **Retries**: Exponential backoff with jitter for transient failures
- **Streaming**: Processes versions as they are fetched rather than loading all in memory
- **Concurrency**: Can be extended for parallel provider queries (current implementation is sequential)

---

## Deployment Model

Verge is a self-contained CLI binary with no external runtime dependencies. Distribution options:

- **Binary releases**: Pre-built for Linux, macOS, Windows
- **Go install**: `go install example.com/verge/cmd/verge@latest`
- **Container images**: GHCR for containerized environments
- **CI/CD integration**: Designed for GitHub Actions, GitLab CI, etc.

