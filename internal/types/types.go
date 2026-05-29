package types

import "example.com/verge/internal/version"

// VersionType defines the parsing and rendering interface for a semantic version format.
type VersionType interface {
	Name() string
	Parse(input string) (*version.Version, error)
	Render(v *version.Version) string
}

// Get finds the matching VersionType by name.
func Get(name string) VersionType {
	switch name {
	case "semver":
		return &SemVer{}
	case "vsemver":
		return &VSemVer{}
	case "pep440":
		return &PEP440{}
	default:
		return nil
	}
}
