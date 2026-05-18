package config

func Default() *Config {
	return &Config{
		Version:   1,
		Ecosystem: "go",
		Format: FormatConfig{
			Input:               "auto",
			Output:              "auto",
			TagPrefix:           "v",
			SequenceInterpreter: "auto",
		},
		Sources: SourcesConfig{
			Precedence: []string{"git-tags"},
			GitTags: GitTagsConfig{
				Enabled:           true,
				Fetch:             false,
				IncludePrerelease: true,
				EcosystemParsing:  "go",
			},
			GitHubReleases: GitHubReleasesConfig{
				Enabled:           false,
				IncludePrerelease: true,
				IncludeDrafts:     false,
			},
			GHCR: GHCRConfig{
				Enabled:           false,
				IncludePrerelease: true,
			},
		},
		Sequence: SequenceConfig{
			HashLength:       7,
			AllowContentHash: true,
			GHBuildPattern:   "gh-",
		},
		Rules: RulesConfig{
			PrereleaseStage:        "dev",
			AllowMajorZeroBreaking: true,
			DefaultBump:            "patch",
		},
		AutoBump: AutoBumpConfig{
			ConventionalCommits: true,
			BreakingTokens:      []string{"BREAKING CHANGE", "!:"},
		},
	}
}
