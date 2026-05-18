package cli

import "github.com/spf13/cobra"

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Version management commands",
		Long:  `Commands for parsing, comparing, bumping, and querying versions.`,
	}

	cmd.AddCommand(
		versionParseCmd(),
		versionCompareCmd(),
		versionCurrentCmd(),
		versionLatestCmd(),
		versionBumpCmd(),
		versionInfoCmd(),
	)

	return cmd
}

func versionInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show verctl version info",
		Run: func(cmd *cobra.Command, args []string) {
			out := NewOutput(OutputFormat(globalFlags.format))
			out.Print(map[string]interface{}{
				"version": versionInfo.Version,
				"commit":  versionInfo.Commit,
				"date":    versionInfo.Date,
			})
		},
	}
}
