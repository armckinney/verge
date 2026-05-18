package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/verge/internal/version"
)

const githubAPIBase = "https://api.github.com"

// GitHubReleasesConfig configures the GitHub Releases provider.
type GitHubReleasesConfig struct {
	Owner             string
	Repo              string
	TagPrefix         string
	IncludePrerelease bool
	IncludeDrafts     bool
	CacheTTL          time.Duration
}

type githubRelease struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

type githubReleasesProvider struct {
	cfg    GitHubReleasesConfig
	client *http.Client
	cache  Cache
	etags  map[string]string
}

// NewGitHubReleasesProvider returns a provider backed by the GitHub Releases API.
func NewGitHubReleasesProvider(cfg GitHubReleasesConfig) VersionProvider {
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 5 * time.Minute
	}
	return &githubReleasesProvider{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
		cache:  NewMemoryCache(),
		etags:  make(map[string]string),
	}
}

func (g *githubReleasesProvider) Name() string { return "github-releases" }

func (g *githubReleasesProvider) Fetch(opts QueryOptions) ([]*VersionResult, error) {
	cacheKey := fmt.Sprintf("gh:%s/%s:pre=%v:draft=%v", g.cfg.Owner, g.cfg.Repo,
		opts.IncludePrerelease, g.cfg.IncludeDrafts)

	if cached, ok := g.cache.Get(cacheKey); ok {
		if results, ok := cached.([]*VersionResult); ok {
			return results, nil
		}
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases?per_page=100", githubAPIBase, g.cfg.Owner, g.cfg.Repo)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if etag, ok := g.etags[cacheKey]; ok {
		req.Header.Set("If-None-Match", etag)
	}

	var releases []githubRelease
	err = DefaultRetryPolicy.Do(context.Background(), func() error {
		resp, err := g.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusNotModified:
			// ETag matched; return cached data
			releases = nil
			return nil
		case http.StatusOK:
			if etag := resp.Header.Get("ETag"); etag != "" {
				g.etags[cacheKey] = etag
			}
			return json.NewDecoder(resp.Body).Decode(&releases)
		case http.StatusUnauthorized:
			return &NetworkError{StatusCode: resp.StatusCode, Message: "invalid or missing GITHUB_TOKEN"}
		case http.StatusForbidden:
			return &NetworkError{StatusCode: resp.StatusCode, Message: "rate limit exceeded or insufficient permissions"}
		case http.StatusNotFound:
			return &NetworkError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("repository %s/%s not found", g.cfg.Owner, g.cfg.Repo)}
		case http.StatusTooManyRequests:
			return &NetworkError{StatusCode: resp.StatusCode, Message: "rate limit exceeded"}
		default:
			return &NetworkError{StatusCode: resp.StatusCode, Message: "unexpected response"}
		}
	})
	if err != nil {
		return nil, err
	}

	parser := version.NewParser()
	tagPrefix := opts.TagPrefix
	if tagPrefix == "" {
		tagPrefix = g.cfg.TagPrefix
	}

	var results []*VersionResult
	for _, rel := range releases {
		if rel.Draft && !g.cfg.IncludeDrafts {
			continue
		}

		raw := rel.TagName
		stripped := strings.TrimPrefix(raw, tagPrefix)

		v, err := parser.Parse(stripped)
		if err != nil {
			continue
		}

		if v.IsPrerelease() && !opts.IncludePrerelease && !g.cfg.IncludePrerelease {
			continue
		}

		results = append(results, &VersionResult{
			Version: v,
			Raw:     raw,
			Source:  g.Name(),
		})
	}

	g.cache.Set(cacheKey, results, g.cfg.CacheTTL)
	return results, nil
}
