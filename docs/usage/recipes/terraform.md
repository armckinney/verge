# Recipe: Terraform Modules & Providers Automation

This recipe guides you through setting up Verge to calculate and tag Terraform Providers and Modules. Like Go, the Terraform registry strictly requires and parses lowercase `v` prefixed semver tags (`vsemver`).

---

## 1. The Configuration Schema (`.verge.yaml`)

Terraform module registries ingest module releases directly via GitHub Releases tags. We parse tags under `vsemver` and query the `ghrelease` provider.

```yaml
version_type: vsemver

default:
  bump_kind: patch
  prerelease_stage: rc

sequence:
  type: increment

provider:
  type: ghrelease
  ghrelease:
    owner: "my-org"
    repo: "terraform-aws-vpc"
    include_prerelease: true
```

---

## 2. Terraform Module Tagging Pipeline

The following workflow fetches the latest registry version tag, calculates the next patch, and registers the release.

```bash
#!/usr/bin/env bash
set -euo pipefail

# 1. Fetch latest release version from Registry (GitHub Releases)
LATEST=$(verge latest)
echo "Current Terraform Module tag: ${LATEST}"

# 2. Bump module patch
NEXT=$(verge bump --kind patch)
echo "Next Terraform Module tag: ${NEXT}"

# 3. Create Git Tag and push to trigger Registry indexing
git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"

git tag "${NEXT}"
git push origin "${NEXT}"
```

---

## 3. Terraform Provider Publishing with Goreleaser

Terraform Providers require signed binary releases. By using Verge in your release pipeline, you can calculate the tag, sign the binaries, and trigger `.goreleaser.yaml` builds automatically:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Calculate the next final provider version
PROVIDER_VERSION=$(verge bump --kind final)
echo "Publishing Provider: ${PROVIDER_VERSION}"

# Run goreleaser release passing the version tag
goreleaser release --clean
```
