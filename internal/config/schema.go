package config

type Config struct {
	Version   int            `yaml:"version"`
	Ecosystem string         `yaml:"ecosystem"`
	Format    FormatConfig   `yaml:"format"`
	Sources   SourcesConfig  `yaml:"sources"`
	Sequence  SequenceConfig `yaml:"sequence"`
	Rules     RulesConfig    `yaml:"rules"`
	AutoBump  AutoBumpConfig `yaml:"autoBump"`
}

type FormatConfig struct {
	Input               string `yaml:"input"`
	Output              string `yaml:"output"`
	TagPrefix           string `yaml:"tagPrefix"`
	SequenceInterpreter string `yaml:"sequenceInterpreter"`
}

type SourcesConfig struct {
	Precedence     []string             `yaml:"precedence"`
	GitTags        GitTagsConfig        `yaml:"git-tags"`
	GitHubReleases GitHubReleasesConfig `yaml:"github-releases"`
	GHCR           GHCRConfig           `yaml:"ghcr"`
}

type GitTagsConfig struct {
	Enabled           bool   `yaml:"enabled"`
	Fetch             bool   `yaml:"fetch"`
	IncludePrerelease bool   `yaml:"includePrerelease"`
	EcosystemParsing  string `yaml:"ecosystemParsing"`
}

type GitHubReleasesConfig struct {
	Enabled           bool   `yaml:"enabled"`
	Owner             string `yaml:"owner"`
	Repo              string `yaml:"repo"`
	IncludePrerelease bool   `yaml:"includePrerelease"`
	IncludeDrafts     bool   `yaml:"includeDrafts"`
}

type GHCRConfig struct {
	Enabled           bool   `yaml:"enabled"`
	Image             string `yaml:"image"`
	IncludePrerelease bool   `yaml:"includePrerelease"`
}

type SequenceConfig struct {
	HashLength       int    `yaml:"hashLength"`
	AllowContentHash bool   `yaml:"allowContentHash"`
	GHBuildPattern   string `yaml:"ghBuildPattern"`
}

type RulesConfig struct {
	PrereleaseStage        string `yaml:"prereleaseStage"`
	AllowMajorZeroBreaking bool   `yaml:"allowMajorZeroBreaking"`
	DefaultBump            string `yaml:"defaultBump"`
}

type AutoBumpConfig struct {
	ConventionalCommits bool     `yaml:"conventionalCommits"`
	BreakingTokens      []string `yaml:"breakingTokens"`
}
