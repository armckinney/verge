package version_test

import (
	"testing"

	"example.com/template-go/internal/version"
)

func TestCommitParser_Feat(t *testing.T) {
	p := version.NewCommitParser()
	c := p.Parse("abc1234", "feat: add new feature")
	if c.Type != version.CommitFeat {
		t.Errorf("expected CommitFeat, got %v", c.Type)
	}
}

func TestCommitParser_Fix(t *testing.T) {
	p := version.NewCommitParser()
	c := p.Parse("abc1234", "fix: correct nil pointer")
	if c.Type != version.CommitFix {
		t.Errorf("expected CommitFix, got %v", c.Type)
	}
}

func TestCommitParser_Breaking_FooterToken(t *testing.T) {
	p := version.NewCommitParser()
	c := p.Parse("abc1234", "feat: something\n\nBREAKING CHANGE: removes API")
	if c.Type != version.CommitBreaking {
		t.Errorf("expected CommitBreaking, got %v", c.Type)
	}
}

func TestCommitParser_Breaking_BangSuffix(t *testing.T) {
	p := version.NewCommitParser()
	c := p.Parse("abc1234", "feat!: drop support for node 6")
	if c.Type != version.CommitBreaking {
		t.Errorf("expected CommitBreaking, got %v", c.Type)
	}
}

func TestCommitParser_WithScope(t *testing.T) {
	p := version.NewCommitParser()
	c := p.Parse("abc1234", "feat(api): new endpoint")
	if c.Type != version.CommitFeat {
		t.Errorf("expected CommitFeat, got %v", c.Type)
	}
	if c.Scope != "api" {
		t.Errorf("expected scope 'api', got %q", c.Scope)
	}
}

func TestCommitParser_Other(t *testing.T) {
	p := version.NewCommitParser()
	c := p.Parse("abc1234", "docs: update README")
	if c.Type != version.CommitOther {
		t.Errorf("expected CommitOther, got %v", c.Type)
	}
}

func TestCommitParser_CustomBreakingToken(t *testing.T) {
	p := version.NewCommitParser("CUSTOM_BREAKING")
	c := p.Parse("abc1234", "feat: change CUSTOM_BREAKING api")
	if c.Type != version.CommitBreaking {
		t.Errorf("expected CommitBreaking with custom token, got %v", c.Type)
	}
}

func TestCommitHistory_DetectedBump_Breaking(t *testing.T) {
	p := version.NewCommitParser()
	h := &version.CommitHistory{
		Commits: []version.ConventionalCommit{
			p.Parse("a", "feat!: breaking change"),
			p.Parse("b", "fix: minor fix"),
		},
	}
	kind, err := h.DetectedBump()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kind != version.BumpMajor {
		t.Errorf("expected BumpMajor, got %v", kind)
	}
}

func TestCommitHistory_DetectedBump_Minor(t *testing.T) {
	p := version.NewCommitParser()
	h := &version.CommitHistory{
		Commits: []version.ConventionalCommit{
			p.Parse("a", "feat: new feature"),
			p.Parse("b", "fix: minor fix"),
		},
	}
	kind, err := h.DetectedBump()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kind != version.BumpMinor {
		t.Errorf("expected BumpMinor, got %v", kind)
	}
}

func TestCommitHistory_DetectedBump_Patch(t *testing.T) {
	p := version.NewCommitParser()
	h := &version.CommitHistory{
		Commits: []version.ConventionalCommit{
			p.Parse("a", "fix: typo"),
			p.Parse("b", "chore: update deps"),
		},
	}
	kind, err := h.DetectedBump()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kind != version.BumpPatch {
		t.Errorf("expected BumpPatch, got %v", kind)
	}
}

func TestCommitHistory_DetectedBump_NoRelevant(t *testing.T) {
	p := version.NewCommitParser()
	h := &version.CommitHistory{
		Commits: []version.ConventionalCommit{
			p.Parse("a", "chore: update deps"),
			p.Parse("b", "docs: update README"),
		},
	}
	_, err := h.DetectedBump()
	if err == nil {
		t.Error("expected error for no relevant commits")
	}
}

func TestCommitHistory_DetectedBump_Empty(t *testing.T) {
	h := &version.CommitHistory{}
	_, err := h.DetectedBump()
	if err == nil {
		t.Error("expected error for empty history")
	}
}
