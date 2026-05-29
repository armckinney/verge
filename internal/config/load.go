package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(configPath string) (*Config, error) {
	cfg := Default()

	if configPath == "" {
		configPath = os.Getenv("VERGE_CONFIG")
	}
	if configPath == "" {
		configPath = ".verge.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			applyEnvOverrides(cfg)
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file %q: %w", configPath, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", configPath, err)
	}

	applyEnvOverrides(cfg)
	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("VERGE_VERSION_TYPE"); v != "" {
		cfg.VersionType = v
	}
	if v := os.Getenv("VERGE_PROVIDER_TYPE"); v != "" {
		cfg.Provider.Type = v
	}
}

func Validate(cfg *Config) []error {
	var errs []error
	var validTypes = map[string]bool{"semver": true, "vsemver": true, "pep440": true}
	if !validTypes[cfg.VersionType] {
		errs = append(errs, fmt.Errorf("invalid version_type: %q (must be semver, vsemver, or pep440)", cfg.VersionType))
	}
	if cfg.Provider.Type == "" {
		errs = append(errs, fmt.Errorf("provider type cannot be empty"))
	}
	return errs
}
