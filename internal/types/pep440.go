package types

import (
	"fmt"
	"regexp"

	"example.com/verge/internal/version"
)

var pep440Re = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:(dev|a|alpha|b|beta|rc)(\d+))?$`)

// PEP440 parses and renders Python PEP440 format versions.
type PEP440 struct{}

func (s *PEP440) Name() string { return "pep440" }

func (s *PEP440) Parse(input string) (*version.Version, error) {
	m := pep440Re.FindStringSubmatch(input)
	if m == nil {
		return nil, fmt.Errorf("invalid pep440: %q", input)
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

func (s *PEP440) Render(v *version.Version) string {
	if v.Stage == version.StageFinal {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}

	// pep440 format: a/b/rc
	stageStr := ""
	switch v.Stage {
	case version.StageDev:
		stageStr = "dev"
	case version.StageA:
		stageStr = "a"
	case version.StageB:
		stageStr = "b"
	case version.StageRC:
		stageStr = "rc"
	}

	return fmt.Sprintf("%d.%d.%d%s%v", v.Major, v.Minor, v.Patch, stageStr, v.Sequence)
}
