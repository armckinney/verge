package version

// Parser parses a version string into a Version struct.
type Parser interface {
	Parse(input string) (*Version, error)
}

// Normalizer normalizes a version (e.g., converts string sequences to int when numeric).
type Normalizer interface {
	Normalize(v *Version) *Version
}

// Comparator compares two versions.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
type Comparator interface {
	Compare(a, b *Version) int
}

// Renderer renders a version to a string for a given ecosystem.
type Renderer interface {
	Render(v *Version) string
}

// BumpKind represents the type of version bump.
type BumpKind string

const (
	BumpMajor      BumpKind = "major"
	BumpMinor      BumpKind = "minor"
	BumpPatch      BumpKind = "patch"
	BumpPrerelease BumpKind = "prerelease"
	BumpFinal      BumpKind = "final"
)

// Bumper bumps a version according to the kind.
type Bumper interface {
	Bump(v *Version, kind BumpKind, stage Stage) (*Version, error)
}
