# Recipe: Devcontainer Features Automation

This recipe guides you through setting up Verge to calculate version tags and floating releases for custom **Devcontainer Features**.

---

## 1. The Configuration Schema (`.verge.yaml`)

Devcontainer features require strict standard `semver` formatting (without a `v` prefix). We use the `ghrelease` provider as feature updates are usually packaged and attached to GitHub Release archives.

```yaml
version_type: semver

default:
  bump_kind: patch
  prerelease_stage: dev

sequence:
  type: increment

provider:
  type: ghrelease
  ghrelease:
    owner: "my-org"
    repo: "devcontainer-features"
    include_prerelease: true
```

---

## 2. Dynamic Feature Release Script

Devcontainer Features are distributed as `.tgz` archives attached to GitHub Releases. The CLI feature manager matches published tags using strict `semver` rules.

The following script automates parsing feature history, bumping the patch, packaging the feature, and tagging the release.

```bash
#!/usr/bin/env bash
set -euo pipefail

# 1. Fetch latest feature version from GH Release
LATEST_VERSION=$(verge latest)
echo "Latest published devcontainer version: ${LATEST_VERSION}"

# 2. Bump patch for the new release
NEXT_VERSION=$(verge bump --kind patch)
echo "Next release version: ${NEXT_VERSION}"

# 3. Compile and Package the Devcontainer Feature
# (Using the devcontainers CLI tool)
npx -y @devcontainers/cli features package ./src --output-folder ./dist --version "${NEXT_VERSION}"

# 4. Git Tag and Push
git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"

git tag "${NEXT_VERSION}"
git push origin "${NEXT_VERSION}"

# 5. Output Floating tags for user references
# E.g., if version is 1.2.3:
# Major floating tag: 1
# Latest tag: latest
MAJOR_TAG="${NEXT_VERSION%%.*}"

echo "Creating floating tags..."
git tag -f "${MAJOR_TAG}"
git push origin -f "${MAJOR_TAG}"

git tag -f "latest"
git push origin -f "latest"
```
