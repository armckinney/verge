package fixtures

import "example.com/verge/internal/version"

// GoldenVersion is a test fixture for version parsing.
type GoldenVersion struct {
	Input    string
	Expected *version.Version
	Valid    bool
}

// GoldenVersions contains canonical version parse fixtures.
var GoldenVersions = []GoldenVersion{
	{
		Input: "v1.2.3",
		Expected: &version.Version{
			Major: 1, Minor: 2, Patch: 3,
			Stage:       version.StageFinal,
			VersionType: "semver",
		},
		Valid: true,
	},
	{
		Input: "1.2.3",
		Expected: &version.Version{
			Major: 1, Minor: 2, Patch: 3,
			Stage:       version.StageFinal,
			VersionType: "semver",
		},
		Valid: true,
	},
	{
		Input: "v1.2.3-dev.1",
		Expected: &version.Version{
			Major: 1, Minor: 2, Patch: 3,
			Stage:        version.StageDev,
			Sequence:     1,
			SequenceType: version.SeqTypeNumeric,
			VersionType:  "semver",
		},
		Valid: true,
	},
	{
		Input: "v1.2.3-alpha.2",
		Expected: &version.Version{
			Major: 1, Minor: 2, Patch: 3,
			Stage:        version.StageA,
			Sequence:     2,
			SequenceType: version.SeqTypeNumeric,
			VersionType:  "semver",
		},
		Valid: true,
	},
	{
		Input: "v1.2.3-beta.3",
		Expected: &version.Version{
			Major: 1, Minor: 2, Patch: 3,
			Stage:        version.StageB,
			Sequence:     3,
			SequenceType: version.SeqTypeNumeric,
			VersionType:  "semver",
		},
		Valid: true,
	},
	{
		Input: "v1.2.3-rc.1",
		Expected: &version.Version{
			Major: 1, Minor: 2, Patch: 3,
			Stage:        version.StageRC,
			Sequence:     1,
			SequenceType: version.SeqTypeNumeric,
			VersionType:  "semver",
		},
		Valid: true,
	},
	{
		Input: "1.2.3a1",
		Expected: &version.Version{
			Major: 1, Minor: 2, Patch: 3,
			Stage:        version.StageA,
			Sequence:     1,
			SequenceType: version.SeqTypeNumeric,
			VersionType:  "pep440",
		},
		Valid: true,
	},
	{
		Input: "not-a-version",
		Valid: false,
	},
}
