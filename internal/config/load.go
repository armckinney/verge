package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(configPath string) (*Config, error) {
	cfg := Default()

	if configPath == "" {
		configPath = os.Getenv("VERCTL_CONFIG")
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
	if v := os.Getenv("VERCTL_ECOSYSTEM"); v != "" {
		cfg.Ecosystem = v
	}
	if v := os.Getenv("VERCTL_FORMAT_OUTPUT"); v != "" {
		cfg.Format.Output = v
	}
	if v := os.Getenv("VERCTL_TAG_PREFIX"); v != "" {
		cfg.Format.TagPrefix = v
	}
}

func Validate(cfg *Config) []error {
	var errs []error
	if cfg.Version != 1 {
		errs = append(errs, fmt.Errorf("unsupported config version: %d", cfg.Version))
	}
	return errs
}
