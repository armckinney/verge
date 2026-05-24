# Container Version Bump (semver output, static input, hash prerelease)

Current state: there is no `static` provider flag on `bump`, and no built-in hash-sequence bump mode yet.
Use verge for core bumping and rendering, then append your hash suffix.

## 1. Direct CLI Call

Patch bump from a known static version:

```sh
verge bump --from 1.2.3 --kind patch --ecosystem containers
```

Hash-based prerelease workflow:

```sh
NEXT_BASE="$(verge --format json bump --from 1.2.3 --kind patch --ecosystem containers | jq -r '.to')"
HASH="$(cat Dockerfile app/config.yaml | sha256sum | cut -c1-12)"
echo "${NEXT_BASE}-dev.${HASH}"
```

Optional validation:

```sh
verge parse "${NEXT_BASE}-dev.${HASH}" --ecosystem containers
```

## 2. Using a Config File

`.verge.yaml`

```yaml
version: 1
ecosystem: containers

rules:
  defaultBump: patch

sequence:
  hashLength: 12
  allowContentHash: true
```

Use config for defaults/documentation; bump still needs explicit `--from` and `--kind`:

```sh
NEXT_BASE="$(verge --format json bump --from 1.2.3 --kind patch --ecosystem containers | jq -r '.to')"
HASH="$(cat Dockerfile app/config.yaml | sha256sum | cut -c1-12)"
echo "${NEXT_BASE}-dev.${HASH}"
```

## 3. Required Auth/Env Vars

- No auth required for local version computation.
- For registry pushes, configure registry credentials separately (for example `DOCKER_USERNAME`, `DOCKER_PASSWORD`, or registry-specific login).
