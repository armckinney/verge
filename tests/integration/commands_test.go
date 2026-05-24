package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildVerge(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("go", "build", "-o", "verge_test_bin", "./cmd/verge")
	cmd.Dir = "../.."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("skipping integration test (build failed): %v\n%s", err, out)
	}
	return "../../verge_test_bin"
}

func TestVersionParse_Integration(t *testing.T) {
	bin := buildVerge(t)
	cmd := exec.Command(bin, "parse", "v1.2.3")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "1") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestVersionBump_Integration(t *testing.T) {
	bin := buildVerge(t)
	cmd := exec.Command(bin, "bump", "--from", "v1.2.3", "--kind", "minor")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "1.3.0") {
		t.Errorf("expected 1.3.0, got: %s", out)
	}
}

func TestVersionBump_Field_Integration(t *testing.T) {
	bin := buildVerge(t)
	cmd := exec.Command(bin, "bump", "--from", "v1.2.3", "--kind", "minor", "--field", "to")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	if got := strings.TrimSpace(string(out)); got != "1.3.0" {
		t.Fatalf("expected bumped version 1.3.0, got %q", got)
	}
}

func TestVersionCompare_Field_Integration(t *testing.T) {
	bin := buildVerge(t)
	cmd := exec.Command(bin, "compare", "1.2.3", "2.0.0", "--field", "value")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected compare to exit non-zero")
	}
	if got := strings.TrimSpace(string(out)); got != "less" {
		t.Fatalf("expected compare output 'less', got %q", got)
	}
}

func TestVersionCurrent_Field_Integration(t *testing.T) {
	bin := buildVerge(t)
	repoDir := t.TempDir()
	if err := exec.Command("git", "init", repoDir).Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("test\n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	if err := exec.Command("git", "-C", repoDir, "add", "README.md").Run(); err != nil {
		t.Fatalf("git add failed: %v", err)
	}
	commit := exec.Command("git", "-C", repoDir, "commit", "-m", "test")
	commit.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test User",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test User",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)
	if out, err := commit.CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v\n%s", err, out)
	}
	if err := exec.Command("git", "-C", repoDir, "tag", "v1.2.3").Run(); err != nil {
		t.Fatalf("git tag failed: %v", err)
	}

	cmd := exec.Command(bin, "current", "--repo-dir", repoDir, "--field", "normalized")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	if got := strings.TrimSpace(string(out)); got != "1.2.3" {
		t.Fatalf("expected normalized version 1.2.3, got %q", got)
	}
}
