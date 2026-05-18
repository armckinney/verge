package sequence

import (
	"strings"

	"example.com/verge/internal/version"
)

// Detect detects the sequence type for a given string.
func Detect(s string) version.SequenceType {
	if strings.HasPrefix(s, "gh-") {
		return version.SeqTypeBuildID
	}

	allDigits := true
	for _, c := range s {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if allDigits && s != "" {
		return version.SeqTypeNumeric
	}

	isHex := true
	for _, c := range strings.ToLower(s) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			isHex = false
			break
		}
	}
	if isHex && len(s) >= 32 {
		return version.SeqTypeContentHash
	}
	if isHex && len(s) >= 7 {
		return version.SeqTypeCommitSHA
	}

	return version.SeqTypeUnknown
}
