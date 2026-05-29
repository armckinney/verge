# Verge

Verge is a deterministic, fast version generation CLI explicitly designed for pipelines and automation scenarios.

It acts as a single-source-of-truth semantic versioning machine that reliably fetches, sequences, and bumps versions using standard ecosystems (semver, vsemver, pep440) across disparate providers (git tags, ghcr, github releases).

## Getting Started

Follow these steps to quickly install, configure, and execute version bumping in your workspace.

### 1. Installation

Install the Verge CLI automatically across Linux, macOS, and Windows (Git Bash/WSL) with our cross-platform one-liner:

```bash
curl -sSL https://raw.githubusercontent.com/armckinney/verge/main/install.sh | bash
```

*Alternatively, manually download compiled archives for your platform from the [GitHub Releases Page](https://github.com/armckinney/verge/releases).*

---

### 2. Configuration

Verge requires a `.verge.yaml` configuration file at the root of your repository to establish the single source of truth for version formats, sequence generators, and providers.

Generate a standard Git-Tag SemVer configuration boilerplate using the `init` helper:

```bash
verge init --template gittag-semver
```

This creates a fresh `.verge.yaml` file:
```yaml
version_type: semver
provider:
  type: gittag
  gittag:
    repo_dir: "."
    include_prerelease: true
```

---

### 3. Basic Setup & Prerequisites (Getting Ready)

Before running version commands, ensure your workspace meets the following minimal prerequisites:

* **Base Git Tag:** Verge calculates bumps *incrementally* from history. If you are using the default `gittag` provider, your Git repository must have at least one valid tag pushed. Create your first base tag to initialize history:
  ```bash
  git tag v0.1.0
  ```
* **Registry Authentication Token:** If you are using remote providers like `ghrelease` or `ghcr` to fetch tags from private repositories, Verge requires credentials. Export your GitHub Personal Access Token inside your terminal:
  ```bash
  export GITHUB_TOKEN="ghp_your_token_here"
  ```

---

### 4. Basic Usage Examples

Once installed and configured, execute Verge to inspect and calculate versions cleanly:

* **Query Latest Tag:** Retrieves the highest parsed tag in your history matching your format rules:
  ```bash
  verge latest
  # Output: 0.1.0
  ```
* **Calculate Next Version:** Computes the next logical tag based on your `.verge.yaml` default pre-release rules:
  ```bash
  verge bump --kind prerelease --stage dev
  # Output: 0.1.1-dev.1
  ```

---

### 5. Centralized CI/CD Integration

To run versioning reliably inside automated pipelines, Verge is designed to integrate cleanly with centralized CI/CD workflows. 

If you are using our organization's central reusable workflows, you can utilize the generic **`rwf-tag-semver.yaml`** workflow hosted inside our [Central CI/CD Repository](https://github.com/armckinney/cicd) to automate version checks, Git tagging, and error cleanups.

For a complete step-by-step setup walkthrough of calling this reusable workflow, refer to our [Central CI/CD Reusable Workflows Guide](docs/usage/cicd-integration.md).

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
  - [verge init](docs/usage/commands/init.md)
  - [Uninstalling Verge](docs/usage/uninstall.md)
- **Recipes**
  - [Container Images](docs/usage/recipes/container-images.md)
  - [Devcontainer Features](docs/usage/recipes/devcontainer-features.md)
  - [Python Packages](docs/usage/recipes/python.md)
  - [Go Modules](docs/usage/recipes/go.md)
  - [Terraform Modules & Providers](docs/usage/recipes/terraform.md)
  - [GitHub Actions / CI](docs/usage/recipes/github-actions.md)
  - [Central CI/CD Reusable Workflows Guide](docs/usage/cicd-integration.md)
- **Development**
  - [Architecture](docs/development/architecture.md)
  - [Product Specification](docs/development/spec.md)
