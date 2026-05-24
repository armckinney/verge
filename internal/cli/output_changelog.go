package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// ChangelogVersionInfo describes a version transition.
type ChangelogVersionInfo struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to"`
	BumpType string `json:"bumpType,omitempty"`
}

// ChangelogCommit represents a commit in the changelog.
type ChangelogCommit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ChangelogMetadata holds provenance for a changelog entry.
type ChangelogMetadata struct {
	Timestamp string            `json:"timestamp"`
	Source    string            `json:"source"`
	Commits   []ChangelogCommit `json:"commits,omitempty"`
}

// ChangelogOutput is the top-level changelog-friendly JSON structure.
type ChangelogOutput struct {
	Version  ChangelogVersionInfo `json:"version"`
	Metadata ChangelogMetadata    `json:"metadata"`
}

// PrintChangelog writes structured changelog JSON to w (defaults to stdout).
func PrintChangelog(w io.Writer, field, from, to, bumpType, source string, commits []ChangelogCommit) error {
	if w == nil {
		w = os.Stdout
	}
	out := ChangelogOutput{
		Version: ChangelogVersionInfo{
			From:     from,
			To:       to,
			BumpType: bumpType,
		},
		Metadata: ChangelogMetadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Source:    source,
			Commits:   commits,
		},
	}
	if field != "" {
		data := map[string]interface{}{
			"version":  out.Version,
			"metadata": out.Metadata,
		}
		value, ok := data[field]
		if !ok {
			return fmt.Errorf("unknown field %q", field)
		}
		return writeSelected(w, value)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encoding changelog JSON: %w", err)
	}
	return nil
}
