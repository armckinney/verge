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

const ghcrAPIBase = "https://ghcr.io/v2"

// GHCRConfig configures the GitHub Container Registry provider.
type GHCRConfig struct {
	Image             string // e.g. "ghcr.io/org/repo"
	TagPrefix         string
	IncludePrerelease bool
	ChannelFilter     string // "rel" | "pr" | "floating" | "" (all)
	CacheTTL          time.Duration
}

type ghcrTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ghcrProvider struct {
	cfg    GHCRConfig
	client *http.Client
	cache  Cache
}

// NewGHCRProvider returns a provider backed by GHCR tag listing.
func NewGHCRProvider(cfg GHCRConfig) VersionProvider {
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 5 * time.Minute
	}
	return &ghcrProvider{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
		cache:  NewMemoryCache(),
	}
}

func (g *ghcrProvider) Name() string { return "ghcr" }

// imageRef parses "ghcr.io/org/repo" → ("org/repo", "ghcr.io")
func imageRef(image string) (name, registry string) {
	image = strings.TrimPrefix(image, "ghcr.io/")
	return image, "ghcr.io"
}

// ghcrToken fetches an anonymous or authenticated bearer token for GHCR.
func (g *ghcrProvider) ghcrToken(imageName string) (string, error) {
	if pat := os.Getenv("GITHUB_TOKEN"); pat != "" {
		// Use PAT directly as Basic auth token
		return pat, nil
	}

	tokenURL := fmt.Sprintf("https://ghcr.io/token?service=ghcr.io&scope=repository:%s:pull", imageName)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, tokenURL, nil)
	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var tok struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", err
	}
	return tok.Token, nil
}

func (g *ghcrProvider) Fetch(opts QueryOptions) ([]*VersionResult, error) {
	cacheKey := fmt.Sprintf("ghcr:%s:chan=%s", g.cfg.Image, g.cfg.ChannelFilter)

	if cached, ok := g.cache.Get(cacheKey); ok {
		if results, ok := cached.([]*VersionResult); ok {
			return results, nil
		}
	}

	imageName, _ := imageRef(g.cfg.Image)
	token, err := g.ghcrToken(imageName)
	if err != nil {
		return nil, fmt.Errorf("authenticating to GHCR: %w", err)
	}

	url := fmt.Sprintf("%s/%s/tags/list", ghcrAPIBase, imageName)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	var tags []string
	err = DefaultRetryPolicy.Do(context.Background(), func() error {
		resp, err := g.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			var body ghcrTagsResponse
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				return err
			}
			tags = body.Tags
			return nil
		case http.StatusUnauthorized:
			return &NetworkError{StatusCode: resp.StatusCode, Message: "unauthorized — set GITHUB_TOKEN for private images"}
		case http.StatusNotFound:
			return &NetworkError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("image %s not found", g.cfg.Image)}
		case http.StatusTooManyRequests:
			return &NetworkError{StatusCode: resp.StatusCode, Message: "rate limit exceeded"}
		default:
			return &NetworkError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status %d", resp.StatusCode)}
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
	for _, tag := range tags {
		if tag == "latest" {
			continue
		}

		stripped := strings.TrimPrefix(tag, tagPrefix)
		v, err := parser.Parse(stripped)
		if err != nil {
			continue
		}

		// Apply channel filter
		switch g.cfg.ChannelFilter {
		case "rel":
			if v.IsPrerelease() {
				continue
			}
		case "pr":
			if !v.IsPrerelease() {
				continue
			}
		case "floating":
			// floating tags are single/double component — our parser requires 3 so these won't parse anyway
			continue
		}

		if v.IsPrerelease() && !opts.IncludePrerelease && !g.cfg.IncludePrerelease {
			continue
		}

		results = append(results, &VersionResult{
			Version: v,
			Raw:     tag,
			Source:  g.Name(),
		})
	}

	g.cache.Set(cacheKey, results, g.cfg.CacheTTL)
	return results, nil
}
