# verctl Usage Documentation

`verctl` is a semantic versioning CLI tool for parsing, comparing, bumping, and querying versions. It supports three generalized version format schemes — `v-semver`, `semver`, and `pep440` — plus ecosystem-specific aliases for Go, Terraform, Containers, GitHub Actions, and Python.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Global Flags](#global-flags)
- [Commands](#commands)
  - [version parse](#version-parse)
  - [version compare](#version-compare)
  - [version bump](#version-bump)
  - [version current](#version-current)
  - [version latest](#version-latest)
  - [version info](#version-info)
- [Output Formats](#output-formats)
- [Version Format Schemes](#version-format-schemes)
- [Exit Codes](#exit-codes)
- [Configuration](#configuration)
- [Environment Variables](#environment-variables)
- [CI/CD Integration](ci-cd.md)

---

## Installation

### From source

```bash
git clone https://github.com/armckinney/template-go
cd template-go
make build          # produces ./verctl
```

### Using Go install

```bash
go install example.com/template-go/cmd/verctl@latest
```

### Pre-built releases

Download from the [releases page](https://github.com/armckinney/template-go/releases). Binaries are available for:

| Platform       | Binary                         |
|----------------|--------------------------------|
| Linux amd64    | `verctl-linux-amd64`           |
| Linux arm64    | `verctl-linux-arm64`           |
| macOS arm64    | `verctl-darwin-arm64`          |
| Windows amd64  | `verctl-windows-amd64.exe`     |
| Windows arm64  | `verctl-windows-arm64.exe`     |

---

## Quick Start

```bash
# Parse a version
verctl version parse v1.2.3-rc.1

# Compare two versions (exit code 10 = left < right)
verctl version compare 1.2.3 2.0.0; echo $?

# Bump a version
verctl version bump --from 1.2.3 --kind minor

# Get current version from git tags
verctl version current

# Auto-detect bump kind from conventional commits
verctl version bump --from 1.2.3 --auto
```

---

## Global Flags

Available on every command:

| Flag               | Default         | Description                              |
|--------------------|-----------------|------------------------------------------|
| `-c, --config`     | `.verctl.yaml`  | Path to config file                      |
| `-f, --format`     | `text`          | Output format: `text` or `json`          |
| `-v, --verbose`    | `false`         | Enable verbose output                    |

---

## Commands

### version parse

Parse a version string and display its components and ecosystem renderings.

```
verctl version parse <version> [flags]
```

**Flags**

| Flag          | Default | Description                                                             |
|---------------|---------|-------------------------------------------------------------------------|
| `--ecosystem` | `all`   | Render for a specific format scheme (`v-semver`, `semver`, `pep440`, `all`) or ecosystem alias (`go`, `terraform`, `containers`, `github-actions`, `python`) |

**Examples**

```bash
# Parse a plain semver
verctl version parse 1.2.3

# Parse a v-prefixed prerelease
verctl version parse v1.2.3-rc.2

# Parse a PEP 440 version
verctl version parse 1.2.3dev4

# Parse and render for PEP 440 (Python) only
verctl version parse 1.2.3-rc.1 --ecosystem pep440

# JSON output
verctl version parse v1.2.3-rc.2 --format json
```

**Text output** (example for `v1.2.3-rc.2`):

```
input         v1.2.3-rc.2
major         1
minor         2
patch         3
stage         rc
sequence      2
sequenceType  numeric
scheme        semver
prerelease    true
core          1.2.3
rendered.v-semver       v1.2.3-rc.2
rendered.semver         1.2.3-rc.2
rendered.pep440         1.2.3rc2
```

---

### version compare

Compare two version strings.

```
verctl version compare <left> <right> [flags]
```

Exits with `0` (equal), `10` (left < right), or `11` (left > right).

**Examples**

```bash
verctl version compare 1.2.3 2.0.0   # exit 10
verctl version compare 2.0.0 1.2.3   # exit 11
verctl version compare 1.2.3 1.2.3   # exit 0

# In a script
if verctl version compare "$CURRENT" "$CANDIDATE"; then
  echo "equal"
elif [ $? -eq 10 ]; then
  echo "$CANDIDATE is newer"
fi
```

---

### version bump

Compute the next version from a given version and bump kind.

```
verctl version bump [flags]
```

**Flags**

| Flag           | Default  | Description                                                             |
|----------------|----------|-------------------------------------------------------------------------|
| `--from`       | *(required)* | Source version to bump from                                        |
| `--kind`       |          | Bump kind: `major`, `minor`, `patch`, `prerelease`, `final`             |
| `--stage`      |          | Prerelease stage for `prerelease` bumps: `dev`, `alpha`, `beta`, `rc`   |
| `--ecosystem`  | `v-semver` | Render output for this format scheme or ecosystem alias                 |
| `--auto`       | `false`  | Auto-detect bump kind from conventional commits (requires git)          |
| `--repo-dir`   | `.`      | Repository directory (used with `--auto`)                               |
| `--changelog`  | `false`  | Output changelog-friendly JSON instead of default output                |

**Bump kinds**

| Kind          | Effect                                              |
|---------------|-----------------------------------------------------|
| `major`       | `1.2.3` → `2.0.0`                                  |
| `minor`       | `1.2.3` → `1.3.0`                                  |
| `patch`       | `1.2.3` → `1.2.4`                                  |
| `prerelease`  | `1.2.3` → `1.2.4-<stage>.1` (next prerelease)      |
| `final`       | `1.2.3-rc.1` → `1.2.3` (drop prerelease suffix)    |

**Examples**

```bash
# Bump minor
verctl version bump --from 1.2.3 --kind minor
# → 1.3.0

# Bump to a prerelease
verctl version bump --from 1.2.3 --kind prerelease --stage rc
# → 1.2.4-rc.1

# Promote a prerelease to final
verctl version bump --from 1.2.3-rc.1 --kind final
# → 1.2.3

# Auto-detect from conventional commits
verctl version bump --from 1.2.3 --auto
# reads git commits since tag v1.2.3

# Changelog JSON output
verctl version bump --from 1.2.3 --kind minor --changelog --format json
```

**Changelog JSON output** (with `--changelog --format json`):

```json
{
  "version": {
    "from": "1.2.3",
    "to": "1.3.0",
    "bumpType": "minor"
  },
  "metadata": {
    "timestamp": "2026-05-18T03:10:00Z",
    "source": "version-bump",
    "commits": []
  }
}
```

**Auto-bump conventional commits**

With `--auto`, verctl reads `git log <from>..HEAD` and determines the bump kind:

| Commit prefix       | Bump kind |
|---------------------|-----------|
| `BREAKING CHANGE:`  | `major`   |
| `feat!:` / `type!:` | `major`   |
| `feat:`             | `minor`   |
| `fix:`              | `patch`   |

---

### version current

Get the highest (current) version from git tags, excluding prereleases by default.

```
verctl version current [flags]
```

**Flags**

| Flag          | Default | Description                                                  |
|---------------|---------|--------------------------------------------------------------|
| `--repo-dir`  | `.`     | Repository directory                                         |
| `--ecosystem` |         | Render output for this ecosystem (falls back to config)      |
| `--explain`   | `false` | Show all candidates and selection reasoning                  |

**Examples**

```bash
verctl version current
verctl version current --ecosystem python
verctl version current --explain
verctl version current --format json
```

---

### version latest

Get the latest (highest) version from git tags, with optional stage/core filtering.

```
verctl version latest [flags]
```

**Flags**

| Flag          | Default | Description                                                    |
|---------------|---------|----------------------------------------------------------------|
| `--repo-dir`  | `.`     | Repository directory                                           |
| `--stage`     |         | Filter by stage: `dev`, `alpha`, `beta`, `rc`, `final`         |
| `--core`      |         | Filter by core version (e.g. `1.2.3`)                          |
| `--ecosystem` |         | Render output for this ecosystem                               |
| `--explain`   | `false` | Show all candidates, filters applied, and selection reasoning  |

**Examples**

```bash
# Latest overall
verctl version latest

# Latest rc candidate
verctl version latest --stage rc

# Latest dev build for a specific core version
verctl version latest --core 1.2.3 --stage dev

# JSON output
verctl version latest --format json

# Show decision reasoning
verctl version latest --explain
```

---

### version info

Show the verctl build version information.

```bash
verctl version info
```

---

## Output Formats

### Text (default)

Key-value pairs, one per line:

```
major       1
minor       2
patch       3
stage       rc
rendered    v1.2.3-rc.1
```

### JSON (`--format json`)

Machine-readable JSON object:

```json
{
  "major": 1,
  "minor": 2,
  "patch": 3,
  "stage": "rc",
  "rendered": "v1.2.3-rc.1"
}
```

Use JSON output when integrating with CI scripts, `jq`, or changelog tools:

```bash
verctl version bump --from 1.2.3 --kind minor --format json | jq -r '.rendered'
# v1.3.0
```

---

## Version Format Schemes

`verctl` renders versions using three generalized format schemes. Ecosystem names are accepted as aliases for convenience.

### Canonical format schemes

| Scheme      | Final format | Prerelease format    | Description                       |
|-------------|--------------|----------------------|-----------------------------------|
| `v-semver`  | `v1.2.3`     | `v1.2.3-rc.1`        | SemVer with `v` prefix            |
| `semver`    | `1.2.3`      | `1.2.3-rc.1`         | Standard SemVer (no prefix)       |
| `pep440`    | `1.2.3`      | `1.2.3rc1`           | Python PEP 440 prerelease format  |

### Ecosystem-to-scheme mapping

Use canonical scheme names in scripts and config for portability. Ecosystem aliases are provided for familiarity and backward compatibility.

| Ecosystem alias  | Maps to   | Example final | Example prerelease  |
|------------------|-----------|---------------|---------------------|
| `go`             | `v-semver` | `v1.2.3`     | `v1.2.3-rc.1`       |
| `terraform`      | `v-semver` | `v1.2.3`     | `v1.2.3-rc.1`       |
| `containers`     | `semver`   | `1.2.3`      | `1.2.3-rc.1`        |
| `github-actions` | `semver`   | `1.2.3`      | `1.2.3-rc.1`        |
| `python`         | `pep440`   | `1.2.3`      | `1.2.3rc1`          |

---

## Exit Codes

| Code | Meaning                          |
|------|----------------------------------|
| `0`  | Success / versions are equal     |
| `1`  | General error                    |
| `2`  | Usage / invalid arguments        |
| `10` | Left version < right version     |
| `11` | Left version > right version     |
| `20` | Version not found                |
| `21` | Invalid version format           |
| `22` | Configuration error              |
| `30` | Network error (remote providers) |

---

## Configuration

`verctl` looks for `.verctl.yaml` in the current directory (or the path given by `--config`). All settings are optional and fall back to defaults.

```yaml
version: 1

# Default format scheme for rendering output
ecosystem: v-semver

format:
  input: auto          # Version scheme detection: auto | semver | pep440
  output: auto         # Output rendering: auto | semver | pep440
  tagPrefix: v         # Prefix stripped from git tags when parsing
  sequenceInterpreter: auto  # Sequence type: auto | numeric | hash

sources:
  precedence:
    - git-tags         # Source search order (git-tags, github-releases, ghcr)
  git-tags:
    enabled: true
    fetch: false                 # Run `git fetch` before listing tags
    includePrerelease: true
    ecosystemParsing: v-semver
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

sequence:
  hashLength: 7
  allowContentHash: true
  ghBuildPattern: "gh-"

rules:
  prereleaseStage: dev           # Default stage for new prereleases
  allowMajorZeroBreaking: true   # Allow breaking changes on v0.x
  defaultBump: patch

autoBump:
  conventionalCommits: true
  breakingTokens:
    - "BREAKING CHANGE"
    - "!:"
```

---

## Environment Variables

| Variable               | Overrides config key           | Description                                    |
|------------------------|--------------------------------|------------------------------------------------|
| `VERCTL_ECOSYSTEM`     | `ecosystem`                    | Default format scheme for rendering (`v-semver`, `semver`, `pep440`) |
| `VERCTL_FORMAT_OUTPUT` | `format.output`                | Output format (`text` or `json`)               |
| `VERCTL_TAG_PREFIX`    | `format.tagPrefix`             | Git tag prefix                                 |
| `GITHUB_TOKEN`         | *(provider auth)*              | Token for GitHub Releases and GHCR providers   |

---

## CI/CD Integration

### GitHub Actions — auto-bump and tag on merge

```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # required for git log history

      - name: Get current version
        run: echo "CURRENT=$(verctl version current --format json | jq -r .version)" >> $GITHUB_ENV

      - name: Auto-detect next version from conventional commits
        run: echo "NEXT=$(verctl version bump --from "$CURRENT" --auto --format json | jq -r .to)" >> $GITHUB_ENV

      - name: Tag and push
        run: |
          git tag "v${{ env.NEXT }}"
          git push --tags
```

### Shell script — guard against version downgrades

```bash
#!/usr/bin/env bash
set -euo pipefail

CURRENT=$(verctl version current --format json | jq -r .version)
PROPOSED="${1:?usage: $0 <proposed-version>}"

verctl version compare "$CURRENT" "$PROPOSED"
CODE=$?

if [ "$CODE" -eq 11 ]; then
  echo "Error: $PROPOSED is less than current $CURRENT" >&2
  exit 1
fi
echo "OK: $PROPOSED is valid next version (current: $CURRENT)"
```

### Parse a version tag and export components

```bash
TAG="v1.3.0-rc.2"
eval "$(verctl version parse "$TAG" --format json | jq -r '
  "MAJOR=\(.parsed.major)",
  "MINOR=\(.parsed.minor)",
  "PATCH=\(.parsed.patch)",
  "STAGE=\(.parsed.stage)"
')"
echo "Building $MAJOR.$MINOR.$PATCH ($STAGE)"
```

### Multi-scheme release matrix

```bash
VERSION="1.3.0-rc.1"
for scheme in v-semver semver pep440; do
  rendered=$(verctl version parse "$VERSION" --ecosystem "$scheme" --format json | jq -r .rendered)
  echo "$scheme: $rendered"
done
```
