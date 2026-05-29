package providers

import "example.com/verge/internal/version"

// VersionProvider defines the contract for accessing version data from a remote or local source.
type VersionProvider interface {
	Name() string

	// GetLatest retrieves the absolute latest version from the source.
	GetLatest(versionType string) (*version.Version, error)

	// GetLatestSpecific retrieves the latest version matching a specific prefix (e.g., "1", "1.2", "1.2.3-a").
	GetLatestSpecific(versionType string, prefix string) (*version.Version, error)
}
