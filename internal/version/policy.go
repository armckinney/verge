package version

import "fmt"

// PolicyRules defines the rules for version policy checking.
type PolicyRules struct {
	// EnforceMonotonicIncrease requires new > current.
	EnforceMonotonicIncrease bool
	// AllowedStages lists stages permitted for this policy (nil = all).
	AllowedStages []Stage
}

// PolicyChecker validates versions against a set of rules.
type PolicyChecker struct {
	rules      PolicyRules
	comparator Comparator
}

// NewPolicyChecker returns a PolicyChecker with the given rules.
func NewPolicyChecker(rules PolicyRules) *PolicyChecker {
	return &PolicyChecker{
		rules:      rules,
		comparator: NewComparator(),
	}
}

// DefaultPolicy returns a minimal policy that only enforces monotonic increase.
func DefaultPolicy() *PolicyChecker {
	return NewPolicyChecker(PolicyRules{EnforceMonotonicIncrease: true})
}

// Validate checks newVersion against the policy, optionally comparing to currentVersion.
// currentVersion may be nil (first release).
func (p *PolicyChecker) Validate(newVersion, currentVersion *Version) error {
	if newVersion == nil {
		return fmt.Errorf("new version must not be nil")
	}

	if p.rules.EnforceMonotonicIncrease && currentVersion != nil {
		cmp := p.comparator.Compare(newVersion, currentVersion)
		if cmp <= 0 {
			return fmt.Errorf("version %s must be greater than current version %s",
				newVersion.String(), currentVersion.String())
		}
	}

	if len(p.rules.AllowedStages) > 0 {
		allowed := false
		for _, s := range p.rules.AllowedStages {
			if newVersion.Stage == s {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("stage %q is not allowed by policy (allowed: %v)",
				newVersion.Stage.String(), p.rules.AllowedStages)
		}
	}

	return nil
}
