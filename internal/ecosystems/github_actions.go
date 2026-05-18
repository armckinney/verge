package ecosystems

import "fmt"

type GitHubActionsRenderer struct{}

func (g *GitHubActionsRenderer) Name() string { return "github-actions" }

func (g *GitHubActionsRenderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s.%v", major, minor, patch, stage, sequence)
}
