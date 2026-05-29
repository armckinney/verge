package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func versionInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show verge version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(os.Stdout, versionInfo.Version)
			return nil
		},
	}
}
