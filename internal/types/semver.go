package types

import (
	"fmt"
	"regexp"

	"example.com/verge/internal/version"
)

var semverRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-(dev|alpha|beta|rc|a|b)\.([a-zA-Z0-9_-]+))?(?:\+[a-zA-Z0-9._-]+)?$`)

// SemVer parses and renders standard Semantic Versioning.
type SemVer struct{}

func (s *SemVer) Name() string { return "semver" }

func (s *SemVer) Parse(input string) (*version.Version, error) {
	m := semverRe.FindStringSubmatch(input)
	if m == nil {
		return nil, fmt.Errorf("invalid semver: %q", input)
	}

	major := mustAtoi(m[1])
	minor := mustAtoi(m[2])
	patch := mustAtoi(m[3])

	v := &version.Version{
		Major:       major,
		Minor:       minor,
		Patch:       patch,
		Stage:       version.StageFinal,
		Original:    input,
		VersionType: s.Name(),
	}

	if m[4] != "" {
		stage, err := version.StageFromString(m[4])
		if err != nil {
			return nil, err
		}
		v.Stage = stage
		seq, seqType := detectSequence(m[5])
		v.Sequence = seq
		v.SequenceType = seqType
	}

	return v, nil
}

func (s *SemVer) Render(v *version.Version) string {
	if v.Stage == version.StageFinal {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}

	// semver format: alpha/beta/rc
	stageStr := ""
	switch v.Stage {
	case version.StageDev:
		stageStr = "dev"
	case version.StageA:
		stageStr = "alpha"
	case version.StageB:
		stageStr = "beta"
	case version.StageRC:
		stageStr = "rc"
	}

	return fmt.Sprintf("%d.%d.%d-%s.%v", v.Major, v.Minor, v.Patch, stageStr, v.Sequence)
}
