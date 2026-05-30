# CLI Commands Reference

Verge provides a concise command-line interface optimized for CI/CD environments. It outputs clean stdout logs and handles exit codes precisely so it can be consumed directly by downstream scripts.

---

## Global Options

These flags are persistent and available across all commands:

* **`-c, --config string`**: Explicit path to the configuration file. (Default: `.verge.yaml`). If not passed, Verge also respects the `VERGE_CONFIG` environment variable.
* **`-f, --format string`**: Output serialization format: `text` or `json`. (Default: `"text"`).
  * **`text`**: Exclusively prints the raw calculated version string directly to `stdout`. All logs go to `stderr`.
  * **`json`**: Outputs structured JSON metadata containing rich execution details (e.g., `version`, `previous`, `bump_kind`, `type`, `provider`).
* **`--json`**: Boolean flag to output in JSON format (convenient shortcut for `--format json`).
* **`--field string`**: Extract a single top-level field from the structured command output (extremely useful in JSON/structured mode to pluck variables).
* **`-v, --verbose`**: Enables detailed verbose logging on `stderr`.
* **`--no-cache`**: Disables transparent remote provider caching (e.g., for GitHub Releases or GHCR) and forces a fresh network lookup. By default, remote queries are cached securely in `~/.cache/verge/` for 5 minutes to eliminate GitHub API rate limits.

---

## Exit Codes

Verge uses standard exit codes to signal status to CI automation tasks:

| Exit Code | Classification | Cause / Context |
| :--- | :--- | :--- |
| **`0`** | Success | Command completed successfully; output is written to `stdout`. |
| **`1`** | Runtime/Command Error | Network request failed, target provider unavailable, or invalid runtime flag supplied. |
| **`2`** | Configuration Error | Syntax errors in `.verge.yaml` or missing mandatory provider options. |

---

## Available Commands

* [`verge current`](current.md): Retrieves the currently active locally-tracked version.
* [`verge latest`](latest.md): Queries the active provider for the highest recorded version.
* [`verge bump`](bump.md): Computes and outputs the next version based on configuration and stage instructions.
* [`verge parse`](parse.md): Parses a static version string into structured component metadata.
* [`verge init`](init.md): Generates a boilerplate config file based on typical ecosystem templates.
