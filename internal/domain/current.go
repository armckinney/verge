package domain

import (
	"fmt"

	"example.com/verge/internal/config"
	"example.com/verge/internal/version"
)

// GetCurrent executes the logic for the "current" action.
func GetCurrent(cfg *config.Config, prefix string) (*version.Version, error) {
	if !IsLocal(cfg.Provider) {
		return nil, fmt.Errorf("current command only supports local tracking providers, got: %s", cfg.Provider.Type)
	}

	p, err := NewFromConfig(cfg.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to init provider: %w", err)
	}

	if prefix != "" {
		return p.GetLatestSpecific(cfg.VersionType, prefix)
	}

	return p.GetLatest(cfg.VersionType)
}
