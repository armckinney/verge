package version_test

import (
	"testing"

	"example.com/verge/internal/version"
)

func TestParse_SemVer(t *testing.T) {
	parser := version.NewParser()

	tests := []struct {
		input   string
		wantErr bool
		major   int
		minor   int
		patch   int
		stage   version.Stage
		seqType version.SequenceType
	}{
		{"v1.2.3", false, 1, 2, 3, version.StageFinal, ""},
		{"1.2.3", false, 1, 2, 3, version.StageFinal, ""},
		{"v0.0.1", false, 0, 0, 1, version.StageFinal, ""},
		{"v1.2.3-dev.1", false, 1, 2, 3, version.StageDev, version.SeqTypeNumeric},
		{"v1.2.3-alpha.2", false, 1, 2, 3, version.StageAlpha, version.SeqTypeNumeric},
		{"v1.2.3-beta.3", false, 1, 2, 3, version.StageBeta, version.SeqTypeNumeric},
		{"v1.2.3-rc.1", false, 1, 2, 3, version.StageRC, version.SeqTypeNumeric},
		{"v1.2.3-a.4", false, 1, 2, 3, version.StageAlpha, version.SeqTypeNumeric},
		{"v1.2.3-b.5", false, 1, 2, 3, version.StageBeta, version.SeqTypeNumeric},
		{"v1.2.3+build.1", false, 1, 2, 3, version.StageFinal, ""},
		{"not-a-version", true, 0, 0, 0, version.StageFinal, ""},
		{"", true, 0, 0, 0, version.StageFinal, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v, err := parser.Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.Major != tt.major || v.Minor != tt.minor || v.Patch != tt.patch {
				t.Errorf("got %d.%d.%d, want %d.%d.%d", v.Major, v.Minor, v.Patch, tt.major, tt.minor, tt.patch)
			}
			if v.Stage != tt.stage {
				t.Errorf("got stage %v, want %v", v.Stage, tt.stage)
			}
			if tt.seqType != "" && v.SequenceType != tt.seqType {
				t.Errorf("got seqType %v, want %v", v.SequenceType, tt.seqType)
			}
		})
	}
}

func TestParse_PEP440(t *testing.T) {
	parser := version.NewParser()

	tests := []struct {
		input   string
		wantErr bool
		stage   version.Stage
		seqVal  int
	}{
		{"1.2.3a1", false, version.StageAlpha, 1},
		{"1.2.3b2", false, version.StageBeta, 2},
		{"1.2.3rc3", false, version.StageRC, 3},
		{"1.2.3dev1", false, version.StageDev, 1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v, err := parser.Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.Stage != tt.stage {
				t.Errorf("got stage %v, want %v", v.Stage, tt.stage)
			}
			if seq, ok := v.Sequence.(int); !ok || seq != tt.seqVal {
				t.Errorf("got sequence %v, want %d", v.Sequence, tt.seqVal)
			}
			if v.Scheme != version.SchemePEP440 {
				t.Errorf("expected PEP440 scheme, got %v", v.Scheme)
			}
		})
	}
}

func TestParse_CommitSHA(t *testing.T) {
	parser := version.NewParser()
	v, err := parser.Parse("v1.2.3-dev.abc1234")
	if err != nil {
		t.Fatal(err)
	}
	if v.SequenceType != version.SeqTypeCommitSHA {
		t.Errorf("expected commit-sha, got %v", v.SequenceType)
	}
}

func TestParse_ContentHash(t *testing.T) {
	parser := version.NewParser()
	v, err := parser.Parse("v1.2.3-dev.abc1234def5678901234567890123456")
	if err != nil {
		t.Fatal(err)
	}
	if v.SequenceType != version.SeqTypeContentHash {
		t.Errorf("expected content-hash, got %v", v.SequenceType)
	}
}

func TestParse_BuildID(t *testing.T) {
	parser := version.NewParser()
	v, err := parser.Parse("v1.2.3-dev.gh-12345")
	if err != nil {
		t.Fatal(err)
	}
	if v.SequenceType != version.SeqTypeBuildID {
		t.Errorf("expected build-id, got %v", v.SequenceType)
	}
}
