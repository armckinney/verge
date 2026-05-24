# Terraform Module Version Bump (v-semver, git tag source)

This repository's `bump` command computes the next version from `--from`.
To use git tags as the source of truth, read current first, then bump.

## 1. Direct CLI Call

Known current version:

```sh
verge bump --from v1.2.3 --kind patch --ecosystem terraform
```

From local git tags in one flow:

```sh
FROM="$(verge --format json current --repo-dir . --ecosystem terraform | jq -r '.version')"
verge bump --from "$FROM" --kind patch --ecosystem terraform
```

## 2. Using a Config File

`.verge.yaml`

```yaml
version: 1
ecosystem: terraform

format:
  tagPrefix: v

sources:
  git-tags:
    enabled: true
    includePrerelease: false

rules:
  defaultBump: patch
```

Use config for source/render defaults, then bump from resolved current:

```sh
FROM="$(verge --config .verge.yaml --format json current --repo-dir . | jq -r '.version')"
verge bump --from "$FROM" --kind patch --ecosystem terraform
```

Note: `bump` does not currently read bump kind from config; pass `--kind` explicitly.

## 3. Required Auth/Env Vars

- No auth is required for local `git tag` reads.
- In CI, ensure tags are present in checkout (for example: `fetch-depth: 0` in GitHub Actions).
- `jq` is used above for convenience when extracting JSON fields.
