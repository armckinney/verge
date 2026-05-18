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
			// Bump patch, set stage and seq=1
			result.Patch++
			result.Stage = stage
			result.Sequence = 1
			result.SequenceType = SeqTypeNumeric
		} else if v.Stage == stage {
			// Same stage: increment sequence
			switch seq := v.Sequence.(type) {
			case int:
				result.Sequence = seq + 1
			default:
				return nil, fmt.Errorf("cannot increment non-numeric sequence: %v", v.Sequence)
			}
		} else {
			// Different stage: keep core, change stage, seq=1
			result.Stage = stage
			result.Sequence = 1
			result.SequenceType = SeqTypeNumeric
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
