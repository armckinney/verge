# `verge current`

The `current` command retrieves the currently active locally-tracked version based on your active configuration.

> [!IMPORTANT]
> To enforce strict local context integrity, the `current` command **only supports local tracking providers** (e.g., `gittag`). If your active provider is remote (such as `ghrelease` or `ghcr`), `verge current` will immediately return an exit code `1` with a descriptive error on `stderr`.

---

## Usage

```bash
verge current [flags]
```

### Examples

#### 1. Text Mode (Default)
In text mode, only the exact parsed tag is returned:
```bash
$ verge current
v1.2.3
```

#### 2. JSON Mode
```bash
$ verge current --format json
{
  "version": "v1.2.3",
  "normalized": "1.2.3",
  "rendered": "v1.2.3"
}
```

#### 3. Extraction via `--field`
Retrieve just the normalized (un-prefixed) semver string:
```bash
$ verge current --format json --field normalized
1.2.3
```

---

## Errors & Edge Cases

* **No git tags exist in repository:** Returns an error `current failed: no valid versions found` with exit code `1`.
* **Remote provider configured:**
  ```bash
  $ verge current
  Error: current command only supports local tracking providers, got: ghrelease
  $ echo $?
  1
  ```
