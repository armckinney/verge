package ecosystems

// EcosystemRenderer renders a version for a specific ecosystem.
type EcosystemRenderer interface {
	Name() string
	Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string
}
