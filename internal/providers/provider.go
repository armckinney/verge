package providers

import "example.com/template-go/internal/version"

// QueryOptions controls how versions are fetched.
type QueryOptions struct {
	IncludePrerelease bool
	TagPrefix         string
	RepoDir           string
}

// VersionResult is a parsed version with metadata.
type VersionResult struct {
	Version *version.Version
	Raw     string
	Source  string
}

// VersionProvider fetches versions from a source.
type VersionProvider interface {
	Name() string
	Fetch(opts QueryOptions) ([]*VersionResult, error)
}
