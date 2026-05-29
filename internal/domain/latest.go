package domain

import (
	"fmt"

	"example.com/verge/internal/config"
	"example.com/verge/internal/version"
)

// LatestOptions holds arguments derived from CLI and mapped config.
type LatestOptions struct {
	Prefix string
}

// GetLatest executes the logic for the "latest" action.
func GetLatest(cfg *config.Config, opts LatestOptions) (*version.Version, error) {
	p, err := NewFromConfig(cfg.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to init provider: %w", err)
	}

	if opts.Prefix != "" {
		return p.GetLatestSpecific(cfg.VersionType, opts.Prefix)
	}

	return p.GetLatest(cfg.VersionType)
}
