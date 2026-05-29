# Verge

Verge is a deterministic, fast version generation CLI explicitly designed for pipelines and automation scenarios.

It acts as a single-source-of-truth semantic versioning machine that reliably fetches, sequences, and bumps versions using standard ecosystems (semver, vsemver, pep440) across disparate providers (git tags, ghcr, github releases).

## Core Principles
1. **Deterministic Operation:** The CLI acts as a pure calculator. Give it a state, and it returns predictably formatted output strings.
2. **Immutable Environment:** Verge operates strictly read-only. It calculates tags but never intrinsically applies them.
3. **Pipeline First:** Native `--format json` flags and strictly defined exit codes ensure 100% interoperability with jq and bash evaluations natively.

## Installation

Install the Verge CLI automatically across Linux (amd64/arm64), macOS (amd64/arm64), and Windows (Git Bash/WSL) with our cross-platform installation script:

```bash
curl -sSL https://raw.githubusercontent.com/armckinney/verge/main/install.sh | bash
```

Alternatively, you can manually download the compiled archive for your specific environment from the [GitHub Releases Page](https://github.com/armckinney/verge/releases).

## Quick Examples

```bash
# 1. Scaffold a Git-Tag configuration boilerplate
verge ini-config --template gittag-semver

# 2. Query the highest tag from your configured provider
verge latest
# Output: 1.2.3

# 3. Calculate next pre-release version with auto-increment
verge bump --kind prerelease --stage dev
# Output: 1.2.4-dev.1
```

## Documentation
- **Configuration**
  - [Overview](docs/usage/configuration/index.md)
  - [Providers](docs/usage/configuration/providers.md)
  - [Version Types](docs/usage/configuration/version_types.md)
  - [Sequences](docs/usage/configuration/sequence.md)
- **Commands**
  - [CLI Index](docs/usage/commands/index.md)
  - [verge current](docs/usage/commands/current.md)
  - [verge latest](docs/usage/commands/latest.md)
  - [verge bump](docs/usage/commands/bump.md)
  - [verge ini-config](docs/usage/commands/ini-config.md)
- **Recipes**
  - [Container Images](docs/usage/recipes/container-images.md)
  - [Devcontainer Features](docs/usage/recipes/devcontainer-features.md)
  - [Python Packages](docs/usage/recipes/python.md)
  - [Go Modules](docs/usage/recipes/go.md)
  - [Terraform Modules & Providers](docs/usage/recipes/terraform.md)
  - [GitHub Actions / CI](docs/usage/recipes/github-actions.md)
- **Development**
  - [Architecture](docs/development/architecture.md)
  - [Product Specification](docs/development/spec.md)
