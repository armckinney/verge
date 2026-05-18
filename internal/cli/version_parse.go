package cli

import (
	"fmt"

	"example.com/template-go/internal/version"
	"github.com/spf13/cobra"
)

func versionParseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "parse <version>",
		Short: "Parse a version string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			parser := version.NewParser()
			v, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing version %q: %w", args[0], err)
			}

			normalizer := version.NewNormalizer()
			v = normalizer.Normalize(v)

			out := NewOutput(OutputFormat(globalFlags.format))
			data := map[string]interface{}{
				"major":        v.Major,
				"minor":        v.Minor,
				"patch":        v.Patch,
				"stage":        v.Stage.String(),
				"sequence":     v.Sequence,
				"sequenceType": string(v.SequenceType),
				"scheme":       string(v.Scheme),
				"prerelease":   v.IsPrerelease(),
				"core":         v.Core(),
				"string":       v.String(),
			}
			return out.Print(data)
		},
	}
}
