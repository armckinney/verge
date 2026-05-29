# Version Providers

Providers define the external or local data sources from which Verge extracts tags and version history. Exactly one provider configuration is active at any time, serving as the single source of truth for the application.

---

## 1. Local Git Tags (`gittag`)

The `gittag` provider reads the local `.git` repository tags state. It is the only provider that supports the `verge current` command, as it operates entirely on local context.

### YAML Schema
```yaml
provider:
  type: gittag
  gittag:
    repo_dir: "."                    # Path to the git repository (default: ".")
    include_prerelease: true         # Whether to include prereleases in history (default: false)
```

### Authentication
No authentication is required as it runs local `git` CLI operations under the hood. The system must have `git` installed and accessible in the environment's `PATH`.

---

## 2. GitHub Releases (`ghrelease`)

The `ghrelease` provider queries the GitHub REST API to list and parse tags associated with repository releases.

### YAML Schema
```yaml
provider:
  type: ghrelease
  ghrelease:
    owner: "owner"                   # GitHub organization or user name
    repo: "repo-name"                # Target repository name
    include_prerelease: true         # Include prerelease tagged releases in history (default: false)
    include_drafts: false            # Include drafts in history (default: false)
```

### Authentication & Token Usage
To prevent API rate-limiting or to query private repositories, supply a GitHub token via the standard **`GITHUB_TOKEN`** environment variable.

```bash
$ export GITHUB_TOKEN="ghp_xxxxxxxxxxxxxxxxxxxxx"
$ verge latest
```

---

## 3. GitHub Container Registry (`ghcr`)

The `ghcr` provider queries the GitHub Container Registry (GHCR.io) to fetch tags of a published container image.

### YAML Schema
```yaml
provider:
  type: ghcr
  ghcr:
    package: "owner/image-name"      # Package container identifier (excluding "ghcr.io/")
    include_prerelease: true         # Include prerelease container tags in history (default: false)
```

### Authentication
Container registries require authentication to pull tag listings. Ensure the environment has the **`GITHUB_TOKEN`** variable exported. Verge will handle docker V2 basic authentication handshakes automatically.

---

## Behavior on Missing History

If a provider queries a target (local or remote) and finds **no valid versions** matching the active parsing rules, Verge initiates the **initialization fallback**:
* For `verge latest`, an error is returned.
* For `verge bump`, the CLI safely assumes an initial version state starting at `0.1.0` (or `v0.1.0` / `0.1.0dev1` depending on `version_type`) and runs the bump calculations from there.
