package version

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	semverRe = regexp.MustCompile(`^(v)?(\d+)\.(\d+)\.(\d+)(?:-(dev|alpha|beta|rc|a|b)\.([a-zA-Z0-9_-]+))?(?:\+[a-zA-Z0-9._-]+)?$`)
	pep440Re = regexp.MustCompile(`^(v)?(\d+)\.(\d+)\.(\d+)(dev|alpha|beta|rc|a|b)(\d+)$`)
)

type defaultParser struct{}

// NewParser returns a new Parser.
func NewParser() Parser {
	return &defaultParser{}
}

func (p *defaultParser) Parse(input string) (*Version, error) {
	if input == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Try SemVer first
	if m := semverRe.FindStringSubmatch(input); m != nil {
		return parseSemVer(input, m)
	}

	// Try PEP440
	if m := pep440Re.FindStringSubmatch(input); m != nil {
		return parsePEP440(input, m)
	}

	return nil, fmt.Errorf("unable to parse version: %q", input)
}

func parseSemVer(input string, m []string) (*Version, error) {
	major := mustAtoi(m[2])
	minor := mustAtoi(m[3])
	patch := mustAtoi(m[4])

	v := &Version{
		Major:    major,
		Minor:    minor,
		Patch:    patch,
		Stage:    StageFinal,
		Original: input,
		Scheme:   SchemeSemVer,
	}

	if m[5] != "" {
		stage, err := stageFromAbbrev(m[5])
		if err != nil {
			return nil, err
		}
		v.Stage = stage
		seq, seqType := detectSequence(m[6])
		v.Sequence = seq
		v.SequenceType = seqType
	}

	return v, nil
}

func parsePEP440(input string, m []string) (*Version, error) {
	major := mustAtoi(m[2])
	minor := mustAtoi(m[3])
	patch := mustAtoi(m[4])

	stage, err := stageFromAbbrev(m[5])
	if err != nil {
		return nil, err
	}

	seqStr := m[6]
	seq, seqType := detectSequence(seqStr)

	return &Version{
		Major:        major,
		Minor:        minor,
		Patch:        patch,
		Stage:        stage,
		Sequence:     seq,
		SequenceType: seqType,
		Original:     input,
		Scheme:       SchemePEP440,
	}, nil
}

func stageFromAbbrev(s string) (Stage, error) {
	switch strings.ToLower(s) {
	case "dev":
		return StageDev, nil
	case "alpha", "a":
		return StageAlpha, nil
	case "beta", "b":
		return StageBeta, nil
	case "rc":
		return StageRC, nil
	}
	return StageFinal, fmt.Errorf("unknown stage abbreviation: %q", s)
}

func detectSequence(s string) (interface{}, SequenceType) {
	if strings.HasPrefix(s, "gh-") {
		return s, SeqTypeBuildID
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
		return n, SeqTypeNumeric
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
		return strings.ToLower(s), SeqTypeContentHash
	}
	if isHex && len(s) >= 7 {
		return strings.ToLower(s), SeqTypeCommitSHA
	}

	return s, SeqTypeUnknown
}

func mustAtoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}
