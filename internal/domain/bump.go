package domain

import (
	"fmt"

	"example.com/verge/internal/config"
	"example.com/verge/internal/sequence"
	"example.com/verge/internal/types"
	"example.com/verge/internal/version"
)

// BumpOptions holds overrides from CLI.
type BumpOptions struct {
	VersionStr      string // Overrides fetching -- parses linearly
	Prefix          string // Prefix filter fetching the latest tracking version
	PrereleaseStage string
	BumpKind        string
	SequenceStr     string
}

// Bump executes the bump action.
func Bump(cfg *config.Config, opts BumpOptions) (*version.Version, error) {
	// 1. Determine Initial Base Version
	var baseVersion *version.Version

	parser := types.Get(cfg.VersionType)
	if parser == nil {
		return nil, fmt.Errorf("invalid version_type: %s", cfg.VersionType)
	}

	if opts.VersionStr != "" {
		// Bypass fetching: process linearly from string
		v, err := parser.Parse(opts.VersionStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse explicit version %q: %w", opts.VersionStr, err)
		}
		baseVersion = v
	} else {
		// Need to fetch from provider to know what to bump
		latestOpts := LatestOptions{
			Prefix: opts.Prefix,
		}
		
		v, err := GetLatest(cfg, latestOpts)
		if err != nil {
			// Init behavior if no version exists
			// "Implement Initialization Behaviors: If a version or sequence doesn't exist, calculate the initial state safely (e.g., default to 0.1.0 for first time use)."
			baseVersion = &version.Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
				Stage: version.StageFinal,
			}
		} else {
			baseVersion = v
		}
	}

	// 2. Perform the Bump
	bumper := version.NewBumper()

	// Priority: CLI > Config File > Default
	kindStr := opts.BumpKind
	if kindStr == "" {
		kindStr = cfg.Default.BumpKind
	}

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
		return nil, fmt.Errorf("unknown bump kind: %s", kindStr)
	}

	stageStr := opts.PrereleaseStage
	if stageStr == "" {
		stageStr = cfg.Default.PrereleaseStage
	}

	var stage version.Stage
	if stageStr != "" {
		parsedStage, err := version.StageFromString(stageStr)
		if err != nil {
			return nil, fmt.Errorf("invalid prerelease stage: %w", err)
		}
		stage = parsedStage
	}

	bumped, err := bumper.Bump(baseVersion, kind, stage)
	if err != nil {
		return nil, fmt.Errorf("bumping version: %w", err)
	}

	// 3. Process Sequence
	if kind == version.BumpPrerelease || bumped.Stage != version.StageFinal {
		// Calculate the next sequence using the configured calculator
		seqEngine, err := sequence.GetCalculator(opts.SequenceStr, cfg.Sequence)
		if err != nil {
			return nil, fmt.Errorf("failed to get sequence calculator: %w", err)
		}

		newSeq, err := seqEngine.Calculate(bumped.Sequence)
		if err != nil {
			return nil, fmt.Errorf("sequence calculation failed: %w", err)
		}
		bumped.Sequence = newSeq
		bumped.SequenceType = version.SeqTypeUnknown
	} else {
		// Explicitly wipe sequence for finals
		bumped.Sequence = nil
		bumped.SequenceType = ""
	}

	return bumped, nil
}
