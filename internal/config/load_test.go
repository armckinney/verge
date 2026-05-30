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

func TestParseOverrides(t *testing.T) {
	pairs := []string{
		"bool_true=true",
		"bool_false=false",
		"int_val=42",
		"str_val=hello",
		"invalid-pair",
	}

	out := ParseOverrides(pairs)

	if out["bool_true"] != true {
		t.Errorf("expected bool_true to be true, got %v", out["bool_true"])
	}
	if out["bool_false"] != false {
		t.Errorf("expected bool_false to be false, got %v", out["bool_false"])
	}
	if out["int_val"] != 42 {
		t.Errorf("expected int_val to be 42, got %v", out["int_val"])
	}
	if out["str_val"] != "hello" {
		t.Errorf("expected str_val to be 'hello', got %v", out["str_val"])
	}
	if _, ok := out["invalid-pair"]; ok {
		t.Errorf("expected invalid-pair to be ignored")
	}
}

func TestMergeOverrides(t *testing.T) {
	t.Run("Flat Merge", func(t *testing.T) {
		raw := &ProviderRaw{
			Type: "gittag",
			Raw: map[string]interface{}{
				"include_prerelease": true,
				"other_field":        "value",
			},
		}

		overrides := map[string]interface{}{
			"include_prerelease": false,
			"new_field":          123,
		}

		MergeOverrides(raw, overrides)

		if raw.Raw["include_prerelease"] != false {
			t.Errorf("expected include_prerelease to be false, got %v", raw.Raw["include_prerelease"])
		}
		if raw.Raw["other_field"] != "value" {
			t.Errorf("expected other_field to remain 'value', got %v", raw.Raw["other_field"])
		}
		if raw.Raw["new_field"] != 123 {
			t.Errorf("expected new_field to be 123, got %v", raw.Raw["new_field"])
		}
	})

	t.Run("Nested Merge", func(t *testing.T) {
		raw := &ProviderRaw{
			Type: "gittag",
			Raw: map[string]interface{}{
				"gittag": map[string]interface{}{
					"include_prerelease": true,
					"other_field":        "value",
				},
			},
		}

		overrides := map[string]interface{}{
			"include_prerelease": false,
			"new_field":          123,
		}

		MergeOverrides(raw, overrides)

		nested, ok := raw.Raw["gittag"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected raw.Raw[\"gittag\"] to be map[string]interface{}")
		}

		if nested["include_prerelease"] != false {
			t.Errorf("expected nested include_prerelease to be false, got %v", nested["include_prerelease"])
		}
		if nested["other_field"] != "value" {
			t.Errorf("expected nested other_field to remain 'value', got %v", nested["other_field"])
		}
		if nested["new_field"] != 123 {
			t.Errorf("expected nested new_field to be 123, got %v", nested["new_field"])
		}
	})
}
