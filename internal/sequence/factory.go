package sequence

import (
	"fmt"

	"example.com/verge/internal/config"
)

// GetCalculator returns the configured Sequence Calculator to use.
func GetCalculator(cliOverrideStr string, cfg config.SequenceConfig) (Calculator, error) {
	if cliOverrideStr != "" {
		return &PassedValueCalculator{Value: cliOverrideStr}, nil
	}

	switch cfg.Type {
	case "increment", "": // Optional default
		return &IncrementCalculator{}, nil
	case "filehash":
		return &FileHashCalculator{
			Targets: cfg.Targets,
			Length:  cfg.Length,
		}, nil
	case "passed":
		// Expects the CLI to pass the value, but since cliOverrideStr was empty:
		return nil, fmt.Errorf("sequence type 'passed' configured but no sequence provided via CLI")
	default:
		return nil, fmt.Errorf("unknown sequence calculator type: %s", cfg.Type)
	}
}
