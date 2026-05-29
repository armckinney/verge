package domain_test

import (
	"os"
	"path/filepath"
	"testing"

	"example.com/verge/internal/config"
	"example.com/verge/internal/domain"
	"example.com/verge/internal/types"
)

func TestDomainBump(t *testing.T) {
	cfg := config.Default()
	cfg.VersionType = "semver"

	// 1. Test Increment Bumping
	t.Run("Increment Bumping", func(t *testing.T) {
		opts := domain.BumpOptions{
			VersionStr:      "1.2.3-dev.1",
			PrereleaseStage: "dev",
			BumpKind:        "prerelease",
		}

		bumped, err := domain.Bump(cfg, opts)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parser := types.Get("semver")
		rendered := parser.Render(bumped)
		if rendered != "1.2.3-dev.2" {
			t.Errorf("expected 1.2.3-dev.2, got %s", rendered)
		}
	})

	// 2. Test Different Stage Reset to 1
	t.Run("Different Stage Reset", func(t *testing.T) {
		opts := domain.BumpOptions{
			VersionStr:      "1.2.3-dev.5",
			PrereleaseStage: "rc",
			BumpKind:        "prerelease",
		}

		bumped, err := domain.Bump(cfg, opts)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parser := types.Get("semver")
		rendered := parser.Render(bumped)
		if rendered != "1.2.3-rc.1" {
			t.Errorf("expected 1.2.3-rc.1, got %s", rendered)
		}
	})

	// 3. Test FileHash Bumping
	t.Run("FileHash Bumping", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "domain-bump-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		f1 := filepath.Join(tmpDir, "A.txt")
		os.WriteFile(f1, []byte("verge-test-content"), 0644)

		fileHashCfg := config.Default()
		fileHashCfg.VersionType = "semver"
		fileHashCfg.Sequence = config.SequenceConfig{
			Type:    "filehash",
			Targets: []string{tmpDir},
			Length:  7,
		}

		opts := domain.BumpOptions{
			VersionStr:      "1.2.3-dev.abc1234",
			PrereleaseStage: "dev",
			BumpKind:        "prerelease",
		}

		bumped, err := domain.Bump(fileHashCfg, opts)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parser := types.Get("semver")
		rendered := parser.Render(bumped)
		expectedSuffix := "-dev.bf88455" // sha256 of "verge-test-content" truncated to 7
		if rendered != "1.2.3"+expectedSuffix {
			t.Errorf("expected 1.2.3-dev.bf88455, got %s", rendered)
		}
	})

	// 4. Test PEP440 Parsing & Rendering
	t.Run("PEP440 Bumping Final & Prerelease", func(t *testing.T) {
		pepCfg := config.Default()
		pepCfg.VersionType = "pep440"

		opts := domain.BumpOptions{
			VersionStr:      "1.2.3",
			PrereleaseStage: "dev",
			BumpKind:        "prerelease",
		}

		bumped, err := domain.Bump(pepCfg, opts)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parser := types.Get("pep440")
		rendered := parser.Render(bumped)
		if rendered != "1.2.4dev1" {
			t.Errorf("expected 1.2.4dev1, got %s", rendered)
		}
	})
}
