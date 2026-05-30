package domain

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"example.com/verge/internal/config"
)

func setupTestGitRepo(t *testing.T, tags []string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "verge-domain-test-*")
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

func TestNewFromConfig_GitTagNested(t *testing.T) {
	tags := []string{"v0.1.0", "v0.1.1-dev.1"}
	dir := setupTestGitRepo(t, tags)
	defer os.RemoveAll(dir)

	t.Run("nested provider config - include_prerelease true", func(t *testing.T) {
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

		p, err := NewFromConfig(cfg)
		if err != nil {
			t.Fatalf("expected no error building provider, got %v", err)
		}

		v, err := p.GetLatest("vsemver")
		if err != nil {
			t.Fatalf("expected no error getting latest, got %v", err)
		}
		if v.Original != "v0.1.1-dev.1" {
			t.Errorf("expected v0.1.1-dev.1, got %s", v.Original)
		}
	})

	t.Run("nested provider config - include_prerelease false", func(t *testing.T) {
		cfg := &config.Config{
			VersionType: "vsemver",
			Provider: config.ProviderRaw{
				Type: "gittag",
				Raw: map[string]interface{}{
					"gittag": map[string]interface{}{
						"repo_dir":           dir,
						"include_prerelease": false,
					},
				},
			},
		}

		p, err := NewFromConfig(cfg)
		if err != nil {
			t.Fatalf("expected no error building provider, got %v", err)
		}

		v, err := p.GetLatest("vsemver")
		if err != nil {
			t.Fatalf("expected no error getting latest, got %v", err)
		}
		if v.Original != "v0.1.0" {
			t.Errorf("expected v0.1.0, got %s", v.Original)
		}
	})

	t.Run("flat provider config - include_prerelease true", func(t *testing.T) {
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

		p, err := NewFromConfig(cfg)
		if err != nil {
			t.Fatalf("expected no error building provider, got %v", err)
		}

		v, err := p.GetLatest("vsemver")
		if err != nil {
			t.Fatalf("expected no error getting latest, got %v", err)
		}
		if v.Original != "v0.1.1-dev.1" {
			t.Errorf("expected v0.1.1-dev.1, got %s", v.Original)
		}
	})
}
