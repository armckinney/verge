package version_test

import (
	"testing"

	"example.com/template-go/internal/version"
)

func TestRender(t *testing.T) {
	tests := []struct {
		ecosystem string
		v         *version.Version
		expected  string
	}{
		{
			"go",
			&version.Version{Major: 1, Minor: 2, Patch: 3, Stage: version.StageFinal},
			"v1.2.3",
		},
		{
			"go",
			&version.Version{Major: 1, Minor: 2, Patch: 3, Stage: version.StageDev, Sequence: 1, SequenceType: version.SeqTypeNumeric},
			"v1.2.3-dev.1",
		},
		{
			"python",
			&version.Version{Major: 1, Minor: 2, Patch: 3, Stage: version.StageFinal},
			"1.2.3",
		},
		{
			"python",
			&version.Version{Major: 1, Minor: 2, Patch: 3, Stage: version.StageAlpha, Sequence: 1, SequenceType: version.SeqTypeNumeric},
			"1.2.3a1",
		},
		{
			"terraform",
			&version.Version{Major: 1, Minor: 2, Patch: 3, Stage: version.StageFinal},
			"v1.2.3",
		},
		{
			"container",
			&version.Version{Major: 1, Minor: 2, Patch: 3, Stage: version.StageFinal},
			"1.2.3",
		},
		{
			"github-actions",
			&version.Version{Major: 1, Minor: 2, Patch: 3, Stage: version.StageFinal},
			"1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.ecosystem+"/"+tt.v.String(), func(t *testing.T) {
			renderer := version.NewRenderer(tt.ecosystem)
			got := renderer.Render(tt.v)
			if got != tt.expected {
				t.Errorf("Render(%v) = %q, want %q", tt.v, got, tt.expected)
			}
		})
	}
}
