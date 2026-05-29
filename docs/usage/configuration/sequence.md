# Dynamic Sequence Calculators

Prereleases require a unique suffix sequence identifier to make them distinct from final releases and other prereleases. Verge uses a **Sequence Calculator Engine** to dynamically resolve these identifiers during the `verge bump` execution.

---

## 1. `increment` (Default Numeric Auto-Increment)

Automatically increments the numeric sequence integer by `+1` based on the previous version retrieved from the provider.

### Config YAML Block
```yaml
sequence:
  type: increment
```

### Bumping Rules:
* **Final to Prerelease:** If the previous version is a final release (e.g., `1.2.3`), and you bump to a prerelease (e.g., `--kind prerelease`), the sequence resets and initializes at `1` (e.g., `1.2.4-dev.1`).
* **Stage Change:** If the prerelease stage shifts (e.g., `1.2.3-dev.4` $\rightarrow$ `--stage rc`), the sequence resets and initializes at `1` (e.g., `1.2.3-rc.1`).
* **Same Stage:** If the stage is identical, the value is incremented by `+1` (e.g., `1.2.3-dev.4` $\rightarrow$ `1.2.3-dev.5`).

---

## 2. `filehash` (Deterministic Directory/File Hashing)

Computes a SHA-256 hash of specified target files or directories and truncates it to the requested length. This is extremely useful in CI/CD pipelines to calculate unique tags for container layers, devcontainer features, or terraform configurations based on actual source file changes.

### Config YAML Block
```yaml
sequence:
  type: filehash
  targets:                   # Files or folders to scan recursively
    - "./Dockerfile"
    - "./src/assets"
  length: 7                  # Truncate hash length (default: 32)
```

### Characteristics:
* Deterministic: If the target files have not changed, the calculated sequence will remain identical, preventing redundant container builds or tag registrations.
* Safe: Unlike numeric auto-increment, `filehash` does not require querying historical counters; it hashes current local workspace state.

---

## 3. `passed` (CLI Suffix Injection)

Injects an arbitrary string passed directly via the CLI interface during runtime (e.g., a CI build number, git commit SHA, or pipeline execution ID).

### Config YAML Block
```yaml
sequence:
  type: passed
```

### Runtime Usage:
When `type: passed` is configured, you must supply the sequence string using the `--sequence` flag, or the execution will return an error:

```bash
$ verge bump --kind prerelease --sequence "$GITHUB_RUN_NUMBER"
1.2.3-dev.42
```
