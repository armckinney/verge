package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("non-existent.yaml")
	if err != nil {
		t.Fatalf("expected no error for missing config, got %v", err)
	}
	if cfg.VersionType != "vsemver" {
		t.Errorf("expected vsemver default, got %s", cfg.VersionType)
	}
	if cfg.Provider.Type != "gittag" {
		t.Errorf("expected gittag default provider, got %s", cfg.Provider.Type)
	}
}

func TestLoadCustom(t *testing.T) {
	yamlContent := []byte(`
version_type: semver
default:
  bump_kind: minor
  prerelease_stage: rc
sequence:
  type: filehash
  targets: ["./Dockerfile"]
  length: 8
provider:
  type: ghcr
  ghcr:
    image: ghcr.io/owner/repo
`)
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(yamlContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.VersionType != "semver" {
		t.Errorf("expected semver, got %s", cfg.VersionType)
	}
	if cfg.Provider.Type != "ghcr" {
		t.Errorf("expected ghcr provider, got %s", cfg.Provider.Type)
	}
	
	// Test unmarshalling of nested raw map
	if ghcrRaw, ok := cfg.Provider.Raw["ghcr"].(map[string]interface{}); ok {
		if image, ok := ghcrRaw["image"].(string); ok {
			if image != "ghcr.io/owner/repo" {
				t.Errorf("expected image to be ghcr.io/owner/repo, got %s", image)
			}
		} else {
			t.Errorf("expected image to be of type string")
		}
	} else {
		t.Errorf("expected ghcr raw config to be available and typed")
	}
}
