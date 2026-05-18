package cli

import (
	"fmt"
	"sort"

	"example.com/template-go/internal/config"
	"example.com/template-go/internal/providers"
	"example.com/template-go/internal/version"
	"github.com/spf13/cobra"
)

func versionCurrentCmd() *cobra.Command {
	var (
		repoDir   string
		ecosystem string
		explain   bool
	)

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Get the current version from git tags",
		Long: `Fetch the current (highest) version from the configured version source.

Examples:
  verge version current
  verge version current --ecosystem python
  verge version current --explain
  verge version current --format json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(globalFlags.configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if ecosystem == "" {
				ecosystem = cfg.Ecosystem
			}

			provider := providers.NewGitTagsProvider()
			results, err := provider.Fetch(providers.QueryOptions{
				IncludePrerelease: cfg.Sources.GitTags.IncludePrerelease,
				TagPrefix:         cfg.Format.TagPrefix,
				RepoDir:           repoDir,
			})
			if err != nil {
				return fmt.Errorf("fetching git tags: %w", err)
			}

			if len(results) == 0 {
				return fmt.Errorf("no versions found in git tags (exit code %d)", ExitNotFound)
			}

			comparator := version.NewComparator()
			sort.Slice(results, func(i, j int) bool {
				return comparator.Compare(results[i].Version, results[j].Version) > 0
			})

			if explain {
				fmt.Println("Candidates from git-tags:")
				for _, r := range results {
					if r.Version.IsPrerelease() {
						fmt.Printf("  %s (prerelease)\n", r.Raw)
					} else {
						fmt.Printf("  %s (final)\n", r.Raw)
					}
				}
				fmt.Println()
			}

			v := results[0].Version
			rendered := version.NewRenderer(ecosystem).Render(v)

			out := NewOutput(OutputFormat(globalFlags.format))
			data := map[string]interface{}{
				"version":    results[0].Raw,
				"normalized": v.String(),
				"ecosystem":  ecosystem,
				"source":     "git-tags",
				"rendered":   rendered,
			}
			return out.Print(data)
		},
	}

	cmd.Flags().StringVar(&repoDir, "repo-dir", ".", "Repository directory")
	cmd.Flags().StringVar(&ecosystem, "ecosystem", "", "Target format scheme for rendering (v-semver, semver, pep440, or ecosystem alias: go, terraform, containers, github-actions, python)")
	cmd.Flags().BoolVar(&explain, "explain", false, "Show selection reasoning")
	return cmd
}

