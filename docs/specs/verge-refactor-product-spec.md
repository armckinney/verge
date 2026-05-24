# Architecture
- CLI has primary entrypoint / interface
- Commands have isolated implementation 
  - easy to understand where to go to modify a commands behavior/logic
  - leverage abstraction where possible and typical software engineering principles / OOP style development
- config is modularized / objectified
  - easy to understand how to modify config models/schemas and loading/parsing processes
  - easy to extend to additional types of config and config parameters
  - easy to find/understand the default configs/behavior of the cli 
- providers have an abstraction/interface layer and are easy to identify how to add additional providers and modify existing ones; they are composable pieces that allow users to extend the integrations of the core versioning app
  - providers are backend implementations for versioning, this is where most of the implementation logic goes
  - these are for example:
    - gittag (default)
    - ghrelease
    - ghcr
    - Can be extended to: pypi, azure blob storage, artifactory, dockerhub
  - providers should also implement their own config which drives them (this doesn't have to be solely in a single file, whatever makes sense to implement, but it should be easy/obvious for users to extend the functionality in src code without having to hunt down files/directories)
  - providers can be a collection of implemented interfaces if it makes composition/code maintainability easier; Ideally a provider should basically be "self-contained" within a module that implements a "core" provider abstraction interfaces;
  - these should be organized by either single files (probably too complex for this), or by module directories, to make it obvious provider code is where
  - new providers should have to implement everything from getting the latest version of a specific version, config values (parsing can use the core config parser)
- version types (vsemver, semver, pep440) should implement how to parse versions and iterate on those
  - version types from a code implementation perspective should be apparent and implement an abstraction / interface in order to make it easy to extend the cli with additional version types


# Implicit Behaviors
- Arguments should be considered in this order, in and be overwriten by subsequent layers
  - cli defaults (from cli src code)
  - config file (.verge.yaml)
  - cli command args (passed as flags like, --kind prerelease)
- all command configs should be able to be passed like flags/arguments, which could make provider arguments extensive; generally though users should not do this and implement a .verge.yaml file instead for provider config
  - providers should have their own config section in the yaml file, i.e. "provider.gittag", "provider.ghcr" to help identify their configuration separate from the application config
- CLI should assume that based on given command if the version doesn't exist, then it should initialize it, i.e.
  - for first time use, it should return a new version like 0.1.0
  - for bumping a prerelease verson to a new stage like a -> b, 1.2.3-a1
- it should also be implied that if a provider's config is identified in the yaml file, then that is what is intended to be used as the single source of truth

# Use Cases

General Use Case Context:
- get current version (local tracking providers only - i.e. gittags/ENV Vars)
  - this should be an optional implementation for providers
- get latest version (could be remote or local depending on provider)
- get latest version of specific major/minor/prerelease (could be remote or local depending on provider)
- (bump) return version bump of latest version:
  - for various kinds: major, minor, patch, prerelease
  - for various prerelease stages: dev, a, b, rc
  - for a given version: i.e. get bump of current version 1.2.3-dev
    - default would be to get the global latest version of specific kind if version is not passed
  - calculated from static versions/sequences: i.e. get bump of dev prerelease version 1.2.3 with sequence filehash of ./Dockerfile
    - also support passed sequence values like build numbers: i.e. "42"
    - default is for increment: i.e. latest determined version sequence +1

Specific Use Case Context Version Releases Scope (this is just to identify needs, not what verge actually needs to implement, implementation should be generalized and not implemented for these specific use cases):
- Containers: 
  - PR: major.minor.patch-stage.<hash>
  - REL: major.minor.patch
  - latest
  - floating: semver
- Devcontainer Features:
  - semver
  - latest
- Python
  - PR: <major>.<minor>.<patch><stage><sequence>
  - REL: <major>.<minor>.<patch>
- Go (v-semver)
  - PR: v<major>.<minor>.<patch>-<stage>.<sequence>
  - REL: v<major>.<minor>.<patch>
- Terraform (v-semver)
- GitHub Actions/RWF: semver

# Supported Version Functionality

Version Context:
- prefix
  - v
- major
- minor
- patch
- stage
  - dev, a, b, rc (note semver is alpha/beta instead of a/b, pep440 is a/b)
- sequence
  - inputs: commit sha, file contents hash, build number
  - interpreted: ghcr, pypi, gh releases, git tag

Version Types:
- semver
- vsemver
- pep440

# Commands & Interfaces

## Global Flags
  -c, --config string   Config file path (default: .verge.yaml)
  -f, --format string   Output format: text or json (default "text")
  -h, --help            help for verge
  -v, --verbose         Enable verbose output

## verge version 
returns the version of verge 
i.e. 0.1.2 or 0.1.2-dev.34

### usage
```
# verge version
1.2.2
```


## verge help
returns the help information of commands possible by verge as well as global flags 

### usage

```
# verge help
verge is a semantic versioning CLI tool for managing and bumping versions across ecosystems.

Usage:
  verge [command]

Available Commands:
  help        Help about any command
  version     Parse a version string
  etc         etc

Flags:
  -h, --help            help for verge
  -v, --verbose         Enable verbose output
  etc                   etc

Use "verge [command] --help" for more information about a command.
```


## verge current
returns the current version as identified by git-tags

### usage

get current local (gittag) version
```
# verge current
1.2.3
```


## verge latest
returns the latest version

Behavior:
Args provide just enough information needed to be able to get the latest version.
Should be able to get the latest versions for:
- global latest
- major latest
- minor latest
- prerelease stage latest

Considerations:
Should default behavior be to get the latest pre-release (how could i make it reliably know the current version?)

Args: 
--type (semver|vsemver|pep440)
--provider (gittags|ghcr|ghrelease)
--format (text|json|rendered)
--version (i.e. 1, 1.2, 1.2.3-dev, 1.2.3-a, 1.2.3-b, 1.2.3-rc)

### usage

get latest global version
```
# verge latest
1.2.3
```

get latest version of a specific version
```
# verge latest --version 1
1.2.3

# verge latest --version 1.3
1.3.7

# verge latest --version 1.2.4-dev
1.2.4-dev42

# verge latest --version 1.2.4-a
1.2.4-a73
```


## verge bump
returns the bumped version

Behavior:
Args provide just enough detail on how to return the bumped version (gets latest & returns bump).

Args:
--type (semver|vsemver|pep440)
--provider (gittags|ghcr|ghrelease)
--format (text|json|rendered)
--prerelease-stage (dev|a|b|rc)
--sequence (filehash|increment|passed value)
--kind (prerelease|patch|minor|major)

### usage
bump version entirely using config file setup, useful for:
- config driven bumping, allowing users to leverage reusable cicd workflows/automation
```
# verge bump
1.2.3-dev42
```

bump version using config file but override kind, useful for situations like:
- prerelease (global/specified)
- prerelease stage progression (global/specified)
- cicd logic drive overrides of default bumps (global/specified)
```
# verge bump --kind prerelease
1.2.3-dev43

# verge bump --kind prerelease --stage rc
1.2.3-rc2

# verge bump --kind prerelease --stage rc --version 1.2.2
1.2.2-rc3

# verge bump --kind minor
1.3.0

# verge bump --kind major
2.0.0
```

## verge ini-config
additional tool to help scaffold out users config files with templated use cases

supports:
- create a separate config file is one already exists (i.e. .verge.template.yaml)
- templates
  - gittag-semver (containers, github actions)
  - gittag-vsemver (terraform)
  - ghrelease-semver (devcontainer features)
  - ghrelease-pep440 (python)
  - ghrealease-vsemver (go)
  - ghcr-semver (containers)

```
# verge ini-config --template gittag-semver
Created config file: .verge.yaml

# verge ini-config --template gittag-semver
Config File Already Exists!
Created config file: .verge.gittag-semver.yaml
```

# Considerations & Open Questions
- Need a way to be able to pass complex config to the cli. For example:
  - provider configuration (i.e. ghcr repository, ghrelease repo)
  - sequence configuration (i.e. filehash length, filehash algorithm)
- Consider identifying "default" values in config that are generally intended to be overwriten in SDLC, like "default.bump.kind=prerelease"

# IMPORTANT Additional information

Potentially a lot of this functionality might already exist in the current implementation. There might just be some small reorganization needed to support clarity for maintainability.

Changes from current implementation:
- No need to support other functionality like commit checking for messages to determine version, we will be more explicit than that
  - this also means no need to check if major bumps require breaking change 
  - this means no need to support "auto" detect
- change "ecosystem" to "type" diction
- no need to support changelog
- no need to support tagPrefix, since that should be inherent to vsemver, if additional prefixes are needed, then they should be implemented in a new provider
- no need to support precedence in providers, there should be a single source of truth that we are imposing
- need to have stronger organization on documentation for usage nad provider extension. currently they are all lumped together in a wall of text. These should be organized better and have easier navigation so users can find what they are looking for faster.
