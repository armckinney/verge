package cli

import (
	"fmt"

	"example.com/verge/internal/config"
	"example.com/verge/internal/domain"
	"example.com/verge/internal/types"
	"github.com/spf13/cobra"
)

func versionCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Get the current locally tracked version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(globalFlags.configPath)
			if err != nil {
				return NewError(ExitConfigError, "loading config: %v", err)
			}

			// Delegate domain logic
			v, err := domain.GetCurrent(cfg, "")
			if err != nil {
				return fmt.Errorf("current failed: %w", err)
			}

			// Render and print
			parser := types.Get(cfg.VersionType)
			if parser == nil {
				return fmt.Errorf("invalid version_type setup")
			}
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

	return cmd
}
