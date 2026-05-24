# GitHub Actions / Reusable Workflows Version Bump (semver, git tag source)

## 1. Direct CLI Call

Known current version:

```sh
verge bump --from 1.2.3 --kind patch --ecosystem github-actions
```

From local git tags in one flow:

```sh
FROM="$(verge --format json current --repo-dir . --ecosystem github-actions | jq -r '.version')"
verge bump --from "$FROM" --kind patch --ecosystem github-actions
```

## 2. Using a Config File

`.verge.yaml`

```yaml
version: 1
ecosystem: github-actions

format:
  tagPrefix: ""

sources:
  git-tags:
    enabled: true
    includePrerelease: false

rules:
  defaultBump: patch
```

```sh
FROM="$(verge --config .verge.yaml --format json current --repo-dir . | jq -r '.version')"
verge bump --from "$FROM" --kind patch --ecosystem github-actions
```

## 3. Required Auth/Env Vars

- No auth is required when reading local git tags.
- Ensure tags are available in CI checkout.
