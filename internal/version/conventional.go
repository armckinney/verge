package version

import (
	"fmt"
	"os/exec"
	"strings"
)

// CommitType represents the conventional commit type.
type CommitType string

const (
	CommitFeat     CommitType = "feat"
	CommitFix      CommitType = "fix"
	CommitBreaking CommitType = "breaking"
	CommitOther    CommitType = "other"
)

// ConventionalCommit holds parsed commit data.
type ConventionalCommit struct {
	Hash     string
	Message  string
	Type     CommitType
	Scope    string
	Subject  string
	Breaking bool
}

// CommitParser parses conventional commit messages.
type CommitParser struct {
	breakingTokens []string
}

// NewCommitParser returns a parser with default breaking change tokens.
func NewCommitParser(breakingTokens ...string) *CommitParser {
	if len(breakingTokens) == 0 {
		breakingTokens = []string{"BREAKING CHANGE", "BREAKING-CHANGE"}
	}
	return &CommitParser{breakingTokens: breakingTokens}
}

// Parse parses a single commit line ("hash message").
func (p *CommitParser) Parse(hash, message string) ConventionalCommit {
	c := ConventionalCommit{Hash: hash, Message: message}

	// Check for breaking change footer or ! marker
	for _, token := range p.breakingTokens {
		if strings.Contains(message, token) {
			c.Breaking = true
			c.Type = CommitBreaking
			return c
		}
	}

	// Parse "type(scope): subject" or "type!: subject"
	msg := message
	if idx := strings.Index(msg, ":"); idx > 0 {
		prefix := msg[:idx]
		subject := strings.TrimSpace(msg[idx+1:])
		c.Subject = subject

		// Handle scope
		if si := strings.Index(prefix, "("); si >= 0 {
			c.Scope = strings.Trim(prefix[si:], "()")
			prefix = prefix[:si]
		}

		// Check for ! (breaking)
		if strings.HasSuffix(prefix, "!") {
			c.Breaking = true
			c.Type = CommitBreaking
			prefix = strings.TrimSuffix(prefix, "!")
		}

		switch strings.ToLower(strings.TrimSpace(prefix)) {
		case "feat", "feature":
			if !c.Breaking {
				c.Type = CommitFeat
			}
		case "fix", "bugfix":
			if !c.Breaking {
				c.Type = CommitFix
			}
		default:
			if !c.Breaking {
				c.Type = CommitOther
			}
		}
	} else {
		c.Type = CommitOther
	}

	return c
}

// CommitHistory holds commits between two refs.
type CommitHistory struct {
	Commits []ConventionalCommit
}

// DetectedBump returns the appropriate BumpKind for this commit history.
func (h *CommitHistory) DetectedBump() (BumpKind, error) {
	if len(h.Commits) == 0 {
		return "", fmt.Errorf("no commits found")
	}

	hasBreaking := false
	hasFeat := false
	hasFix := false

	for _, c := range h.Commits {
		switch c.Type {
		case CommitBreaking:
			hasBreaking = true
		case CommitFeat:
			hasFeat = true
		case CommitFix:
			hasFix = true
		}
	}

	switch {
	case hasBreaking:
		return BumpMajor, nil
	case hasFeat:
		return BumpMinor, nil
	case hasFix:
		return BumpPatch, nil
	default:
		return "", fmt.Errorf("no feat/fix/breaking commits found")
	}
}

// FetchCommitsSince returns commits between fromRef (e.g. a tag) and HEAD.
// If fromRef is empty, all commits are returned.
func FetchCommitsSince(repoDir, fromRef string, parser *CommitParser) (*CommitHistory, error) {
	var args []string
	if fromRef != "" {
		args = []string{"log", "--format=%H %s", fromRef + "..HEAD"}
	} else {
		args = []string{"log", "--format=%H %s"}
	}

	cmd := exec.Command("git", args...)
	if repoDir != "" {
		cmd.Dir = repoDir
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git log: %w", err)
	}

	var commits []ConventionalCommit
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		hash := parts[0]
		message := ""
		if len(parts) > 1 {
			message = parts[1]
		}
		commits = append(commits, parser.Parse(hash, message))
	}

	return &CommitHistory{Commits: commits}, nil
}
