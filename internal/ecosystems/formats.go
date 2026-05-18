package ecosystems

import "fmt"

// VSemVerRenderer renders versions with a `v` prefix (SemVer with v-prefix).
// Used by: go, terraform
type VSemVerRenderer struct{}

func (r *VSemVerRenderer) Name() string { return "v-semver" }

func (r *VSemVerRenderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	}
	return fmt.Sprintf("v%d.%d.%d-%s.%v", major, minor, patch, stage, sequence)
}

// SemVerRenderer renders versions as standard SemVer (no v-prefix).
// Used by: containers, github-actions
type SemVerRenderer struct{}

func (r *SemVerRenderer) Name() string { return "semver" }

func (r *SemVerRenderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s.%v", major, minor, patch, stage, sequence)
}

// PEP440Renderer renders versions in Python PEP 440 format.
// Used by: python
type PEP440Renderer struct{}

func (r *PEP440Renderer) Name() string { return "pep440" }

func (r *PEP440Renderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}
	abbr := pep440StageAbbr(stage)
	return fmt.Sprintf("%d.%d.%d%s%v", major, minor, patch, abbr, sequence)
}

func pep440StageAbbr(stage string) string {
	switch stage {
	case "dev":
		return "dev"
	case "alpha":
		return "a"
	case "beta":
		return "b"
	case "rc":
		return "rc"
	}
	return stage
}
