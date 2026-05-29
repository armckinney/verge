# Recipe: Python PEP 440 Packages Automation

This recipe guides you through setting up Verge to manage and build Python package versions adhering strictly to the **PEP 440** standard.

---

## 1. The Configuration Schema (`.verge.yaml`)

Python packages strictly mandate PEP 440 version grammar. By setting `version_type: pep440`, Verge automatically renders pre-release punctuation-less sequence suffixes (e.g., `1.2.3rc1` instead of `1.2.3-rc.1`).

```yaml
version_type: pep440

default:
  bump_kind: prerelease
  prerelease_stage: dev      # Mapped to PEP 440 "dev" suffix

sequence:
  type: increment

provider:
  type: ghrelease
  ghrelease:
    owner: "my-org"
    repo: "my-python-lib"
    include_prerelease: true
```

---

## 2. Python Bumping Variations

PEP 440 pre-releases are formatted without separation dots or hyphens:

### Pull Request pre-release (`<major>.<minor>.<patch><stage><sequence>`)
Useful to build test packages for test-PyPI on pull request commits:
```bash
$ verge bump --version 1.2.3 --kind prerelease --stage dev
1.2.4dev1

# Transitioning from dev to release candidate stage
$ verge bump --version 1.2.4dev4 --kind prerelease --stage rc
1.2.4rc1
```

### Production release (`<major>.<minor>.<patch>`)
Promoting features to final PyPI releases:
```bash
$ verge bump --version 1.2.4rc2 --kind final
1.2.4
```

---

## 3. Python Packaging Build & Publish Pipeline

The following workflow script extracts the calculated PEP 440 version from Verge, injects it dynamically into package configs, builds source distributions, and publishes to PyPI using `hatch` or `twine`.

```bash
#!/usr/bin/env bash
set -euo pipefail

# 1. Fetch latest version and calculate release bump
TAG=$(verge bump --kind prerelease)
echo "Next PEP 440 version: ${TAG}"

# 2. Inject version into hatchling or pyproject.toml
# If using Hatch, we can override version during build:
# Or rewrite version file:
echo "__version__ = \"${TAG}\"" > src/my_package/__about__.py

# 3. Build distributions (source and wheel)
python -m pip install --upgrade build twine
python -m build

# 4. Upload to Test PyPI or Production PyPI
# Export PyPI credentials securely in environment beforehand
export TWINE_USERNAME="__token__"
export TWINE_PASSWORD="${PYPI_API_TOKEN}"

python -m twine upload --repository pypi dist/*
```
