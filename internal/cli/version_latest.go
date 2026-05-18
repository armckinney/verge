package cli

import (
	"fmt"
	"sort"

	"example.com/verge/internal/config"
	"example.com/verge/internal/providers"
	"example.com/verge/internal/version"
	"github.com/spf13/cobra"
)

func versionLatestCmd() *cobra.Command {
	var (
		repoDir   string
		coreStr   string
		stageStr  string
		ecosystem string
		explain   bool
	)

	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Get the latest version from git tags",
		Long: `Fetch the latest (highest) version from the configured version source.

Supports filtering by core version and prerelease stage.

Examples:
  verge version latest
  verge version latest --stage rc
  verge version latest --core 1.2.3 --stage dev
  verge version latest --explain
  verge version latest --format json`,
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
				IncludePrerelease: true,
				TagPrefix:         cfg.Format.TagPrefix,
				RepoDir:           repoDir,
			})
			if err != nil {
				return fmt.Errorf("fetching git tags: %w", err)
			}

			// Filter by core version if specified
			if coreStr != "" {
				var filtered []*providers.VersionResult
				for _, r := range results {
					if r.Version.Core() == coreStr {
						filtered = append(filtered, r)
					}
				}
				results = filtered
			}

			// Filter by stage if specified
			if stageStr != "" {
				targetStage, err := version.StageFromString(stageStr)
				if err != nil {
					return fmt.Errorf("invalid stage: %w", err)
				}
				var filtered []*providers.VersionResult
				for _, r := range results {
					if r.Version.Stage == targetStage {
						filtered = append(filtered, r)
					}
				}
				results = filtered
			}

			if len(results) == 0 {
				return fmt.Errorf("no matching versions found (exit code %d)", ExitNotFound)
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
				fmt.Printf("\nSelected: %s (highest version)\n\n", results[0].Raw)
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
				"core":       coreStr,
				"stage":      stageStr,
			}
			return out.Print(data)
		},
	}

	cmd.Flags().StringVar(&repoDir, "repo-dir", ".", "Repository directory")
	cmd.Flags().StringVar(&coreStr, "core", "", "Filter by core version (e.g., 1.2.3)")
	cmd.Flags().StringVar(&stageStr, "stage", "", "Filter by stage (dev, alpha, beta, rc, final)")
	cmd.Flags().StringVar(&ecosystem, "ecosystem", "", "Target ecosystem for rendering")
	cmd.Flags().BoolVar(&explain, "explain", false, "Show filtering and selection reasoning")
	return cmd
}
