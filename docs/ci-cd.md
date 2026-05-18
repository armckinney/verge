# CI/CD Integration Guide

`verctl` is designed to work seamlessly in CI pipelines. This guide covers common patterns for GitHub Actions, shell scripts, and release automation.

---

## Table of Contents

- [GitHub Actions](#github-actions)
- [Shell Script Patterns](#shell-script-patterns)
- [Auto-Bump with Conventional Commits](#auto-bump-with-conventional-commits)
- [Remote Providers](#remote-providers)
- [Changelog Generation](#changelog-generation)

---

## GitHub Actions

### Basic version bump on push to main

```yaml
name: Release

on:
  push:
    branches: [main]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0        # required for git tag history

      - name: Install verctl
        run: |
          curl -sSL https://github.com/armckinney/template-go/releases/latest/download/verctl-linux-amd64 \
            -o /usr/local/bin/verctl
          chmod +x /usr/local/bin/verctl

      - name: Get current version
        id: current
        run: echo "version=$(verctl version current --format json | jq -r '.rendered')" >> $GITHUB_OUTPUT

      - name: Auto-detect bump kind
        id: bump
        run: |
          NEW=$(verctl version bump --from ${{ steps.current.outputs.version }} --auto --format json | jq -r '.to')
          echo "new_version=$NEW" >> $GITHUB_OUTPUT

      - name: Tag and push
        run: |
          git tag "v${{ steps.bump.outputs.new_version }}"
          git push origin "v${{ steps.bump.outputs.new_version }}"
```

### Parse version and expose components

```yaml
- name: Parse version
  id: ver
  run: |
    JSON=$(verctl version parse ${{ github.ref_name }} --format json)
    echo "major=$(echo $JSON | jq -r '.parsed.major')" >> $GITHUB_OUTPUT
    echo "minor=$(echo $JSON | jq -r '.parsed.minor')" >> $GITHUB_OUTPUT
    echo "patch=$(echo $JSON | jq -r '.parsed.patch')" >> $GITHUB_OUTPUT
    echo "stage=$(echo $JSON | jq -r '.parsed.stage')" >> $GITHUB_OUTPUT
    echo "prerelease=$(echo $JSON | jq -r '.parsed.stage != "final"')" >> $GITHUB_OUTPUT

- name: Build with version
  run: |
    docker build \
      --label "version=${{ github.ref_name }}" \
      --label "major=${{ steps.ver.outputs.major }}" \
      -t myimage:${{ github.ref_name }} .
```

### Conditional release based on version stage

```yaml
- name: Parse version
  id: ver
  run: |
    STAGE=$(verctl version parse ${{ github.ref_name }} --format json | jq -r '.parsed.stage')
    echo "stage=$STAGE" >> $GITHUB_OUTPUT

- name: Publish to production
  if: steps.ver.outputs.stage == 'final'
  run: echo "Publishing final release..."

- name: Publish to staging
  if: steps.ver.outputs.stage == 'rc'
  run: echo "Publishing RC to staging..."
```

### Compare versions (gate on regression)

```yaml
- name: Check version is newer
  run: |
    CURRENT=$(verctl version current --format json | jq -r '.normalized')
    CANDIDATE="${{ github.ref_name }}"
    verctl version compare "$CURRENT" "$CANDIDATE"
    STATUS=$?
    if [ $STATUS -eq 11 ]; then
      echo "Error: $CANDIDATE is older than current $CURRENT"
      exit 1
    fi
    echo "Version $CANDIDATE is valid (>= $CURRENT)"
```

---

## Shell Script Patterns

### Get current version and render for ecosystem

```bash
#!/usr/bin/env bash
set -euo pipefail

# Get the current stable version rendered for Python
PYTHON_VERSION=$(verctl version current --ecosystem python --format json | jq -r '.rendered')
echo "Current Python version: $PYTHON_VERSION"

# Write to version file
echo "__version__ = \"$PYTHON_VERSION\"" > src/_version.py
```

### Bump and create git tag

```bash
#!/usr/bin/env bash
set -euo pipefail

CURRENT=$(verctl version current --format json | jq -r '.normalized')
NEW=$(verctl version bump --from "$CURRENT" --kind "$1" --ecosystem go --format json | jq -r '.rendered')

echo "Bumping $CURRENT → $NEW"
git tag "$NEW"
git push origin "$NEW"
echo "Tagged $NEW"
```

Usage: `./scripts/release.sh minor`

### Check if a version is a prerelease

```bash
#!/usr/bin/env bash
VERSION=$1
STAGE=$(verctl version parse "$VERSION" --format json | jq -r '.parsed.stage')

if [ "$STAGE" = "final" ]; then
  echo "$VERSION is a stable release"
else
  echo "$VERSION is a prerelease (stage: $STAGE)"
  exit 1
fi
```

### Sort a list of version strings

```bash
#!/usr/bin/env bash
# Print versions in ascending order using verctl compare as sort key
versions=("v2.0.0" "v1.2.3" "v1.3.0-rc.1" "v1.3.0")

printf '%s\n' "${versions[@]}" | sort -t. -k1,1V -k2,2n -k3,3n
# For correct prerelease ordering, compare pairs:
for v in "${versions[@]}"; do
  echo "$(verctl version parse "$v" --format json | jq -r '"\(.parsed.major).\(.parsed.minor).\(.parsed.patch) \(.parsed.stage) \(.parsed.sequence // 0)"') $v"
done | sort -k1,1n -k2,2n -k3,3n | awk '{print $NF}'
```

---

## Auto-Bump with Conventional Commits

`verctl version bump --auto` reads commit messages since the `--from` tag and determines the bump kind automatically.

### Commit message format

```
<type>[optional scope][!]: <description>

[optional body]

[optional footer: BREAKING CHANGE: <description>]
```

| Pattern                         | Detected bump |
|---------------------------------|---------------|
| `feat: add new endpoint`        | `minor`       |
| `fix: correct nil pointer`      | `patch`       |
| `feat!: remove deprecated API`  | `major`       |
| `fix(auth)!: revoke old tokens` | `major`       |
| `BREAKING CHANGE: ...` footer   | `major`       |
| `chore: update deps`            | *(no bump)*   |
| `docs: update README`           | *(no bump)*   |

### Usage in CI

```bash
# Get current version from latest tag
CURRENT=$(verctl version current --format json | jq -r '.normalized')

# Auto-detect bump and compute new version
NEW=$(verctl version bump --from "$CURRENT" --auto --format json | jq -r '.to')

echo "Bumping $CURRENT → $NEW"
```

### Custom breaking tokens

Configure in `.verctl.yaml`:

```yaml
autoBump:
  conventionalCommits: true
  breakingTokens:
    - "BREAKING CHANGE"
    - "BREAKING-CHANGE"
    - "!:"
    - "MAJOR:"       # custom token
```

---

## Remote Providers

### GitHub Releases provider

Configure in `.verctl.yaml`:

```yaml
sources:
  precedence:
    - github-releases
  github-releases:
    enabled: true
    owner: my-org
    repo: my-app
    includePrerelease: true
```

Set `GITHUB_TOKEN` for private repos or to avoid rate limits:

```bash
export GITHUB_TOKEN=ghp_xxxx
verctl version current
```

### GHCR (Container Registry) provider

```yaml
sources:
  precedence:
    - ghcr
  ghcr:
    enabled: true
    image: ghcr.io/my-org/my-image
    includePrerelease: true
```

Filter to release-only tags:

```yaml
ghcr:
  enabled: true
  image: ghcr.io/my-org/my-image
  channelFilter: rel   # rel | pr | floating
```

---

## Changelog Generation

### Generate changelog JSON on bump

```bash
verctl version bump \
  --from 1.2.3 \
  --kind minor \
  --changelog \
  --format json > changelog-entry.json
```

Output schema:

```json
{
  "version": {
    "from": "1.2.3",
    "to": "1.3.0",
    "bumpType": "minor"
  },
  "metadata": {
    "timestamp": "2026-05-18T00:00:00Z",
    "source": "version-bump",
    "commits": []
  }
}
```

### Integrating with conventional-changelog

```bash
# Generate changelog data
verctl version bump --from "$CURRENT" --auto --changelog --format json > bump.json

# Extract fields
FROM=$(jq -r '.version.from' bump.json)
TO=$(jq -r '.version.to' bump.json)
BUMP=$(jq -r '.version.bumpType' bump.json)

echo "## [$TO] - $(date +%Y-%m-%d)" >> CHANGELOG.md
echo "" >> CHANGELOG.md
echo "**Bump type:** $BUMP (from $FROM)" >> CHANGELOG.md
```

### GitHub Actions release notes

```yaml
- name: Generate release notes
  id: notes
  run: |
    JSON=$(verctl version bump --from "${{ steps.current.outputs.version }}" \
      --auto --changelog --format json)
    FROM=$(echo $JSON | jq -r '.version.from')
    TO=$(echo $JSON | jq -r '.version.to')
    TYPE=$(echo $JSON | jq -r '.version.bumpType')
    echo "body=Bumped $FROM → $TO ($TYPE)" >> $GITHUB_OUTPUT

- name: Create release
  uses: actions/create-release@v1
  with:
    tag_name: "v${{ steps.bump.outputs.new_version }}"
    release_name: "v${{ steps.bump.outputs.new_version }}"
    body: ${{ steps.notes.outputs.body }}
```
