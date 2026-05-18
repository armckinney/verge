package cli

import (
	"github.com/spf13/cobra"
)

var (
	versionInfo struct {
		Version string
		Commit  string
		Date    string
	}

	globalFlags struct {
		verbose    bool
		configPath string
		format     string
	}
)

// SetVersionInfo sets the version info injected via ldflags.
func SetVersionInfo(version, commit, date string) {
	versionInfo.Version = version
	versionInfo.Commit = commit
	versionInfo.Date = date
}

var rootCmd = &cobra.Command{
	Use:   "verctl",
	Short: "A semantic versioning CLI tool",
	Long:  `verctl is a semantic versioning CLI tool for managing and bumping versions across ecosystems.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.configPath, "config", "c", "", "Config file path (default: .verctl.yaml)")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.format, "format", "f", "text", "Output format: text or json")

	rootCmd.AddCommand(versionCmd())
}
