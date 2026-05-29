package cli

import (
	"fmt"

	"example.com/verge/internal/config"
	"example.com/verge/internal/domain"
	"example.com/verge/internal/types"
	"github.com/spf13/cobra"
)

func versionLatestCmd() *cobra.Command {
	var (
		versionType string
		providerStr string
		versionArg  string
	)

	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Get the latest version from tracking provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(globalFlags.configPath)
			if err != nil {
				return NewError(ExitConfigError, "loading config: %v", err)
			}
			cfg.NoCache = globalFlags.noCache

			if versionType != "" {
				cfg.VersionType = versionType
			}
			if providerStr != "" {
				cfg.Provider.Type = providerStr
				cfg.Provider.Raw = nil
			}

			opts := domain.LatestOptions{
				Prefix: versionArg,
			}

			v, err := domain.GetLatest(cfg, opts)
			if err != nil {
				return fmt.Errorf("latest failed: %w", err)
			}

			parser := types.Get(cfg.VersionType)
			rendered := parser.Render(v)

			out := NewOutput(OutputFormat(globalFlags.format))
			out.Field = globalFlags.field
			data := map[string]interface{}{
				"version":    v.Original,
				"normalized": v.String(),
				"rendered":   rendered,
			}
			return out.Print(data)
		},
	}

	cmd.Flags().StringVarP(&versionType, "type", "t", "", "Override version_type")
	cmd.Flags().StringVarP(&providerStr, "provider", "p", "", "Override provider type")
	cmd.Flags().StringVar(&versionArg, "version", "", "Prefix filter for version (e.g. 1.2)")

	return cmd
}
