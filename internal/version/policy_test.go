package version_test

import (
	"testing"

	"example.com/verge/internal/version"
)

func TestPolicyChecker_MonotonicIncrease_Pass(t *testing.T) {
	parser := version.NewParser()
	checker := version.DefaultPolicy()

	current, _ := parser.Parse("1.2.3")
	next, _ := parser.Parse("1.3.0")

	if err := checker.Validate(next, current); err != nil {
		t.Errorf("unexpected policy violation: %v", err)
	}
}

func TestPolicyChecker_MonotonicIncrease_Fail(t *testing.T) {
	parser := version.NewParser()
	checker := version.DefaultPolicy()

	current, _ := parser.Parse("1.3.0")
	next, _ := parser.Parse("1.2.3")

	if err := checker.Validate(next, current); err == nil {
		t.Error("expected policy violation for downgrade, got none")
	}
}

func TestPolicyChecker_MonotonicIncrease_Equal(t *testing.T) {
	parser := version.NewParser()
	checker := version.DefaultPolicy()

	current, _ := parser.Parse("1.2.3")
	next, _ := parser.Parse("1.2.3")

	if err := checker.Validate(next, current); err == nil {
		t.Error("expected policy violation for equal version, got none")
	}
}

func TestPolicyChecker_NilCurrentAllowed(t *testing.T) {
	parser := version.NewParser()
	checker := version.DefaultPolicy()

	next, _ := parser.Parse("1.0.0")

	if err := checker.Validate(next, nil); err != nil {
		t.Errorf("unexpected error for nil current: %v", err)
	}
}

func TestPolicyChecker_AllowedStages_Pass(t *testing.T) {
	parser := version.NewParser()
	checker := version.NewPolicyChecker(version.PolicyRules{
		AllowedStages: []version.Stage{version.StageFinal},
	})

	v, _ := parser.Parse("1.2.3")
	if err := checker.Validate(v, nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPolicyChecker_AllowedStages_Fail(t *testing.T) {
	parser := version.NewParser()
	checker := version.NewPolicyChecker(version.PolicyRules{
		AllowedStages: []version.Stage{version.StageFinal},
	})

	v, _ := parser.Parse("1.2.3-rc.1")
	if err := checker.Validate(v, nil); err == nil {
		t.Error("expected stage policy violation, got none")
	}
}

func TestPolicyChecker_NilVersion_Error(t *testing.T) {
	checker := version.DefaultPolicy()
	if err := checker.Validate(nil, nil); err == nil {
		t.Error("expected error for nil version, got none")
	}
}
