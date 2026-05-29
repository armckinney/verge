package gittag

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"example.com/verge/internal/providers"
	"example.com/verge/internal/types"
	"example.com/verge/internal/version"
)

// Config represents the raw config block passed for this provider.
type Config struct {
	RepoDir           string `yaml:"repo_dir"`
	IncludePrerelease bool   `yaml:"include_prerelease"`
}

type Provider struct {
	config Config
}

func NewProvider(cfg Config) providers.VersionProvider {
	return &Provider{
		config: cfg,
	}
}

func (p *Provider) Name() string {
	return "gittag"
}

func (p *Provider) fetchAndParse(versionType string) ([]*version.Version, error) {
	dir := p.config.RepoDir
	if dir == "" {
		dir = "."
	}

	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch git tags: %w", err)
	}

	parser := types.Get(versionType)
	if parser == nil {
		return nil, fmt.Errorf("parser not found: %s", versionType)
	}

	var results []*version.Version
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		v, err := parser.Parse(line)
		if err != nil {
			continue
		}

		if !p.config.IncludePrerelease && v.IsPrerelease() {
			continue
		}

		results = append(results, v)
	}

	// Sort highest first
	comp := version.NewComparator()
	sort.Slice(results, func(i, j int) bool {
		return comp.Compare(results[i], results[j]) > 0
	})

	return results, nil
}

func (p *Provider) GetLatest(versionType string) (*version.Version, error) {
	results, err := p.fetchAndParse(versionType)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no valid versions found")
	}

	return results[0], nil
}

func (p *Provider) GetLatestSpecific(versionType string, prefix string) (*version.Version, error) {
	results, err := p.fetchAndParse(versionType)
	if err != nil {
		return nil, err
	}
	
	// Example of matching specific prefix: 1.2
	for _, v := range results {
		// Does original match prefix?
		// Note we should consider that prefix might mean matching Major or Minor
		// For now we check if Original starts with prefix or Core starts with prefix.
		if strings.HasPrefix(v.Original, prefix) || strings.HasPrefix(v.Core(), prefix) {
			return v, nil
		}
		
		// If vsemver was used, prefix "1.2" might mean v1.2
		if versionType == "vsemver" && strings.HasPrefix(v.Original, "v"+prefix) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no version found matching prefix %s", prefix)
}
