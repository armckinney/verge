package version

import (
	"fmt"
	"strings"
)

type Stage int

const (
	StageFinal Stage = iota
	StageDev
	StageA
	StageB
	StageRC
)

func (s Stage) String() string {
	switch s {
	case StageFinal:
		return "final"
	case StageDev:
		return "dev"
	case StageA:
		return "a"
	case StageB:
		return "b"
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
		return StageA, nil
	case "beta", "b":
		return StageB, nil
	case "rc":
		return StageRC, nil
	}
	return StageFinal, fmt.Errorf("unknown stage: %q", s)
}

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
	VersionType  string // e.g., "semver", "vsemver", "pep440"
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
