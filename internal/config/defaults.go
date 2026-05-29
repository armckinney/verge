package config

func Default() *Config {
	return &Config{
		VersionType: "vsemver",
		Default: DefaultConfig{
			BumpKind:        "prerelease",
			PrereleaseStage: "dev",
		},
		Sequence: SequenceConfig{
			Type:    "increment",
			Targets: []string{},
			Length:  7,
		},
		Provider: ProviderRaw{
			Type: "gittag",
			Raw:  make(map[string]interface{}),
		},
	}
}
