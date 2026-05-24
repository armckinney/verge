# Devcontainer Feature Version Bump (semver, GitHub Releases source)

Current state: `current/latest` are currently wired to git tags, not GitHub Releases.
For a releases-based flow today, fetch the source version externally, then bump with verge.

## 1. Direct CLI Call

```sh
FROM="$(gh release list --repo org/devcontainer-feature --limit 1 --json tagName --jq '.[0].tagName')"
verge bump --from "$FROM" --kind minor --ecosystem semver
```

## 2. Using a Config File

`.verge.yaml`

```yaml
version: 1
ecosystem: semver

sources:
  github-releases:
    enabled: true
    owner: org
    repo: devcontainer-feature
    includePrerelease: true

rules:
  defaultBump: minor
```

Use config to record source intent; bump still requires explicit `--from`:

```sh
FROM="$(gh release list --repo org/devcontainer-feature --limit 1 --json tagName --jq '.[0].tagName')"
verge bump --from "$FROM" --kind minor --ecosystem semver
```

## 3. Required Auth/Env Vars

- `GITHUB_TOKEN` for private repos and higher API limits.

```sh
export GITHUB_TOKEN=ghp_xxx
```
