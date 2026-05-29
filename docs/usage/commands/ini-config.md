# `verge ini-config`

The `ini-config` command assists in scaffolding standard, out-of-the-box configuration files for your workspace.

---

## Usage

```bash
verge ini-config [flags]
```

### Command-Specific Flags

* **`--template string`**: The target boilerplate ecosystem template to create. Supported templates include:
  * `gittag-semver` (Default): Sets up standard SemVer parsing query against local Git tags.
  * `gittag-vsemver`: Configures prefixed `vsemver` parsing against local Git tags.
  * `ghrelease-semver`: Configures standard SemVer parsing against the remote GitHub Releases provider.

---

## Conflict Resolution & File Safety

To ensure that existing configuration configurations are **never accidentally overwritten**, Verge checks if `.verge.yaml` already exists in the current directory:

1. **If `.verge.yaml` does NOT exist:** Generates the boilerplate file directly at `.verge.yaml`.
2. **If `.verge.yaml` DOES exist:** 
   * If a specific template was requested (e.g. `--template gittag-vsemver`), creates a separate conflict file named `.verge.gittag-vsemver.yaml`.
   * If no specific template was requested, creates a separate conflict file named `.verge.generated.yaml`.

---

## Examples

### 1. Simple Scaffolding
Create a fresh config:
```bash
$ verge ini-config
Wrote config to .verge.yaml
```

### 2. Scaffold with specific template when `.verge.yaml` already exists
```bash
$ verge ini-config --template gittag-vsemver
Wrote config to .verge.gittag-vsemver.yaml
```

### 3. Generated Boilerplate File Contents (`gittag-vsemver`)
```yaml
version_type: vsemver
provider:
  type: gittag
```
