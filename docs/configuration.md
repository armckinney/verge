# Configuration Reference

Full reference for `.verctl.yaml` configuration options.

---

## File Location

`verctl` searches for the config file in this order:

1. Path given by `--config <path>`
2. `.verctl.yaml` in the current working directory

If no config file is found, all defaults apply.

---

## Top-level Fields

| Field       | Type   | Default    | Description                                    |
|-------------|--------|------------|------------------------------------------------|
| `version`   | int    | `1`        | Config schema version (always `1`)             |
| `ecosystem` | string | `v-semver` | Default format scheme for rendering. Accepts canonical names (`v-semver`, `semver`, `pep440`) or ecosystem aliases (`go`, `terraform`, `containers`, `github-actions`, `python`). |

---

## `format`

Controls how versions are parsed and displayed.

```yaml
format:
  input: auto
  output: auto
  tagPrefix: v
  sequenceInterpreter: auto
```

| Field                 | Values                    | Default  | Description                                               |
|-----------------------|---------------------------|----------|-----------------------------------------------------------|
| `input`               | `auto` \| `semver` \| `pep440` | `auto` | Force a specific parsing scheme, or auto-detect         |
| `output`              | `auto` \| `semver` \| `pep440` | `auto` | Force a specific output scheme                          |
| `tagPrefix`           | string                    | `v`      | Prefix stripped from git tags before parsing (e.g. `v`) |
| `sequenceInterpreter` | `auto` \| `numeric` \| `hash`  | `auto` | How to interpret prerelease sequences                   |

**Environment variable override:** `VERCTL_TAG_PREFIX`

---

## `sources`

Configures version providers and their priority order.

```yaml
sources:
  precedence:
    - git-tags
    - github-releases
    - ghcr
  git-tags:
    enabled: true
    fetch: false
    includePrerelease: true
    ecosystemParsing: v-semver
  github-releases:
    enabled: false
    owner: ""
    repo: ""
    includePrerelease: true
    includeDrafts: false
  ghcr:
    enabled: false
    image: ""
    includePrerelease: true
```

### `sources.precedence`

Ordered list of provider names. The first enabled provider that returns results wins.

Available providers: `git-tags`, `github-releases`, `ghcr`

### `sources.git-tags`

| Field               | Type    | Default | Description                                             |
|---------------------|---------|---------|---------------------------------------------------------|
| `enabled`           | bool    | `true`  | Enable this provider                                    |
| `fetch`             | bool    | `false` | Run `git fetch --tags` before listing                   |
| `includePrerelease` | bool    | `true`  | Include prerelease tags in results                      |
| `ecosystemParsing`  | string  | `v-semver` | Format scheme hint for tag detection (`v-semver`, `semver`, `pep440`, or ecosystem alias) |

### `sources.github-releases`

Fetches versions from the GitHub Releases API.

| Field               | Type    | Default | Description                                          |
|---------------------|---------|---------|------------------------------------------------------|
| `enabled`           | bool    | `false` | Enable this provider                                 |
| `owner`             | string  | `""`    | GitHub repository owner (user or org)                |
| `repo`              | string  | `""`    | GitHub repository name                               |
| `includePrerelease` | bool    | `true`  | Include prerelease releases                          |
| `includeDrafts`     | bool    | `false` | Include draft releases                               |

**Authentication:** Set `GITHUB_TOKEN` environment variable. Without it, requests are unauthenticated (60 req/hr rate limit).

### `sources.ghcr`

Fetches versions from GitHub Container Registry image tags.

| Field               | Type    | Default | Description                                                          |
|---------------------|---------|---------|----------------------------------------------------------------------|
| `enabled`           | bool    | `false` | Enable this provider                                                 |
| `image`             | string  | `""`    | Full image reference (e.g. `ghcr.io/my-org/my-image`)               |
| `includePrerelease` | bool    | `true`  | Include prerelease tags                                              |
| `channelFilter`     | string  | `""`    | Filter to a channel: `rel` (final only), `pr` (prerelease only), `floating` |

**Authentication:** Set `GITHUB_TOKEN` for private images. Public images use anonymous token auth.

---

## `sequence`

Controls how prerelease sequence identifiers are interpreted.

```yaml
sequence:
  hashLength: 7
  allowContentHash: true
  ghBuildPattern: "gh-"
```

| Field              | Type   | Default | Description                                                  |
|--------------------|--------|---------|--------------------------------------------------------------|
| `hashLength`       | int    | `7`     | Number of characters to use from content hashes              |
| `allowContentHash` | bool   | `true`  | Allow hash-based sequences (e.g. `dev.abc1234`)              |
| `ghBuildPattern`   | string | `gh-`   | Prefix pattern that identifies GitHub Actions build IDs      |

---

## `rules`

Governs default bump behavior and release rules.

```yaml
rules:
  prereleaseStage: dev
  allowMajorZeroBreaking: true
  defaultBump: patch
```

| Field                    | Type   | Default | Description                                                         |
|--------------------------|--------|---------|---------------------------------------------------------------------|
| `prereleaseStage`        | string | `dev`   | Default stage for new prereleases (`dev`, `alpha`, `beta`, `rc`)    |
| `allowMajorZeroBreaking` | bool   | `true`  | Allow breaking changes (major bump) on `v0.x` versions             |
| `defaultBump`            | string | `patch` | Default bump kind when none is specified                            |

---

## `autoBump`

Controls automatic bump detection from conventional commits.

```yaml
autoBump:
  conventionalCommits: true
  breakingTokens:
    - "BREAKING CHANGE"
    - "!:"
```

| Field                  | Type     | Default                                | Description                                              |
|------------------------|----------|----------------------------------------|----------------------------------------------------------|
| `conventionalCommits`  | bool     | `true`                                 | Enable conventional commits detection for `--auto`       |
| `breakingTokens`       | []string | `["BREAKING CHANGE", "!:"]`            | Strings in commit messages that signal a major bump      |

---

## Environment Variables

All environment variables override the corresponding config file setting at runtime without modifying the file.

| Variable               | Overrides              | Example                          |
|------------------------|------------------------|----------------------------------|
| `VERCTL_ECOSYSTEM`     | `ecosystem`            | `VERCTL_ECOSYSTEM=pep440`        |
| `VERCTL_FORMAT_OUTPUT` | `format.output`        | `VERCTL_FORMAT_OUTPUT=json`      |
| `VERCTL_TAG_PREFIX`    | `format.tagPrefix`     | `VERCTL_TAG_PREFIX=""`           |
| `GITHUB_TOKEN`         | *(provider auth)*      | `GITHUB_TOKEN=ghp_xxxx`          |

---

## Complete Example

A `.verctl.yaml` for a Python project published to GHCR and GitHub Releases:

```yaml
version: 1
ecosystem: pep440

format:
  tagPrefix: v
  sequenceInterpreter: auto

sources:
  precedence:
    - git-tags
    - github-releases
  git-tags:
    enabled: true
    includePrerelease: true
  github-releases:
    enabled: true
    owner: my-org
    repo: my-python-app
    includePrerelease: true

rules:
  prereleaseStage: rc
  defaultBump: minor

autoBump:
  conventionalCommits: true
  breakingTokens:
    - "BREAKING CHANGE"
    - "!:"
```
