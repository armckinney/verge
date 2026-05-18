package cli

import (
	"fmt"

	"example.com/template-go/internal/version"
	"github.com/spf13/cobra"
)

func versionBumpCmd() *cobra.Command {
	var (
		fromVersion string
		kindStr     string
		stageStr    string
		ecosystem   string
		autoDetect  bool
		repoDir     string
		changelog   bool
	)

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump a version",
		Long: `Compute the next version from a given version and bump kind.

Kinds: major, minor, patch, prerelease, final

Examples:
  verctl version bump --from 1.2.3 --kind minor
  verctl version bump --from 1.2.3 --kind prerelease --stage dev
  verctl version bump --from 1.2.3-rc.1 --kind final
  verctl version bump --from 1.2.3 --auto
  verctl version bump --from 1.2.3 --kind minor --changelog --format json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromVersion == "" {
				return fmt.Errorf("--from flag is required")
			}

			parser := version.NewParser()
			v, err := parser.Parse(fromVersion)
			if err != nil {
				return fmt.Errorf("parsing version: %w", err)
			}

			var kind version.BumpKind

			if autoDetect {
				commitParser := version.NewCommitParser()
				history, err := version.FetchCommitsSince(repoDir, fromVersion, commitParser)
				if err != nil {
					return fmt.Errorf("reading commit history: %w", err)
				}
				kind, err = history.DetectedBump()
				if err != nil {
					return fmt.Errorf("detecting bump kind: %w", err)
				}
				kindStr = string(kind)
			} else {
				if kindStr == "" {
					return fmt.Errorf("--kind flag is required (or use --auto)")
				}
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
					return fmt.Errorf("unknown bump kind %q (use: major, minor, patch, prerelease, final)", kindStr)
				}
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

			if ecosystem == "" {
				ecosystem = "go"
			}
			rendered := version.NewRenderer(ecosystem).Render(bumped)

			if changelog {
				return PrintChangelog(nil, fromVersion, bumped.String(), kindStr, "version-bump", nil)
			}

			out := NewOutput(OutputFormat(globalFlags.format))
			data := map[string]interface{}{
				"from":      fromVersion,
				"kind":      kindStr,
				"to":        bumped.String(),
				"ecosystem": ecosystem,
				"rendered":  rendered,
			}
			if stageStr != "" {
				data["stage"] = stageStr
			}
			return out.Print(data)
		},
	}

	cmd.Flags().StringVar(&fromVersion, "from", "", "Source version to bump from (required)")
	cmd.Flags().StringVar(&kindStr, "kind", "", "Bump kind: major, minor, patch, prerelease, final")
	cmd.Flags().StringVar(&stageStr, "stage", "", "Prerelease stage for prerelease bumps (dev, alpha, beta, rc)")
	cmd.Flags().StringVar(&ecosystem, "ecosystem", "go", "Target ecosystem for rendering (go, python, containers, terraform, github-actions)")
	cmd.Flags().BoolVar(&autoDetect, "auto", false, "Auto-detect bump kind from conventional commits since --from tag")
	cmd.Flags().StringVar(&repoDir, "repo-dir", ".", "Repository directory (used with --auto)")
	cmd.Flags().BoolVar(&changelog, "changelog", false, "Output changelog-friendly JSON instead of default output")
	return cmd
}

