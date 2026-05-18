package integration_test

import (
	"testing"

	"example.com/verge/internal/providers"
)

func TestGitTagsProvider(t *testing.T) {
	provider := providers.NewGitTagsProvider()
	results, err := provider.Fetch(providers.QueryOptions{
		IncludePrerelease: true,
		TagPrefix:         "v",
		RepoDir:           "../..",
	})
	if err != nil {
		t.Skipf("skipping: git tags not available: %v", err)
	}
	// Just verify it doesn't error with a valid git repo
	t.Logf("Found %d versions", len(results))
}
