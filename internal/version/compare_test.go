package version_test

import (
	"testing"

	"example.com/verge/internal/version"
)

func TestCompare(t *testing.T) {
	parser := version.NewParser()
	cmp := version.NewComparator()

	tests := []struct {
		a, b string
		want int
	}{
		{"v1.0.0", "v1.0.0", 0},
		{"v1.0.0", "v2.0.0", -1},
		{"v2.0.0", "v1.0.0", 1},
		{"v1.1.0", "v1.2.0", -1},
		{"v1.0.1", "v1.0.2", -1},
		// Stage ordering: dev < alpha < beta < rc < final
		{"v1.0.0-dev.1", "v1.0.0-alpha.1", -1},
		{"v1.0.0-alpha.1", "v1.0.0-beta.1", -1},
		{"v1.0.0-beta.1", "v1.0.0-rc.1", -1},
		{"v1.0.0-rc.1", "v1.0.0", -1},
		{"v1.0.0", "v1.0.0-rc.1", 1},
		// Same stage, different sequence
		{"v1.0.0-dev.1", "v1.0.0-dev.2", -1},
		{"v1.0.0-dev.2", "v1.0.0-dev.1", 1},
		{"v1.0.0-dev.1", "v1.0.0-dev.1", 0},
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
