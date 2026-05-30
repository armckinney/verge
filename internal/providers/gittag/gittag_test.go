package gittag

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupTestGitRepo(t *testing.T, tags []string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "verge-git-test-*")
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
	// Set dummy git config for commit
	runCmd("git", "config", "user.name", "Test User")
	runCmd("git", "config", "user.email", "test@example.com")

	// We need a commit to tag
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

func TestGetLatest(t *testing.T) {
	tags := []string{"v0.1.0", "v0.1.1-dev.1"}
	dir := setupTestGitRepo(t, tags)
	defer os.RemoveAll(dir)

	t.Run("include_prerelease true", func(t *testing.T) {
		p := NewProvider(Config{
			RepoDir:           dir,
			IncludePrerelease: true,
		})
		v, err := p.GetLatest("vsemver")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if v.Original != "v0.1.1-dev.1" {
			t.Errorf("expected v0.1.1-dev.1, got %s", v.Original)
		}
	})

	t.Run("include_prerelease false", func(t *testing.T) {
		p := NewProvider(Config{
			RepoDir:           dir,
			IncludePrerelease: false,
		})
		v, err := p.GetLatest("vsemver")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if v.Original != "v0.1.0" {
			t.Errorf("expected v0.1.0, got %s", v.Original)
		}
	})
}

func TestGetLatestSpecific(t *testing.T) {
	tags := []string{"v0.1.0", "v0.1.1-dev.1"}
	dir := setupTestGitRepo(t, tags)
	defer os.RemoveAll(dir)

	t.Run("matching pre-release with include_prerelease true", func(t *testing.T) {
		p := NewProvider(Config{
			RepoDir:           dir,
			IncludePrerelease: true,
		})
		v, err := p.GetLatestSpecific("vsemver", "v0.1.1-dev")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if v.Original != "v0.1.1-dev.1" {
			t.Errorf("expected v0.1.1-dev.1, got %s", v.Original)
		}
	})

	t.Run("matching pre-release with include_prerelease false", func(t *testing.T) {
		p := NewProvider(Config{
			RepoDir:           dir,
			IncludePrerelease: false,
		})
		_, err := p.GetLatestSpecific("vsemver", "v0.1.1-dev")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
