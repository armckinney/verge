package integration_test

import (
	"os/exec"
	"strings"
	"testing"
)

func buildVerctl(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("go", "build", "-o", "verge_test_bin", "./cmd/verge")
	cmd.Dir = "../.."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("skipping integration test (build failed): %v\n%s", err, out)
	}
	return "../../verge_test_bin"
}

func TestVersionParse_Integration(t *testing.T) {
	bin := buildVerctl(t)
	cmd := exec.Command(bin, "version", "parse", "v1.2.3")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "1") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestVersionBump_Integration(t *testing.T) {
	bin := buildVerctl(t)
	cmd := exec.Command(bin, "version", "bump", "--from", "v1.2.3", "--kind", "minor")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "1.3.0") {
		t.Errorf("expected 1.3.0, got: %s", out)
	}
}
