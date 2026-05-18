package version_test

import (
	"testing"

	"example.com/template-go/internal/version"
)

func TestNormalize(t *testing.T) {
	normalizer := version.NewNormalizer()

	t.Run("lowercase commit sha", func(t *testing.T) {
		v := &version.Version{
			Major:        1,
			Minor:        2,
			Patch:        3,
			Stage:        version.StageDev,
			Sequence:     "ABC1234",
			SequenceType: version.SeqTypeCommitSHA,
		}
		normalized := normalizer.Normalize(v)
		if seq, ok := normalized.Sequence.(string); !ok || seq != "abc1234" {
			t.Errorf("expected lowercase commit SHA, got %v", normalized.Sequence)
		}
	})

	t.Run("passthrough numeric", func(t *testing.T) {
		v := &version.Version{
			Major:        1,
			Sequence:     42,
			SequenceType: version.SeqTypeNumeric,
		}
		normalized := normalizer.Normalize(v)
		if seq, ok := normalized.Sequence.(int); !ok || seq != 42 {
			t.Errorf("expected numeric 42, got %v", normalized.Sequence)
		}
	})
}
