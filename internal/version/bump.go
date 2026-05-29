package version

import "fmt"

type defaultBumper struct{}

// NewBumper returns a new Bumper.
func NewBumper() Bumper {
	return &defaultBumper{}
}

func (b *defaultBumper) Bump(v *Version, kind BumpKind, stage Stage) (*Version, error) {
	result := *v
	result.Original = ""

	switch kind {
	case BumpMajor:
		result.Major++
		result.Minor = 0
		result.Patch = 0
		result.Stage = StageFinal
		result.Sequence = nil
		result.SequenceType = ""
	case BumpMinor:
		result.Minor++
		result.Patch = 0
		result.Stage = StageFinal
		result.Sequence = nil
		result.SequenceType = ""
	case BumpPatch:
		result.Patch++
		result.Stage = StageFinal
		result.Sequence = nil
		result.SequenceType = ""
	case BumpPrerelease:
		if v.Stage == StageFinal {
			// Bump patch, set stage and seq=nil
			result.Patch++
			result.Stage = stage
			result.Sequence = nil
			result.SequenceType = ""
		} else if v.Stage == stage {
			// Same stage: copy sequence as-is, let sequence calculator handle it
			result.Sequence = v.Sequence
		} else {
			// Different stage: keep core, change stage, seq=nil
			result.Stage = stage
			result.Sequence = nil
			result.SequenceType = ""
		}
	case BumpFinal:
		result.Stage = StageFinal
		result.Sequence = nil
		result.SequenceType = ""
	default:
		return nil, fmt.Errorf("unknown bump kind: %q", kind)
	}

	return &result, nil
}
