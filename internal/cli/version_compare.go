package cli

import (
	"fmt"
	"os"

	"example.com/template-go/internal/version"
	"github.com/spf13/cobra"
)

func versionCompareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "compare <left> <right>",
		Short: "Compare two version strings",
		Long: `Compare two version strings.
Exit codes:
  0  - versions are equal
  10 - left is less than right
  11 - left is greater than right`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			parser := version.NewParser()
			left, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing left version: %w", err)
			}
			right, err := parser.Parse(args[1])
			if err != nil {
				return fmt.Errorf("parsing right version: %w", err)
			}

			comparator := version.NewComparator()
			result := comparator.Compare(left, right)

			out := NewOutput(OutputFormat(globalFlags.format))
			switch result {
			case -1:
				out.PrintValue("less")
				os.Exit(ExitCompareLeft)
			case 1:
				out.PrintValue("greater")
				os.Exit(ExitCompareRight)
			default:
				out.PrintValue("equal")
			}
			return nil
		},
	}
}
