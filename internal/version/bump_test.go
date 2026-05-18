package version_test

import (
	"testing"

	"example.com/template-go/internal/version"
)

func TestBump(t *testing.T) {
	parser := version.NewParser()
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
		// Final to prerelease: bump patch, add stage.1
		{"v1.2.3", version.BumpPrerelease, version.StageDev, "1.2.4-dev.1"},
		{"v1.2.3", version.BumpPrerelease, version.StageAlpha, "1.2.4-alpha.1"},
		// Same stage: increment sequence
		{"v1.2.3-dev.1", version.BumpPrerelease, version.StageDev, "1.2.3-dev.2"},
		{"v1.2.3-alpha.3", version.BumpPrerelease, version.StageAlpha, "1.2.3-alpha.4"},
		// Different stage: keep core, change stage, seq=1
		{"v1.2.3-dev.1", version.BumpPrerelease, version.StageAlpha, "1.2.3-alpha.1"},
		{"v1.2.3-alpha.2", version.BumpPrerelease, version.StageBeta, "1.2.3-beta.1"},
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
