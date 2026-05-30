package cli

import (
	"example.com/verge/internal/config"
	"example.com/verge/internal/types"
	"example.com/verge/internal/version"
	"github.com/spf13/cobra"
)

func versionParseCmd() *cobra.Command {
	var versionType string

	cmd := &cobra.Command{
		Use:   "parse [version]",
		Short: "Parse a version string into structured components",
		Long:  `parse extracts semantic components (major, minor, patch, stage, sequence) from a version string.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versionArg := args[0]

			cfg, err := config.Load(globalFlags.configPath)
			if err != nil {
				return NewError(ExitConfigError, "loading config: %v", err)
			}
			cfg.NoCache = globalFlags.noCache

			var v *version.Version
			var parseErr error
			var chosenType string

			// If type is explicitly overridden
			if versionType != "" {
				chosenType = versionType
				parser := types.Get(chosenType)
				if parser == nil {
					return NewError(ExitConfigError, "invalid version type: %q", chosenType)
				}
				v, parseErr = parser.Parse(versionArg)
				if parseErr != nil {
					return NewError(ExitError, "parsing version %q: %v", versionArg, parseErr)
				}
			} else {
				// 1. Try the configured version type first
				chosenType = cfg.VersionType
				parser := types.Get(chosenType)
				if parser != nil {
					v, parseErr = parser.Parse(versionArg)
				}

				// 2. If it fails, try other available types
				if parseErr != nil || parser == nil {
					allTypes := []string{"semver", "vsemver", "pep440"}
					for _, t := range allTypes {
						if t == cfg.VersionType {
							continue
						}
						p := types.Get(t)
						if p == nil {
							continue
						}
						if parsed, err := p.Parse(versionArg); err == nil {
							v = parsed
							chosenType = t
							parseErr = nil
							break
						}
					}
				}

				if v == nil {
					return NewError(ExitError, "unable to parse version %q with any supported parser", versionArg)
				}
			}

			parser := types.Get(chosenType)
			rendered := parser.Render(v)

			defaultFormat := OutputFormat(globalFlags.format)
			if !cmd.Flags().Changed("format") && !globalFlags.json {
				defaultFormat = FormatJSON
			}
			out := NewOutput(defaultFormat)
			out.Field = globalFlags.field

			data := map[string]interface{}{
				"major":         v.Major,
				"minor":         v.Minor,
				"patch":         v.Patch,
				"stage":         v.Stage.String(),
				"sequence":      v.Sequence,
				"sequence_type": string(v.SequenceType),
				"is_prerelease": v.IsPrerelease(),
				"core":          v.Core(),
				"version_type":  v.VersionType,
				"version":       v.Original,
				"normalized":    v.String(),
				"rendered":      rendered,
			}

			return out.Print(data)
		},
	}

	cmd.Flags().StringVarP(&versionType, "type", "t", "", "Override version_type (semver | vsemver | pep440)")

	return cmd
}
