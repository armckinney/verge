package version_test

import (
	"testing"

	"example.com/verge/internal/types"
	"example.com/verge/internal/version"
)

func TestBump(t *testing.T) {
	parser := types.Get("vsemver")
	bumper := version.NewBumper()

	tests := []struct {
		input    string
		kind     version.BumpKind
		stage    version.Stage
		expected string
	}{
		{"v1.2.3", version.BumpMajor, version.StageFinal, "2.0.0"},
		{"v1.2.3", version.BumpMinor, version.StageFinal, "1.3.0"},
		{"v1.2.3", version.BumpPatch, version.StageFinal, "1.2.4"},
		// Final to prerelease: bump patch, add stage.nil (sequence calculator will set it)
		{"v1.2.3", version.BumpPrerelease, version.StageDev, "1.2.4-dev.<nil>"},
		{"v1.2.3", version.BumpPrerelease, version.StageA, "1.2.4-a.<nil>"},
		// Same stage: keep sequence as-is (sequence calculator will increment it)
		{"v1.2.3-dev.1", version.BumpPrerelease, version.StageDev, "1.2.3-dev.1"},
		{"v1.2.3-alpha.3", version.BumpPrerelease, version.StageA, "1.2.3-a.3"},
		// Different stage: keep core, change stage, seq=nil
		{"v1.2.3-dev.1", version.BumpPrerelease, version.StageA, "1.2.3-a.<nil>"},
		{"v1.2.3-alpha.2", version.BumpPrerelease, version.StageB, "1.2.3-b.<nil>"},
		// Final: remove prerelease
		{"v1.2.3-rc.2", version.BumpFinal, version.StageFinal, "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.input+" "+string(tt.kind), func(t *testing.T) {
			v, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("parsing %q: %v", tt.input, err)
			}
			bumped, err := bumper.Bump(v, tt.kind, tt.stage)
			if err != nil {
				t.Fatalf("bumping: %v", err)
			}
			if bumped.String() != tt.expected {
				t.Errorf("Bump(%q, %v, %v) = %q, want %q", tt.input, tt.kind, tt.stage, bumped.String(), tt.expected)
			}
		})
	}
}
