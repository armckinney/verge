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
	var repoDir string

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Get the current version from git tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(globalFlags.configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
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
				return fmt.Errorf("no versions found")
			}

			// Sort and get latest
			comparator := version.NewComparator()
			sort.Slice(results, func(i, j int) bool {
				return comparator.Compare(results[i].Version, results[j].Version) > 0
			})

			v := results[0].Version
			out := NewOutput(OutputFormat(globalFlags.format))
			out.PrintValue(cfg.Format.TagPrefix + v.String())
			return nil
		},
	}

	cmd.Flags().StringVar(&repoDir, "repo-dir", ".", "Repository directory")
	return cmd
}
