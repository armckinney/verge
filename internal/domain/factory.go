package domain

import (
	"example.com/verge/internal/config"
	"example.com/verge/internal/providers"
	"example.com/verge/internal/providers/ghcr"
	"example.com/verge/internal/providers/ghrelease"
	"example.com/verge/internal/providers/gittag"
	"fmt"
	"gopkg.in/yaml.v3"
)

// NewFromConfig builds a VersionProvider given the raw provider configuration.
func NewFromConfig(raw config.ProviderRaw) (providers.VersionProvider, error) {
	// Re-marshal and unmarshal to parse strictly into the target structs
	data, err := yaml.Marshal(raw.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider config: %w", err)
	}

	switch raw.Type {
	case "gittag":
		var cfg gittag.Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
		return gittag.NewProvider(cfg), nil
	case "ghrelease":
		var cfg ghrelease.Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
		return ghrelease.NewProvider(cfg), nil
	case "ghcr":
		var cfg ghcr.Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
		return ghcr.NewProvider(cfg), nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", raw.Type)
	}
}

// IsLocal tracks whether the provider reads from local state or network.
func IsLocal(raw config.ProviderRaw) bool {
	return raw.Type == "gittag"
}
