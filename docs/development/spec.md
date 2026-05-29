# Verge CLI Product Specification

This document details the definitive product specification for **Verge**, a lightweight, fast, and deterministic semantic version calculation tool built specifically for automated build environments and CI/CD pipelines.

---

## 1. Product Goals & Core Philosophy

Verge addresses a common problem in continuous integration pipelines: parsing, calculating, and incrementing version strings consistently across heterogeneous ecosystems (Docker, Python, Go, Terraform) and version tracking registries (Git tags, GitHub Releases, GitHub Container Registry).

### Core Tenants:
1. **Side-Effect Free (Immutable Environment):** Verge strictly calculates and formats version strings. It **never** creates git tags, pushes container images, or commits file modifications. Downstream workflow tasks consume its stdout to execute these mutations.
2. **Bash-First Interoperability:** To support clean bash pipe consumption:
   * **Stdout Hygiene:** Only the exact calculated or requested version string is written to `stdout`. All trace logs, diagnostics, and errors are outputted to `stderr`.
   * **Structured Outputs:** Supports raw string representations (`text`) or enriched metadata payloads (`json`).
3. **Precedence Override Priority:** Settings are cascaded and resolved in a strict order:
   $$\text{CLI Default Code Constants} < \text{Configuration File (.verge.yaml)} < \text{Environment Variables} < \text{CLI Flags}$$

---

## 2. Command-Line Interface (CLI) Spec

### Global Options

* `-c, --config string`: Explicit path to the configuration file (default: `.verge.yaml`). Also respects the `VERGE_CONFIG` environment variable.
* `-f, --format string`: Serialization output format: `text` or `json` (default: `"text"`).
* `-v, --verbose`: Enable diagnostic verbose logging on `stderr`.
* `--field string`: Pluck a single top-level field from the structured command output (useful to extract parameters in shell scripts).

---

### Command: `verge current`
Retrieves the currently active version from a local repository.

* **Support Limits:** Only supports local tracking providers (e.g. `gittag`). If configured with a remote provider (`ghrelease`, `ghcr`), returns an error with exit code `1`.
* **Output Format:**
  * `text` (Default): Returns raw tag string (e.g. `v1.2.3`).
  * `json`: Emits:
    ```json
    {
      "version": "v1.2.3",
      "normalized": "1.2.3",
      "rendered": "v1.2.3"
    }
    ```

---

### Command: `verge latest`
Queries the active provider for the highest recorded version matching format rules.

* **Flags:**
  * `-t, --type string`: Override version format (`semver` | `vsemver` | `pep440`).
  * `-p, --provider string`: Override active provider (`gittag` | `ghrelease` | `ghcr`).
  * `-v, --version string`: Match prefix filter (e.g., query `--version 1.2` returns `1.2.8` rather than `2.0.1`).
* **Output Format:**
  * `text`: Prints the raw highest version.
  * `json`: Emits:
    ```json
    {
      "version": "v1.2.3",
      "normalized": "1.2.3",
      "rendered": "v1.2.3"
    }
    ```

---

### Command: `verge bump`
Calculates and returns the logically bumped next version.

* **Flags:**
  * `-t, --type string`: Override version format.
  * `-p, --provider string`: Override active provider.
  * `-v, --version string`: Bypass fetching and directly parse a static string (e.g., `verge bump --version 1.2.3`).
  * `--prefix string`: Prefix filter to apply when fetching the latest version.
  * `--kind string`: The type of semantic bump to execute (`major` | `minor` | `patch` | `prerelease` | `final`).
  * `--stage string`: The target pre-release stage label (`dev` | `a` | `b` | `rc`).
  * `-s, --sequence string`: Static sequence value to override calculators.
* **Output Format:**
  * `text`: Prints the next version string.
  * `json`: Emits:
    ```json
    {
      "kind": "prerelease",
      "to": "1.2.4-dev.1",
      "rendered": "1.2.4-dev.1"
    }
    ```

---

### Command: `verge ini-config`
Generates boilerplate configurations to assist in fast workspace setup.

* **Flags:**
  * `--template string`: Template identifier (`gittag-semver` | `gittag-vsemver` | `ghrelease-semver`).
* **Conflict Prevention:** If `.verge.yaml` already exists, writes to a separate suffix file (e.g., `.verge.gittag-vsemver.yaml` or `.verge.generated.yaml`) to protect existing configurations.

---

## 3. Configuration Schema Spec (`.verge.yaml`)

Verge loads configuration from a single `.verge.yaml` file parsed across Go YAML specifications. The schema is defined as:

```yaml
# Supported format logic: semver | vsemver | pep440
version_type: semver

# Bumping Defaults
default:
  bump_kind: prerelease      # (major | minor | patch | prerelease | final)
  prerelease_stage: dev      # (dev | a | b | rc)

# Sequence Generator
sequence:
  type: increment            # (increment | filehash | passed)
  targets:                   # Directory or file paths (for filehash)
    - "./go.mod"
  length: 7                  # Truncate length (for filehash)

# History Source (Exactly one provider type configuration)
provider:
  type: gittag               # Active provider (gittag | ghrelease | ghcr)
  gittag:
    repo_dir: "."
    include_prerelease: true
```

---

## 4. Supported Formats & Bumping Logic

### Version Formats
1. **`semver`**: Standard semantic versioning. Matches:
   `^(\d+)\.(\d+)\.(\d+)(?:-(dev|alpha|beta|rc|a|b)\.([a-zA-Z0-9_-]+))?(?:\+[a-zA-Z0-9._-]+)?$`
2. **`vsemver`**: Semantic versioning with a required `v` prefix. Matches:
   `^v(\d+)\.(\d+)\.(\d+)(?:-(dev|alpha|beta|rc|a|b)\.([a-zA-Z0-9_-]+))?(?:\+[a-zA-Z0-9._-]+)?$`
3. **`pep440`**: Python PEP 440 specification. Matches:
   `^(\d+)\.(\d+)\.(\d+)(?:(dev|a|alpha|b|beta|rc)(\d+))?$`

### Bumping Transition Table

When performing a pre-release bump, components and sequence increments behave according to the following logic:

| Base Version | Bumping Kind | Requested Stage | Computed Version Output | Sequence Rule Applied |
| :--- | :--- | :--- | :--- | :--- |
| `1.2.3` | `prerelease` | `dev` | `1.2.4-dev.1` | Resets sequence to initial `1` |
| `1.2.4-dev.1` | `prerelease` | `dev` | `1.2.4-dev.2` | Auto-increments sequence |
| `1.2.4-dev.2` | `prerelease` | `rc` | `1.2.4-rc.1` | Transition resets sequence to `1` |
| `1.2.4-rc.1` | `final` | — | `1.2.4` | Promotes core, wipes sequence |
| `1.2.4-rc.1` | `minor` | — | `1.3.0` | Bumps minor, wipes sequence |

---

## 5. Exit Codes Specification

To ensure robust error handling in automation, Verge mandates the following exit codes:

* **`0`**: Success. Command completed cleanly; output is printed to `stdout`.
* **`1`**: Runtime/Command Error. Network request failed, target provider was unreachable, or invalid runtime flag supplied.
* **`2`**: Configuration Error. Syntax errors in `.verge.yaml` or missing mandatory options.
