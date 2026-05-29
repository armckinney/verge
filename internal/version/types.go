package version

// BumpKind represents the type of version bump.
type BumpKind string

const (
	BumpMajor      BumpKind = "major"
	BumpMinor      BumpKind = "minor"
	BumpPatch      BumpKind = "patch"
	BumpPrerelease BumpKind = "prerelease"
	BumpFinal      BumpKind = "final"
)

// Comparator compares two versions.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
type Comparator interface {
	Compare(a, b *Version) int
}

// Bumper bumps a version.
type Bumper interface {
	Bump(v *Version, kind BumpKind, stage Stage) (*Version, error)
}
