package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var templateName string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate configuration boilerplate",
		RunE: func(cmd *cobra.Command, args []string) error {
			content := "version_type: semver\nprovider:\n  type: gittag\n"

			switch templateName {
			case "gittag-semver":
				// Used default
			case "gittag-vsemver":
				content = "version_type: vsemver\nprovider:\n  type: gittag\n"
			case "ghrelease-semver":
				content = "version_type: semver\nprovider:\n  type: ghrelease\n"
			default:
				if templateName != "" {
					return NewError(ExitError, "unknown template name: %q", templateName)
				}
			}

			target := ".verge.yaml"
			if _, err := os.Stat(target); err == nil {
				target = ".verge." + templateName + ".yaml"
				if templateName == "" {
					target = ".verge.generated.yaml"
				}
			}

			if err := os.WriteFile(target, []byte(content), 0644); err != nil {
				return NewError(ExitError, "writing config file: %v", err)
			}

			fmt.Fprintf(os.Stderr, "Wrote config to %s\n", target)
			return nil
		},
	}

	cmd.Flags().StringVar(&templateName, "template", "", "Template name (e.g. gittag-semver, ghrelease-semver)")
	return cmd
}
