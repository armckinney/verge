# Recipe: Container Image Tags Automation

This recipe guides you through setting up Verge to calculate dynamic tags for Docker containers, covering PR builds, final releases, latest pointers, and floating version tag updates.

---

## 1. The Configuration Schema (`.verge.yaml`)

We use standard `semver` parsing alongside the `filehash` sequence calculator. By hashing the `Dockerfile` and the package configs, Verge ensures that container images are only bumped and built when there are actual source file mutations.

```yaml
version_type: semver

default:
  bump_kind: prerelease
  prerelease_stage: dev

sequence:
  type: filehash
  targets:
    - "./Dockerfile"
    - "./go.mod"
  length: 7

provider:
  type: ghcr
  ghcr:
    package: "owner/app"
    include_prerelease: true
```

---

## 2. Pull Request Bumping (PR Tagging)
* **Goal:** Calculate a unique pre-release container tag: `major.minor.patch-stage.<hash>`.
* **Behavior:** When building on a PR branch, we run the default pre-release bump. Since `sequence.type` is set to `filehash`, the suffix is dynamically resolved to a 7-character file content hash.

### Execution:
```bash
$ verge bump
1.2.4-dev.bf88455
```

---

## 3. Production Release Bumping (REL Tagging)
* **Goal:** Promote the pre-release version to a final production tag: `major.minor.patch`.
* **Behavior:** When merging to the master production branch, we explicitly pass `--kind final` to wipe the pre-release stage and sequence markers.

### Execution:
```bash
$ verge bump --kind final
1.2.4
```

---

## 4. Complete Automation Pipeline Script

The following shell script builds and tags container images for both PR branches and production releases, pushing final, latest, and floating semver tags.

```bash
#!/usr/bin/env bash
set -euo pipefail

# Define registry URL
REGISTRY_IMAGE="ghcr.io/owner/app"

# Detect if current build is a production release merge vs a PR build
IS_RELEASE=false
if [[ "${GITHUB_REF:-}" == "refs/heads/main" ]]; then
  IS_RELEASE=true
fi

if [ "$IS_RELEASE" = true ]; then
  # 1. Promote to final release: major.minor.patch
  TAG=$(verge bump --kind final)
  log "Resolved production tag: ${TAG}"
  
  # Build and push exact release version
  docker build -t "${REGISTRY_IMAGE}:${TAG}" .
  docker push "${REGISTRY_IMAGE}:${TAG}"
  
  # 2. Tag and push "latest" pointer
  docker tag "${REGISTRY_IMAGE}:${TAG}" "${REGISTRY_IMAGE}:latest"
  docker push "${REGISTRY_IMAGE}:latest"
  
  # 3. Generate Floating Tags (Major and Minor)
  # Extract components (e.g. from 1.2.4 -> MAJOR=1, MINOR=1.2)
  MAJOR="${TAG%%.*}"
  MINOR_PART="${TAG%.*}"
  
  # Tag and push floating major
  docker tag "${REGISTRY_IMAGE}:${TAG}" "${REGISTRY_IMAGE}:${MAJOR}"
  docker push "${REGISTRY_IMAGE}:${MAJOR}"
  
  # Tag and push floating minor
  docker tag "${REGISTRY_IMAGE}:${TAG}" "${REGISTRY_IMAGE}:${MINOR_PART}"
  docker push "${REGISTRY_IMAGE}:${MINOR_PART}"
  
else
  # 4. PR Build: major.minor.patch-dev.<hash>
  TAG=$(verge bump)
  log "Resolved PR tag: ${TAG}"
  
  docker build -t "${REGISTRY_IMAGE}:${TAG}" .
  docker push "${REGISTRY_IMAGE}:${TAG}"
fi
```
