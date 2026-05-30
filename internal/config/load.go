package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

	if errs := Validate(cfg); len(errs) > 0 {
		return nil, fmt.Errorf("config validation failed: %v", errs[0])
	}

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

// ParseOverrides converts a slice of "key=value" strings into a strongly-typed map.
func ParseOverrides(pairs []string) map[string]interface{} {
	out := make(map[string]interface{})
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])

		// Handle Boolean Conversion
		if boolVal, err := strconv.ParseBool(valStr); err == nil {
			out[key] = boolVal
			continue
		}

		// Handle Integer Conversion
		if intVal, err := strconv.Atoi(valStr); err == nil {
			out[key] = intVal
			continue
		}

		// Fallback to String
		out[key] = valStr
	}
	return out
}

// MergeOverrides merges the parsed overrides into the ProviderRaw config block.
func MergeOverrides(raw *ProviderRaw, overrides map[string]interface{}) {
	if len(overrides) == 0 {
		return
	}
	if raw.Raw == nil {
		raw.Raw = make(map[string]interface{})
	}

	// If there is a nested sub-map for the active provider (e.g., raw.Raw["gittag"])
	if nested, ok := raw.Raw[raw.Type].(map[string]interface{}); ok {
		for k, v := range overrides {
			nested[k] = v
		}
		return
	}

	// Otherwise, merge flatly
	for k, v := range overrides {
		raw.Raw[k] = v
	}
}
