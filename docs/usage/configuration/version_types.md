# Version Types

The `version_type` setting determines the grammar, parsing constraints, sorting comparator rules, and rendering formats for all version strings processed by Verge.

Verge supports three primary version types:

| Type Name | Format Syntax | Stage Keywords | Example Final | Example Prerelease |
| :--- | :--- | :--- | :--- | :--- |
| **`semver`** | Standard Semantic Versioning | `dev`, `alpha`, `beta`, `rc` | `1.2.3` | `1.2.3-dev.42` |
| **`vsemver`** | Prefixed Semantic Versioning | `dev`, `alpha`, `beta`, `rc` | `v1.2.3` | `v1.2.3-dev.42` |
| **`pep440`** | Python PEP 440 Standard | `dev`, `a`, `b`, `rc` | `1.2.3` | `1.2.3dev42` |

---

## 1. `semver` (Standard Semantic Versioning)

Adheres strictly to the SemVer 2.0.0 specification.
* **Regex Pattern:** `^(\d+)\.(\d+)\.(\d+)(?:-(dev|alpha|beta|rc|a|b)\.([a-zA-Z0-9_-]+))?(?:\+[a-zA-Z0-9._-]+)?$`
* **Prelease Stage rendering:** Prefixed with `-`, stage label expanded (`dev` $\rightarrow$ `dev`, `alpha` $\rightarrow$ `alpha`), followed by a dot `.` and sequence identifier.
* **Example:** `1.4.0-alpha.3`

---

## 2. `vsemver` (Prefixed Semantic Versioning)

Identical to standard `semver`, but requires and enforces a lowercase `v` prefix. This is widely used in ecosystems like Go modules and Terraform providers.
* **Regex Pattern:** `^v(\d+)\.(\d+)\.(\d+)(?:-(dev|alpha|beta|rc|a|b)\.([a-zA-Z0-9_-]+))?(?:\+[a-zA-Z0-9._-]+)?$`
* **Example:** `v1.4.0-rc.1`

---

## 3. `pep440` (Python PEP 440 Standard)

Adheres to Python's PEP 440 version standard, omitting punctuation separation between the patch component, prerelease stage, and sequence identifier.
* **Regex Pattern:** `^(\d+)\.(\d+)\.(\d+)(?:(dev|a|alpha|b|beta|rc)(\d+))?$`
* **Stage Keywords mapping:**
  * `dev` $\rightarrow$ `dev`
  * `a` / `alpha` $\rightarrow$ `a`
  * `b` / `beta` $\rightarrow$ `b`
  * `rc` $\rightarrow$ `rc`
* **Rendering format:** `<major>.<minor>.<patch><stage><sequence>`
* **Example:** `2.3.0b12` (equivalent to beta 12)

---

## Comparison Order & Stages

Verge sorts versions according to component weight, from major down to the sequence level. The internal prerelease stage ordering weight is defined as:

$$\text{dev} < \text{alpha (a)} < \text{beta (b)} < \text{rc} < \text{final release}$$

For example, the following PEP440 sequence is sorted from lowest to highest:
1. `1.0.0dev1`
2. `1.0.0a1`
3. `1.0.0a2`
4. `1.0.0b1`
5. `1.0.0rc1`
6. `1.0.0` (final release)
