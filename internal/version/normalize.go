package version

import "strings"

type defaultNormalizer struct{}

// NewNormalizer returns a new Normalizer.
func NewNormalizer() Normalizer {
	return &defaultNormalizer{}
}

func (n *defaultNormalizer) Normalize(v *Version) *Version {
	result := *v // shallow copy

	if v.Sequence != nil {
		switch seq := v.Sequence.(type) {
		case string:
			// Lowercase hex for commit SHAs and content hashes
			if v.SequenceType == SeqTypeCommitSHA || v.SequenceType == SeqTypeContentHash {
				result.Sequence = strings.ToLower(seq)
			}
		}
	}

	return &result
}
