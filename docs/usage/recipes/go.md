# Recipe: Go Modules (v-semver) Automation

This recipe guides you through setting up Verge to manage and tag Go module versions, which strictly require a `v` prefixed semantic format.

---

## 1. The Configuration Schema (`.verge.yaml`)

Go modules use Git tags directly to resolve module dependencies. We configure Verge with `version_type: vsemver` and query the `gittag` provider directly to access local module tag histories.

```yaml
version_type: vsemver

default:
  bump_kind: prerelease
  prerelease_stage: dev

sequence:
  type: increment

provider:
  type: gittag
  gittag:
    repo_dir: "."
    include_prerelease: true
```

---

## 2. Go Bumping Variations

Go requires lowercase `v` prefixes on all release tags (e.g., `v1.2.3`). Pre-release pseudo-versions or test tags must preserve this structure.

### Pull Request pre-release (`v<major>.<minor>.<patch>-<stage>.<sequence>`)
Used to publish test modules for PR integration validations:
```bash
$ verge bump --version v1.2.3 --kind prerelease --stage dev
v1.2.4-dev.1
```

### Production release (`v<major>.<minor>.<patch>`)
Promoting the module to a final tag:
```bash
$ verge bump --version v1.2.4-dev.3 --kind final
v1.2.4
```

---

## 3. Go Module Tagging Pipeline Script

The following shell script fetches the historical module tags, calculates the next tag, and publishes it using git commands.

```bash
#!/usr/bin/env bash
set -euo pipefail

# 1. Fetch latest tag from git repo
LATEST=$(verge latest)
echo "Latest module tag: ${LATEST}"

# 2. Bump version depending on branch
if [[ "${GITHUB_REF:-}" == "refs/heads/main" ]]; then
  # Production merge: promote to final release vX.Y.Z
  NEXT=$(verge bump --kind final)
else
  # Feature branch / PR: publish dev tag vX.Y.Z-dev.N
  NEXT=$(verge bump --kind prerelease --stage dev)
fi

echo "Calculated next Go Module version: ${NEXT}"

# 3. Create and push tag to GitHub
git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"

git tag "${NEXT}"
git push origin "${NEXT}"
```
