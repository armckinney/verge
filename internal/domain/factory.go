package domain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"example.com/verge/internal/config"
	"example.com/verge/internal/providers"
	"example.com/verge/internal/providers/ghcr"
	"example.com/verge/internal/providers/ghrelease"
	"example.com/verge/internal/providers/gittag"
	"gopkg.in/yaml.v3"
)

// NewFromConfig builds a VersionProvider given the raw provider configuration.
func NewFromConfig(cfg *config.Config) (providers.VersionProvider, error) {
	// Re-marshal and unmarshal to parse strictly into the target structs
	var rawData interface{} = cfg.Provider.Raw
	if nested, ok := cfg.Provider.Raw[cfg.Provider.Type]; ok {
		rawData = nested
	}
	data, err := yaml.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider config: %w", err)
	}

	var p providers.VersionProvider
	switch cfg.Provider.Type {
	case "gittag":
		var pCfg gittag.Config
		if err := yaml.Unmarshal(data, &pCfg); err != nil {
			return nil, err
		}
		p = gittag.NewProvider(pCfg)
	case "ghrelease":
		var pCfg ghrelease.Config
		if err := yaml.Unmarshal(data, &pCfg); err != nil {
			return nil, err
		}
		p = ghrelease.NewProvider(pCfg)
	case "ghcr":
		var pCfg ghcr.Config
		if err := yaml.Unmarshal(data, &pCfg); err != nil {
			return nil, err
		}
		p = ghcr.NewProvider(pCfg)
	default:
		return nil, fmt.Errorf("unknown provider type: %s", cfg.Provider.Type)
	}

	// Wrap in a caching provider decorator if it's a remote provider
	if !IsLocal(cfg.Provider) {
		cacheKey := fmt.Sprintf("%s:%s", cfg.Provider.Type, hashConfig(cfg.Provider.Raw))
		p = providers.NewCachingProvider(p, cacheKey, cfg.NoCache, 5*time.Minute)
	}

	return p, nil
}

// IsLocal tracks whether the provider reads from local state or network.
func IsLocal(raw config.ProviderRaw) bool {
	return raw.Type == "gittag"
}

func hashConfig(raw map[string]interface{}) string {
	data, _ := json.Marshal(raw)
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h[:8]) // Keep it short and readable
}
