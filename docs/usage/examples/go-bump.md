# Go Module Version Bump (v-semver, GitHub Releases source)

Current state: this repo has a GitHub Releases provider implementation, but `current/latest` are currently wired to git tags.
For a GitHub Releases source workflow today, fetch the release tag externally, then bump with verge.

## 1. Direct CLI Call

```sh
FROM="$(gh release list --repo org/repo --limit 1 --json tagName --jq '.[0].tagName')"
verge bump --from "$FROM" --kind minor --ecosystem go
```

Alternative without `gh` CLI:

```sh
FROM="$(curl -fsSL https://api.github.com/repos/org/repo/releases | jq -r 'map(select(.draft|not)) | .[0].tag_name')"
verge bump --from "$FROM" --kind minor --ecosystem go
```

## 2. Using a Config File

`.verge.yaml`

```yaml
version: 1
ecosystem: go

format:
  tagPrefix: v

sources:
  github-releases:
    enabled: true
    owner: org
    repo: repo
    includePrerelease: true
```

Use config for intent/documentation; bump still needs explicit `--from`:

```sh
FROM="$(gh release list --repo org/repo --limit 1 --json tagName --jq '.[0].tagName')"
verge bump --from "$FROM" --kind minor --ecosystem go
```

## 3. Required Auth/Env Vars

- `GITHUB_TOKEN` for private repos and higher API limits.

```sh
export GITHUB_TOKEN=ghp_xxx
```
