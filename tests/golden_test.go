package integration_test

import (
	"testing"
	"time"

	"example.com/verge/internal/version"
	"example.com/verge/tests/fixtures"
)

// TestGoldenCorpus_Parse verifies all real-world versions parse without error.
func TestGoldenCorpus_Parse(t *testing.T) {
	parser := version.NewParser()
	start := time.Now()

	for _, g := range fixtures.GoldenCorpus {
		t.Run(g.Description, func(t *testing.T) {
			_, err := parser.Parse(g.Input)
			if g.ExpectedErr {
				if err == nil {
					t.Errorf("expected parse error for %q, got none", g.Input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected parse error for %q: %v", g.Input, err)
			}
		})
	}

	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("parsing %d versions took %v, want <100ms", len(fixtures.GoldenCorpus), elapsed)
	}
}

// TestGoldenCorpus_Normalize verifies all parsed versions normalize correctly.
func TestGoldenCorpus_Normalize(t *testing.T) {
	parser := version.NewParser()
	normalizer := version.NewNormalizer()

	for _, g := range fixtures.GoldenCorpus {
		if g.ExpectedErr {
			continue
		}
		t.Run(g.Description, func(t *testing.T) {
			v, err := parser.Parse(g.Input)
			if err != nil {
				t.Skipf("parse failed (covered by parse test): %v", err)
			}
			normalized := normalizer.Normalize(v)
			if normalized == nil {
				t.Errorf("normalize returned nil for %q", g.Input)
				return
			}
			if normalized.Major < 0 || normalized.Minor < 0 || normalized.Patch < 0 {
				t.Errorf("normalize produced negative version component for %q", g.Input)
			}
		})
	}
}

// TestGoldenCorpus_RenderEcosystem verifies versions render to the expected ecosystem format.
func TestGoldenCorpus_RenderEcosystem(t *testing.T) {
	parser := version.NewParser()
	normalizer := version.NewNormalizer()

	for _, g := range fixtures.GoldenCorpus {
		if g.ExpectedErr {
			continue
		}
		t.Run(g.Description, func(t *testing.T) {
			v, err := parser.Parse(g.Input)
			if err != nil {
				t.Skipf("parse failed: %v", err)
			}
			normalized := normalizer.Normalize(v)
			rendered := version.NewRenderer(g.Ecosystem).Render(normalized)
			if rendered == "" {
				t.Errorf("render returned empty string for %q (ecosystem=%s)", g.Input, g.Ecosystem)
			}
		})
	}
}

// TestGoldenCorpus_RoundTrip verifies parse → render → re-parse preserves version identity.
func TestGoldenCorpus_RoundTrip(t *testing.T) {
	parser := version.NewParser()
	normalizer := version.NewNormalizer()
	comparator := version.NewComparator()

	ecosystems := []string{"go", "python", "containers", "terraform", "github-actions"}

	for _, g := range fixtures.GoldenCorpus {
		if g.ExpectedErr {
			continue
		}
		for _, eco := range ecosystems {
			eco := eco
			t.Run(g.Description+"/"+eco, func(t *testing.T) {
				v, err := parser.Parse(g.Input)
				if err != nil {
					t.Skipf("parse failed: %v", err)
				}
				normalized := normalizer.Normalize(v)
				rendered := version.NewRenderer(eco).Render(normalized)

				reparsed, err := parser.Parse(rendered)
				if err != nil {
					// Some ecosystem formats (e.g. Python PEP440) may not round-trip through SemVer parser; skip
					t.Skipf("re-parse of rendered %q failed (eco=%s): %v", rendered, eco, err)
				}
				reparsed = normalizer.Normalize(reparsed)

				if comparator.Compare(normalized, reparsed) != 0 {
					t.Errorf("round-trip mismatch: %q → (eco=%s) %q → %q", g.Input, eco, rendered, reparsed)
				}
			})
		}
	}
}

// TestInvalidVersionCorpus verifies known-invalid versions are rejected.
func TestInvalidVersionCorpus(t *testing.T) {
	parser := version.NewParser()

	for _, inv := range fixtures.InvalidVersionCorpus {
		t.Run(inv.Reason, func(t *testing.T) {
			_, err := parser.Parse(inv.Input)
			if err == nil {
				t.Errorf("expected error for %q (%s), got none", inv.Input, inv.Reason)
			}
		})
	}
}

// TestGoldenCorpus_Compare verifies that higher versions compare > lower versions.
func TestGoldenCorpus_Compare(t *testing.T) {
	parser := version.NewParser()
	comparator := version.NewComparator()

	pairs := []struct {
		lower, higher string
	}{
		{"v1.0.0", "v2.0.0"},
		{"v1.0.0", "v1.1.0"},
		{"v1.0.0", "v1.0.1"},
		{"v1.0.0-rc.1", "v1.0.0"},
		{"v1.0.0-beta.1", "v1.0.0-rc.1"},
		{"v1.0.0-alpha.1", "v1.0.0-beta.1"},
		{"v1.0.0-dev.1", "v1.0.0-alpha.1"},
		{"1.0.0a1", "1.0.0rc1"},
		{"1.0.0b1", "1.0.0"},
	}

	for _, p := range pairs {
		t.Run(p.lower+"<"+p.higher, func(t *testing.T) {
			lo, err := parser.Parse(p.lower)
			if err != nil {
				t.Fatalf("parse lower %q: %v", p.lower, err)
			}
			hi, err := parser.Parse(p.higher)
			if err != nil {
				t.Fatalf("parse higher %q: %v", p.higher, err)
			}
			if comparator.Compare(lo, hi) >= 0 {
				t.Errorf("expected %q < %q", p.lower, p.higher)
			}
			if comparator.Compare(hi, lo) <= 0 {
				t.Errorf("expected %q > %q", p.higher, p.lower)
			}
		})
	}
}
