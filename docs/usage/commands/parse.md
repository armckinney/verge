# `verge parse`

The `parse` command extracts and outputs structured components (major, minor, patch, stage, sequence) from a static version string. It is highly optimized for CI/CD pipelines to easily query and perform "floating" tags based on the version.

By default, the command auto-detects the version scheme (`semver`, `vsemver`, or `pep440`) if not explicitly specified.

---

## Usage

```bash
verge parse [version] [flags]
```

### Command-Specific Flags

* **`-t, --type string`**: Explicitly specifies the version format to parse against (`semver` | `vsemver` | `pep440`). If omitted, Verge will auto-detect the correct format.

---

## Output Fields

When requested with `--format json` or plucked via `--field`, the following fields are returned:

| Field Name | Type | Description |
| :--- | :--- | :--- |
| **`major`** | `integer` | The major version component. |
| **`minor`** | `integer` | The minor version component. |
| **`patch`** | `integer` | The patch version component. |
| **`stage`** | `string` | The pre-release stage (`dev` \| `a` \| `b` \| `rc` \| `final`). |
| **`sequence`** | `integer\|string\|null` | The pre-release sequence value. |
| **`sequence_type`** | `string` | The sequence type (`numeric` \| `commit-sha` \| `content-hash` \| `build-id` \| `unknown` \| `""`). |
| **`is_prerelease`** | `boolean` | `true` if the version is a pre-release version, `false` otherwise. |
| **`core`** | `string` | The core `major.minor.patch` version string. |
| **`version_type`** | `string` | The matched format scheme (`semver` \| `vsemver` \| `pep440`). |
| **`version`** | `string` | The original input version string. |
| **`normalized`** | `string` | The normalized canonical version representation. |
| **`rendered`** | `string` | The rendered version string according to the scheme's rules. |

---

## Examples

#### 1. Parse Version with Auto-Detection (Default)
By default, the parse command outputs fully structured JSON:
```bash
$ verge parse 1.2.3dev4
{
  "core": "1.2.3",
  "is_prerelease": true,
  "major": 1,
  "minor": 2,
  "normalized": "1.2.3-dev.4",
  "patch": 3,
  "rendered": "1.2.3dev4",
  "sequence": 4,
  "sequence_type": "numeric",
  "stage": "dev",
  "version": "1.2.3dev4",
  "version_type": "pep440"
}
```

#### 2. Get Raw Rendered Version String in Text Mode
If you prefer raw text output, explicitly pass `--format text`:
```bash
$ verge parse 1.2.3dev4 --format text
1.2.3dev4
```

#### 3. Pluck Major Version for Floating Tags in Scripts
Using the global `--field` flag:
```bash
$ verge parse 1.2.3dev4 --field major
1
```

#### 4. Pluck Core and Pre-release Fields
```bash
$ verge parse 1.2.3dev4 --field core
1.2.3

$ verge parse 1.2.3dev4 --field stage
dev
```

#### 5. Integration with jq
Plucking with `jq`:
```bash
$ verge parse 1.2.3dev4 | jq -r '.is_prerelease'
true
```
