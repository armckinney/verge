package cli

import (
	"fmt"

	"example.com/verge/internal/config"
	"example.com/verge/internal/domain"
	"example.com/verge/internal/types"
	"github.com/spf13/cobra"
)

func versionBumpCmd() *cobra.Command {
	var (
		versionType    string
		providerStr    string
		versionArg     string // Base string
		prefixArg      string // Base prefix
		kindStr        string
		stageStr       string
		sequenceStr    string
		providerConfig []string
	)

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump a version based on deterministic rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(globalFlags.configPath)
			if err != nil {
				return NewError(ExitConfigError, "loading config: %v", err)
			}
			cfg.NoCache = globalFlags.noCache

			// CLI Overrides Config
			if versionType != "" {
				cfg.VersionType = versionType
			}
			if providerStr != "" {
				cfg.Provider.Type = providerStr
				cfg.Provider.Raw = nil
			}

			// Apply inline provider config overrides
			if len(providerConfig) > 0 {
				overrides := config.ParseOverrides(providerConfig)
				config.MergeOverrides(&cfg.Provider, overrides)
			}

			opts := domain.BumpOptions{
				VersionStr:      versionArg,
				Prefix:          prefixArg,
				PrereleaseStage: stageStr,
				BumpKind:        kindStr,
				SequenceStr:     sequenceStr,
			}

			bumped, err := domain.Bump(cfg, opts)
			if err != nil {
				return fmt.Errorf("bump failed: %w", err)
			}

			parser := types.Get(cfg.VersionType)
			rendered := parser.Render(bumped)

			out := NewOutput(OutputFormat(globalFlags.format))
			out.Field = globalFlags.field
			data := map[string]interface{}{
				"kind":     kindStr,
				"to":       bumped.String(),
				"rendered": rendered,
			}
			if stageStr != "" {
				data["stage"] = stageStr
			}
			return out.Print(data)
		},
	}

	cmd.Flags().StringVarP(&versionType, "type", "t", "", "Override version_type")
	cmd.Flags().StringVarP(&providerStr, "provider", "p", "", "Override provider type")
	cmd.Flags().StringVar(&versionArg, "version", "", "Bypass fetch and use static version to bump")
	cmd.Flags().StringVar(&prefixArg, "prefix", "", "Prefix filter fetching the latest tracking version")
	cmd.Flags().StringVar(&kindStr, "kind", "", "Bump kind: major, minor, patch, prerelease, final")
	cmd.Flags().StringVar(&stageStr, "stage", "", "Prerelease stage (dev, a, b, rc)")
	cmd.Flags().StringVarP(&sequenceStr, "sequence", "s", "", "Static sequence value to override calculators")
	cmd.Flags().StringSliceVar(&providerConfig, "provider-config", nil, "Provider config overrides (e.g. key=val)")

	return cmd
}
