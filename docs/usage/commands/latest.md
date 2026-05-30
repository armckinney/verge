# `verge latest`

The `latest` command queries your active tracking provider, processes the returned history across the configured `version_type` comparator rules, and outputs the highest definitive version.

Unlike `current`, the `latest` command supports remote providers (`ghrelease`, `ghcr`) and handles HTTP rate-limiting, authentication, and caching automatically under the hood.

---

## Usage

```bash
verge latest [flags]
```

### Command-Specific Flags

* **`-t, --type string`**: Overrides the `version_type` setting in `.verge.yaml` (`semver` | `vsemver` | `pep440`).
* **`-p, --provider string`**: Overrides the active `provider.type` (`gittag` | `ghrelease` | `ghcr`).
* **`-v, --version string`**: Applies a prefix filter when querying. This allows you to find the latest version within a specific major or minor release line (e.g., finding the latest `1.2` release).
* **`--provider-config strings`**: Comma-separated `key=value` pairs to provide fine-grained inline overrides of provider settings (e.g., `--provider-config repo_dir="/tmp/project"`).

---

## Examples

#### 1. Global Latest (Default)
```bash
$ verge latest
v2.4.1
```

#### 2. Specific Minor Release Filtering
Filter tag lists to return the highest version starting with prefix `1.2`:
```bash
$ verge latest --version 1.2
v1.2.8
```

#### 3. Overriding Provider on the Fly
Query the container tags in GitHub Container Registry directly, bypassing the local git tag state:
```bash
$ verge latest --provider ghcr --type semver
1.4.2-dev.c7bb24f
```

#### 4. JSON Mode
```bash
$ verge latest --format json
{
  "version": "v2.4.1",
  "normalized": "2.4.1",
  "rendered": "v2.4.1"
}
```

#### 5. Fine-grained Overrides
Query the latest stable tags by overriding `include_prerelease` on the fly:
```bash
$ verge latest --provider-config include_prerelease=false,repo_dir="/tmp/project"
v1.2.3
```
