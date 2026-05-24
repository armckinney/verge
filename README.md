# Verge

Version Merge at the bleeding edge.

`verge` is a semantic versioning CLI for parsing, comparing, bumping, and rendering versions across ecosystems.

## Install

Download the latest release and install (Linux example):

```bash
curl -sSL https://github.com/armckinney/verge/releases/latest/download/verge-linux-amd64 -o verge
chmod +x verge
sudo mv verge /usr/local/bin/
verge info
```

Replace `verge-linux-amd64` with the appropriate binary for your platform (e.g., `verge-darwin-amd64` for macOS, `verge-windows-amd64.exe` for Windows). See [GitHub Releases](https://github.com/armckinney/verge/releases) for all available downloads.

Install with Go:

```bash
go install example.com/verge/cmd/verge@latest
verge info
```

## Quick Usage

Parse a version:

```bash
verge parse v1.2.3-rc.1
```

Compare two versions (`10` means left is older, `11` means left is newer):

```bash
verge compare 1.2.3 2.0.0; echo $?
```

Bump a version:

```bash
verge bump --from 1.2.3 --kind minor
```

Get current version from local git tags:

```bash
verge current
```

Use JSON output for scripting:

```bash
verge --format json latest
```

## Documentation

- Usage guide: [docs/usage/usage.md](docs/usage/usage.md)
- Configuration reference: [docs/usage/configuration.md](docs/usage/configuration.md)
- CI/CD notes: [docs/usage/cicd.md](docs/usage/cicd.md)
- Use-case examples: [docs/usage/examples/README.md](docs/usage/examples/README.md)
