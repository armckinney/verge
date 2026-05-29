package integration_test

import (
	"testing"

	"example.com/verge/internal/domain"
	"example.com/verge/internal/config"
)

func TestGitTagsProvider(t *testing.T) {
	cfg := &config.Config{
		VersionType: "vsemver",
		Provider: config.ProviderRaw{
			Type: "gittag",
			Raw:  map[string]interface{}{"repo_dir": "../.."},
		},
	}
	v, err := domain.GetCurrent(cfg, "")
	if err != nil {
		t.Skipf("skipping: git tags not available: %v", err)
	}
	t.Logf("Found version: %v", v.String())
}
