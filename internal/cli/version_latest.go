package cli

import (
	"fmt"
	"sort"

	"example.com/template-go/internal/config"
	"example.com/template-go/internal/providers"
	"example.com/template-go/internal/version"
	"github.com/spf13/cobra"
)

func versionLatestCmd() *cobra.Command {
	var (
		repoDir  string
		coreOnly bool
		stageStr string
	)

	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Get the latest version from git tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(globalFlags.configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
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
				return fmt.Errorf("no versions found")
			}

			comparator := version.NewComparator()
			sort.Slice(results, func(i, j int) bool {
				return comparator.Compare(results[i].Version, results[j].Version) > 0
			})

			v := results[0].Version
			out := NewOutput(OutputFormat(globalFlags.format))
			if coreOnly {
				out.PrintValue(v.Core())
			} else {
				out.PrintValue(cfg.Format.TagPrefix + v.String())
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&repoDir, "repo-dir", ".", "Repository directory")
	cmd.Flags().BoolVar(&coreOnly, "core", false, "Output only the core version (M.m.p)")
	cmd.Flags().StringVar(&stageStr, "stage", "", "Filter by stage (dev, alpha, beta, rc, final)")
	return cmd
}
