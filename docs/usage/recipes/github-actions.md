# Recipe: GitHub Actions & Reusable Workflows Automation

This recipe guides you through setting up Verge to calculate and tag custom **GitHub Actions** and **Reusable Workflows** (RWF). 

Actions are referenced by consumers using exact or floating semantic versions (e.g. `actions/checkout@v4`). Having dynamic floating tags is critical to ensure consumers automatically receive non-breaking patch updates.

---

## 1. The Configuration Schema (`.verge.yaml`)

We use standard `semver` (without prefix, or with depending on your preference, though standard is `semver` and users reference them using `v` prefix like `@v4` on the tag itself).

```yaml
version_type: semver

default:
  bump_kind: patch
  prerelease_stage: rc

sequence:
  type: increment

provider:
  type: ghrelease
  ghrelease:
    owner: "my-org"
    repo: "custom-action"
    include_prerelease: false
```

---

## 2. Complete GitHub Actions Tagging Workflow

The following pipeline fetches the latest action tag, calculates the next patch, tags the exact release, and updates the **floating major and minor tags** (e.g. `v1`, `v1.2`).

```yaml
name: Release Action

on:
  push:
    branches:
      - main

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Build Verge CLI
        run: go build -o /usr/local/bin/verge cmd/verge/main.go

      - name: Calculate Version Bump
        id: verge
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          NEXT_VERSION=$(verge bump --kind patch)
          echo "version=${NEXT_VERSION}" >> $GITHUB_OUTPUT
          
          # Extract major version (e.g. from 1.2.3 -> v1)
          MAJOR_VERSION="v${NEXT_VERSION%%.*}"
          echo "major=${MAJOR_VERSION}" >> $GITHUB_OUTPUT
          
          # Extract minor version (e.g. from 1.2.3 -> v1.2)
          MINOR_VERSION="v${NEXT_VERSION%.*}"
          echo "minor=${MINOR_VERSION}" >> $GITHUB_OUTPUT

      - name: Create Git Release Tags
        run: |
          git config --global user.name "github-actions[bot]"
          git config --global user.email "github-actions[bot]@users.noreply.github.com"
          
          # 1. Tag exact version
          git tag "v${{ steps.verge.outputs.version }}"
          git push origin "v${{ steps.verge.outputs.version }}"
          
          # 2. Update floating major tag (e.g. v1)
          git tag -f "${{ steps.verge.outputs.major }}"
          git push origin -f "${{ steps.verge.outputs.major }}"
          
          # 3. Update floating minor tag (e.g. v1.2)
          git tag -f "${{ steps.verge.outputs.minor }}"
          git push origin -f "${{ steps.verge.outputs.minor }}"
```
