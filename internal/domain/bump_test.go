package domain_test

import (
	"os"
	"os/exec"
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

func setupTestGitRepo(t *testing.T, tags []string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "verge-domain-bump-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Helper to run command in temp dir
	runCmd := func(name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatalf("command failed: %s %v: %v", name, args, err)
		}
	}

	runCmd("git", "init")
	runCmd("git", "config", "user.name", "Test User")
	runCmd("git", "config", "user.email", "test@example.com")

	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	runCmd("git", "add", "file.txt")
	runCmd("git", "commit", "-m", "initial commit")

	for _, tag := range tags {
		runCmd("git", "tag", tag)
	}

	return dir
}

func TestStableBumpIgnoresPrerelease(t *testing.T) {
	tags := []string{"v1.2.3", "v1.2.4-dev.5"}
	dir := setupTestGitRepo(t, tags)
	defer os.RemoveAll(dir)

	t.Run("nested provider config - include_prerelease true ignored on stable bump", func(t *testing.T) {
		cfg := &config.Config{
			VersionType: "vsemver",
			Provider: config.ProviderRaw{
				Type: "gittag",
				Raw: map[string]interface{}{
					"gittag": map[string]interface{}{
						"repo_dir":           dir,
						"include_prerelease": true,
					},
				},
			},
		}

		opts := domain.BumpOptions{
			BumpKind: "patch",
		}

		bumped, err := domain.Bump(cfg, opts)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parser := types.Get("vsemver")
		rendered := parser.Render(bumped)
		// Expected: v1.2.3 is latest stable, bump patch -> v1.2.4 (NOT v1.2.5!)
		if rendered != "v1.2.4" {
			t.Errorf("expected v1.2.4, got %s", rendered)
		}
	})

	t.Run("flat provider config - include_prerelease true ignored on stable bump", func(t *testing.T) {
		cfg := &config.Config{
			VersionType: "vsemver",
			Provider: config.ProviderRaw{
				Type: "gittag",
				Raw: map[string]interface{}{
					"repo_dir":           dir,
					"include_prerelease": true,
				},
			},
		}

		opts := domain.BumpOptions{
			BumpKind: "patch",
		}

		bumped, err := domain.Bump(cfg, opts)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parser := types.Get("vsemver")
		rendered := parser.Render(bumped)
		// Expected: v1.2.3 is latest stable, bump patch -> v1.2.4 (NOT v1.2.5!)
		if rendered != "v1.2.4" {
			t.Errorf("expected v1.2.4, got %s", rendered)
		}
	})
}
