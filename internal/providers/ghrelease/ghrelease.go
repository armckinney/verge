package ghrelease

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"example.com/verge/internal/providers"
	"example.com/verge/internal/types"
	"example.com/verge/internal/version"
)

const githubAPIBase = "https://api.github.com"

// Config represents the raw config block passed for this provider.
type Config struct {
	Owner             string        `yaml:"owner"`
	Repo              string        `yaml:"repo"`
	IncludePrerelease bool          `yaml:"include_prerelease"`
	IncludeDrafts     bool          `yaml:"include_drafts"`
}

type Provider struct {
	config Config
	client *http.Client
}

func NewProvider(cfg Config) providers.VersionProvider {
	return &Provider{
		config: cfg,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *Provider) Name() string {
	return "ghrelease"
}

type githubRelease struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

func (p *Provider) fetchAndParse(versionType string) ([]*version.Version, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases?per_page=100", githubAPIBase, p.config.Owner, p.config.Repo)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Simplification: Omitting retry and cache logic for brevity per spec constraints (or using it if required later)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: %d", resp.StatusCode)
	}

	var releases []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	parser := types.Get(versionType)
	if parser == nil {
		return nil, fmt.Errorf("parser not found: %s", versionType)
	}

	var results []*version.Version
	for _, rel := range releases {
		if rel.Draft && !p.config.IncludeDrafts {
			continue
		}

		v, err := parser.Parse(rel.TagName)
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
	
	for _, v := range results {
		if strings.HasPrefix(v.Original, prefix) || strings.HasPrefix(v.Core(), prefix) {
			return v, nil
		}
		
		if versionType == "vsemver" && strings.HasPrefix(v.Original, "v"+prefix) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no version found matching prefix %s", prefix)
}
