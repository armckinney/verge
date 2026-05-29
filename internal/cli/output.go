package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

type Output struct {
	Format OutputFormat
	Field  string
	Writer io.Writer
}

func NewOutput(format OutputFormat) *Output {
	return &Output{Format: format, Writer: os.Stdout}
}

// Print outputs data strictly based on spec 009.
func (o *Output) Print(data map[string]interface{}) error {
	switch o.Format {
	case FormatJSON:
		enc := json.NewEncoder(o.Writer)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	default: // text
		// In text mode, we exclusively output the exact rendered version string bare.
		if rendered, ok := data["rendered"]; ok {
			fmt.Fprintln(o.Writer, rendered)
		} else if versionStr, ok := data["version"]; ok {
			fmt.Fprintln(o.Writer, versionStr)
		}
		return nil
	}
}

func (o *Output) PrintErr(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
