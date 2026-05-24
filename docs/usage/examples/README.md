# Usage Examples Index

These examples are aligned with the current CLI implementation.

## Use Cases

- Terraform modules (v-semver, git tags): `terraform-bump.md`
- Go modules (v-semver, GitHub Releases source intent): `go-bump.md`
- GitHub Actions/reusable workflows (semver, git tags): `github-actions-bump.md`
- Containers (semver output + manual hash prerelease workflow): `container-bump.md`
- Devcontainer features (semver, GitHub Releases source intent): `devcontainer-feature-bump.md`
- Python packages (pep440 rendering + external PyPI source lookup): `python-bump.md`

## Important Current Limitations

- `verge bump` requires `--from` and `--kind`; it does not read these from config.
- `verge current` and `verge latest` are currently wired to git tags.
- GitHub Releases, GHCR, and PyPI source examples use external fetch commands for now.

## Recommended Pattern

1. Resolve current version from your source of truth (local git tags or external API).
2. Run `verge bump --from <current> --kind <kind> --ecosystem <ecosystem>`.
3. Use `--format json` and `jq` in automation.
