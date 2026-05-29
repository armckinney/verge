package version_test

import (
	"example.com/verge/internal/types"
	"testing"

	"example.com/verge/internal/version"
)

func TestCompare(t *testing.T) {
	parser := types.Get("semver")
	cmp := version.NewComparator()

	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.1.0", "1.2.0", -1},
		{"1.0.1", "1.0.2", -1},
		// Stage ordering: dev < alpha < beta < rc < final
		{"1.0.0-dev.1", "1.0.0-alpha.1", -1},
		{"1.0.0-alpha.1", "1.0.0-beta.1", -1},
		{"1.0.0-beta.1", "1.0.0-rc.1", -1},
		{"1.0.0-rc.1", "1.0.0", -1},
		{"1.0.0", "1.0.0-rc.1", 1},
		// Same stage, different sequence
		{"1.0.0-dev.1", "1.0.0-dev.2", -1},
		{"1.0.0-dev.2", "1.0.0-dev.1", 1},
		{"1.0.0-dev.1", "1.0.0-dev.1", 0},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			a, err := parser.Parse(tt.a)
			if err != nil {
				t.Fatalf("parsing %q: %v", tt.a, err)
			}
			b, err := parser.Parse(tt.b)
			if err != nil {
				t.Fatalf("parsing %q: %v", tt.b, err)
			}
			got := cmp.Compare(a, b)
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
