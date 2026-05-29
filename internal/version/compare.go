package version

type defaultComparator struct{}

// NewComparator returns a new Comparator.
func NewComparator() Comparator {
	return &defaultComparator{}
}

var stageOrder = map[Stage]int{
	StageDev:   0,
	StageA: 1,
	StageB:  2,
	StageRC:    3,
	StageFinal: 4,
}

func (c *defaultComparator) Compare(a, b *Version) int {
	if a.Major != b.Major {
		return cmpInt(a.Major, b.Major)
	}
	if a.Minor != b.Minor {
		return cmpInt(a.Minor, b.Minor)
	}
	if a.Patch != b.Patch {
		return cmpInt(a.Patch, b.Patch)
	}

	// Compare stages
	sa := stageOrder[a.Stage]
	sb := stageOrder[b.Stage]
	if sa != sb {
		return cmpInt(sa, sb)
	}

	// Same stage - compare sequence
	return compareSequence(a.Sequence, b.Sequence)
}

func cmpInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareSequence(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	switch av := a.(type) {
	case int:
		switch bv := b.(type) {
		case int:
			return cmpInt(av, bv)
		case string:
			return 1 // numeric > string (arbitrary but consistent)
		}
	case string:
		switch bv := b.(type) {
		case string:
			if av < bv {
				return -1
			}
			if av > bv {
				return 1
			}
			return 0
		case int:
			return -1
		}
	}
	return 0
}
