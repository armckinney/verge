package cli

import (
	"fmt"

	"example.com/template-go/internal/version"
	"github.com/spf13/cobra"
)

func versionBumpCmd() *cobra.Command {
	var stageStr string

	cmd := &cobra.Command{
		Use:   "bump <kind> <version>",
		Short: "Bump a version",
		Long: `Bump a version by a given kind.
Kinds: major, minor, patch, prerelease, final`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			kindStr := args[0]
			versionStr := args[1]

			var kind version.BumpKind
			switch kindStr {
			case "major":
				kind = version.BumpMajor
			case "minor":
				kind = version.BumpMinor
			case "patch":
				kind = version.BumpPatch
			case "prerelease":
				kind = version.BumpPrerelease
			case "final":
				kind = version.BumpFinal
			default:
				return fmt.Errorf("unknown bump kind: %q (use major, minor, patch, prerelease, final)", kindStr)
			}

			parser := version.NewParser()
			v, err := parser.Parse(versionStr)
			if err != nil {
				return fmt.Errorf("parsing version: %w", err)
			}

			stage := version.StageDev
			if stageStr != "" {
				stage, err = version.StageFromString(stageStr)
				if err != nil {
					return fmt.Errorf("invalid stage: %w", err)
				}
			}

			bumper := version.NewBumper()
			bumped, err := bumper.Bump(v, kind, stage)
			if err != nil {
				return fmt.Errorf("bumping version: %w", err)
			}

			out := NewOutput(OutputFormat(globalFlags.format))
			out.PrintValue(bumped.String())
			return nil
		},
	}

	cmd.Flags().StringVar(&stageStr, "stage", "dev", "Prerelease stage (dev, alpha, beta, rc)")
	return cmd
}
