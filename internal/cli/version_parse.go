package cli

import (
	"fmt"

	"example.com/template-go/internal/version"
	"github.com/spf13/cobra"
)

func versionParseCmd() *cobra.Command {
	var ecosystem string

	cmd := &cobra.Command{
		Use:   "parse <version>",
		Short: "Parse a version string",
		Long: `Parse and validate a version string, displaying its components and rendered forms.

Examples:
  verctl version parse 1.2.3
  verctl version parse v1.2.3-rc.2
  verctl version parse 1.2.3dev4
  verctl version parse 1.2.3 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			parser := version.NewParser()
			v, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing version %q: %w", args[0], err)
			}

			normalizer := version.NewNormalizer()
			v = normalizer.Normalize(v)

			// Render once per canonical format scheme.
			ecosystemNames := []string{"v-semver", "semver", "pep440"}
			rendered := map[string]string{}
			for _, eco := range ecosystemNames {
				rendered[eco] = version.NewRenderer(eco).Render(v)
			}

			out := NewOutput(OutputFormat(globalFlags.format))
			if out.Format == FormatJSON {
				data := map[string]interface{}{
					"input": args[0],
					"parsed": map[string]interface{}{
						"major":        v.Major,
						"minor":        v.Minor,
						"patch":        v.Patch,
						"stage":        v.Stage.String(),
						"sequence":     v.Sequence,
						"sequenceType": string(v.SequenceType),
					},
					"schemeDetected": string(v.Scheme),
					"rendered":       rendered,
				}
				if ecosystem != "" && ecosystem != "all" {
					data["rendered"] = rendered[ecosystem]
				}
				return out.Print(data)
			}

			data := map[string]interface{}{
				"input":        args[0],
				"major":        v.Major,
				"minor":        v.Minor,
				"patch":        v.Patch,
				"stage":        v.Stage.String(),
				"sequence":     v.Sequence,
				"sequenceType": string(v.SequenceType),
				"scheme":       string(v.Scheme),
				"prerelease":   v.IsPrerelease(),
				"core":         v.Core(),
			}
			if ecosystem != "" && ecosystem != "all" {
				data["rendered"] = rendered[ecosystem]
			} else {
				for eco, r := range rendered {
					data["rendered."+eco] = r
				}
			}
			return out.Print(data)
		},
	}

	cmd.Flags().StringVar(&ecosystem, "ecosystem", "all", "Target format scheme for rendering (v-semver, semver, pep440, or ecosystem alias: go, terraform, containers, github-actions, python)")
	return cmd
}

