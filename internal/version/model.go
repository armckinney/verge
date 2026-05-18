package version

import (
	"fmt"
	"strings"
)

type Stage int

const (
	StageFinal Stage = iota
	StageDev
	StageAlpha
	StageBeta
	StageRC
)

func (s Stage) String() string {
	switch s {
	case StageFinal:
		return "final"
	case StageDev:
		return "dev"
	case StageAlpha:
		return "alpha"
	case StageBeta:
		return "beta"
	case StageRC:
		return "rc"
	}
	return "unknown"
}

func StageFromString(s string) (Stage, error) {
	switch strings.ToLower(s) {
	case "final", "release", "":
		return StageFinal, nil
	case "dev":
		return StageDev, nil
	case "alpha", "a":
		return StageAlpha, nil
	case "beta", "b":
		return StageBeta, nil
	case "rc":
		return StageRC, nil
	}
	return StageFinal, fmt.Errorf("unknown stage: %q", s)
}

type Scheme string

const (
	SchemeSemVer Scheme = "semver"
	SchemePEP440 Scheme = "pep440"
	SchemeAuto   Scheme = "auto"
)

type SequenceType string

const (
	SeqTypeNumeric     SequenceType = "numeric"
	SeqTypeCommitSHA   SequenceType = "commit-sha"
	SeqTypeContentHash SequenceType = "content-hash"
	SeqTypeBuildID     SequenceType = "build-id"
	SeqTypeUnknown     SequenceType = "unknown"
)

type Version struct {
	Major        int
	Minor        int
	Patch        int
	Stage        Stage
	Sequence     interface{} // int or string
	SequenceType SequenceType
	Original     string
	Scheme       Scheme
}

func (v *Version) String() string {
	if v.Stage == StageFinal {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s.%v", v.Major, v.Minor, v.Patch, v.Stage, v.Sequence)
}

func (v *Version) IsPrerelease() bool {
	return v.Stage != StageFinal
}

func (v *Version) Core() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) Validate() error {
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		return fmt.Errorf("negative version components not allowed")
	}
	if v.Stage == StageFinal && v.Sequence != nil {
		return fmt.Errorf("final release must not have a sequence")
	}
	if v.Stage != StageFinal && v.Sequence == nil {
		return fmt.Errorf("prerelease stage %s requires a sequence", v.Stage)
	}
	return nil
}
