package cli

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupTestGitRepo(t *testing.T, tags []string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "verge-cli-test-*")
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

func TestCLIProviderConfigOverride(t *testing.T) {
	tags := []string{"v1.2.3", "v1.2.4-dev.5"}
	repoDir := setupTestGitRepo(t, tags)
	defer os.RemoveAll(repoDir)

	// Create temp verge config with include_prerelease=false
	configContent := `
version_type: vsemver
provider:
  type: gittag
  gittag:
    repo_dir: "` + repoDir + `"
    include_prerelease: false
`
	tmpConfig, err := os.CreateTemp("", "verge-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpConfig.Name())
	if _, err := tmpConfig.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpConfig.Close()

	t.Run("latest command defaults to false from config", func(t *testing.T) {
		// Reset globalFlags
		globalFlags.field = ""
		globalFlags.format = "text"
		globalFlags.json = false
		globalFlags.configPath = tmpConfig.Name()

		root := &cobra.Command{Use: "verge"}
		root.PersistentFlags().StringVarP(&globalFlags.configPath, "config", "c", "", "")
		root.PersistentFlags().StringVarP(&globalFlags.format, "format", "f", "text", "")
		root.PersistentFlags().BoolVar(&globalFlags.json, "json", false, "")
		root.PersistentFlags().StringVar(&globalFlags.field, "field", "", "")

		cmd := versionLatestCmd()
		root.AddCommand(cmd)

		root.SetArgs([]string{"latest", "--config", tmpConfig.Name(), "--format", "text"})

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := root.Execute()
		w.Close()

		var outBuf bytes.Buffer
		_, _ = io.Copy(&outBuf, r)
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got := outBuf.String()
		// Expected to resolve v1.2.3 because include_prerelease is false in config
		if !strings.Contains(got, "v1.2.3") {
			t.Errorf("expected output to contain v1.2.3, got %q", got)
		}
	})

	t.Run("latest command overrides include_prerelease=true via --provider-config", func(t *testing.T) {
		// Reset globalFlags
		globalFlags.field = ""
		globalFlags.format = "text"
		globalFlags.json = false
		globalFlags.configPath = tmpConfig.Name()

		root := &cobra.Command{Use: "verge"}
		root.PersistentFlags().StringVarP(&globalFlags.configPath, "config", "c", "", "")
		root.PersistentFlags().StringVarP(&globalFlags.format, "format", "f", "text", "")
		root.PersistentFlags().BoolVar(&globalFlags.json, "json", false, "")
		root.PersistentFlags().StringVar(&globalFlags.field, "field", "", "")

		cmd := versionLatestCmd()
		root.AddCommand(cmd)

		root.SetArgs([]string{"latest", "--config", tmpConfig.Name(), "--format", "text", "--provider-config", "include_prerelease=true"})

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := root.Execute()
		w.Close()

		var outBuf bytes.Buffer
		_, _ = io.Copy(&outBuf, r)
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got := outBuf.String()
		// Expected to resolve v1.2.4-dev.5 because of --provider-config override
		if !strings.Contains(got, "v1.2.4-dev.5") {
			t.Errorf("expected output to contain v1.2.4-dev.5, got %q", got)
		}
	})

	t.Run("bump command overrides include_prerelease=true via --provider-config", func(t *testing.T) {
		// Reset globalFlags
		globalFlags.field = ""
		globalFlags.format = "text"
		globalFlags.json = false
		globalFlags.configPath = tmpConfig.Name()

		root := &cobra.Command{Use: "verge"}
		root.PersistentFlags().StringVarP(&globalFlags.configPath, "config", "c", "", "")
		root.PersistentFlags().StringVarP(&globalFlags.format, "format", "f", "text", "")
		root.PersistentFlags().BoolVar(&globalFlags.json, "json", false, "")
		root.PersistentFlags().StringVar(&globalFlags.field, "field", "", "")

		cmd := versionBumpCmd()
		root.AddCommand(cmd)

		// Bump prerelease on top of v1.2.4-dev.5
		root.SetArgs([]string{"bump", "--config", tmpConfig.Name(), "--kind", "prerelease", "--stage", "dev", "--format", "text", "--provider-config", "include_prerelease=true"})

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := root.Execute()
		w.Close()

		var outBuf bytes.Buffer
		_, _ = io.Copy(&outBuf, r)
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got := outBuf.String()
		// Expected to resolve v1.2.4-dev.5 and bump to v1.2.4-dev.6
		if !strings.Contains(got, "v1.2.4-dev.6") {
			t.Errorf("expected output to contain v1.2.4-dev.6, got %q", got)
		}
	})
}
