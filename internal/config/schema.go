package config

// Config represents the root configuration of `.verge.yaml`.
type Config struct {
	VersionType string         `yaml:"version_type"`
	Default     DefaultConfig  `yaml:"default"`
	Sequence    SequenceConfig `yaml:"sequence"`
	Provider    ProviderRaw    `yaml:"provider"`
	NoCache     bool           `yaml:"-"`
}

// DefaultConfig holds default values for bumping logic.
type DefaultConfig struct {
	BumpKind        string `yaml:"bump_kind"`
	PrereleaseStage string `yaml:"prerelease_stage"`
}

// SequenceConfig specifies how to compute sequence strings dynamically.
type SequenceConfig struct {
	Type    string   `yaml:"type"`
	Targets []string `yaml:"targets"`
	Length  int      `yaml:"length"`
}

// ProviderRaw holds the raw, unparsed provider section of the config.
// The `Type` acts as a discriminator so we know which provider to pass the remainder of the block to.
type ProviderRaw struct {
	Type string                 `yaml:"type"`
	Raw  map[string]interface{} `yaml:",inline"`
}
