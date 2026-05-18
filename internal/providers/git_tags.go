package providers

import (
	"os/exec"
	"strings"

	"example.com/template-go/internal/version"
)

type GitTagsProvider struct{}

func NewGitTagsProvider() *GitTagsProvider {
	return &GitTagsProvider{}
}

func (g *GitTagsProvider) Name() string { return "git-tags" }

func (g *GitTagsProvider) Fetch(opts QueryOptions) ([]*VersionResult, error) {
	dir := opts.RepoDir
	if dir == "" {
		dir = "."
	}

	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	parser := version.NewParser()
	prefix := opts.TagPrefix
	if prefix == "" {
		prefix = "v"
	}

	var results []*VersionResult
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		raw := line
		// Try stripping prefix for parsing
		toParse := line
		if strings.HasPrefix(toParse, prefix) {
			toParse = strings.TrimPrefix(toParse, prefix)
			// re-add v for the parser since it expects it for semver
			toParse = prefix + toParse
		}

		v, err := parser.Parse(toParse)
		if err != nil {
			// Try without prefix
			v, err = parser.Parse(line)
			if err != nil {
				continue
			}
		}

		if !opts.IncludePrerelease && v.IsPrerelease() {
			continue
		}

		results = append(results, &VersionResult{
			Version: v,
			Raw:     raw,
			Source:  "git-tags",
		})
	}

	return results, nil
}
