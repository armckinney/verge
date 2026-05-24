package cli

import "github.com/spf13/cobra"

func versionInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show verge version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := NewOutput(OutputFormat(globalFlags.format))
			out.Field = globalFlags.field
			return out.Print(map[string]interface{}{
				"version": versionInfo.Version,
				"commit":  versionInfo.Commit,
				"date":    versionInfo.Date,
			})
		},
	}
}
