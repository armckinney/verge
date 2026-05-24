# Python Package Version Bump (pep440, PyPI source)

Current state: there is no built-in `pypi` provider command path yet.
For a PyPI-based workflow today, fetch current from PyPI externally, then bump with verge using Python rendering.

## 1. Direct CLI Call

```sh
FROM="$(curl -fsSL https://pypi.org/pypi/mypackage/json | jq -r '.info.version')"
verge bump --from "$FROM" --kind patch --ecosystem python
```

Prerelease example:

```sh
verge bump --from 1.2.3 --kind prerelease --stage rc --ecosystem python
# rendered output includes pep440 form (for example 1.2.4rc1)
```

## 2. Using a Config File

`.verge.yaml`

```yaml
version: 1
ecosystem: python

rules:
  defaultBump: patch
  prereleaseStage: rc
```

Use config for defaults/documentation; bump still requires explicit `--from` and `--kind`:

```sh
FROM="$(curl -fsSL https://pypi.org/pypi/mypackage/json | jq -r '.info.version')"
verge bump --from "$FROM" --kind patch --ecosystem python
```

## 3. Required Auth/Env Vars

- Public package reads from `pypi.org` require no auth.
- For private indexes, configure your index credentials for your package tooling (for example `PIP_INDEX_URL`, `PIP_EXTRA_INDEX_URL`, token-based URL, etc.).
