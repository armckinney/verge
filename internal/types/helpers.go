package types

import (
	"strings"

	"example.com/verge/internal/version"
)

func mustAtoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func detectSequence(s string) (interface{}, version.SequenceType) {
	if strings.HasPrefix(s, "gh-") {
		return s, version.SeqTypeBuildID
	}

	// Check if all digits
	allDigits := true
	for _, c := range s {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if allDigits && s != "" {
		n := mustAtoi(s)
		return n, version.SeqTypeNumeric
	}

	// Check if hex
	isHex := true
	for _, c := range strings.ToLower(s) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			isHex = false
			break
		}
	}
	if isHex && len(s) >= 32 {
		return strings.ToLower(s), version.SeqTypeContentHash
	}
	if isHex && len(s) >= 7 {
		return strings.ToLower(s), version.SeqTypeCommitSHA
	}

	return s, version.SeqTypeUnknown
}
