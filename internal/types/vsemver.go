package types

import (
	"fmt"
	"regexp"

	"example.com/verge/internal/version"
)

var vsemverRe = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)(?:-(dev|alpha|beta|rc|a|b)\.([a-zA-Z0-9_-]+))?(?:\+[a-zA-Z0-9._-]+)?$`)

// VSemVer parses and renders Semantic Versioning with a generic 'v' prefix.
type VSemVer struct{}

func (s *VSemVer) Name() string { return "vsemver" }

func (s *VSemVer) Parse(input string) (*version.Version, error) {
	m := vsemverRe.FindStringSubmatch(input)
	if m == nil {
		return nil, fmt.Errorf("invalid vsemver: %q", input)
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

func (s *VSemVer) Render(v *version.Version) string {
	if v.Stage == version.StageFinal {
		return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	}

	// vsemver format: alpha/beta/rc
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

	return fmt.Sprintf("v%d.%d.%d-%s.%v", v.Major, v.Minor, v.Patch, stageStr, v.Sequence)
}
