package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// OutputFormat represents the output format.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

// Output handles formatted output.
type Output struct {
	Format OutputFormat
	Writer io.Writer
}

// NewOutput creates a new Output writer.
func NewOutput(format OutputFormat) *Output {
	return &Output{Format: format, Writer: os.Stdout}
}

// Print prints key-value pairs.
func (o *Output) Print(data map[string]interface{}) error {
	switch o.Format {
	case FormatJSON:
		enc := json.NewEncoder(o.Writer)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	default: // text
		for k, v := range data {
			fmt.Fprintf(o.Writer, "%s\t%v\n", k, v)
		}
		return nil
	}
}

// PrintValue prints a single value.
func (o *Output) PrintValue(v interface{}) {
	switch o.Format {
	case FormatJSON:
		data := map[string]interface{}{"value": v}
		enc := json.NewEncoder(o.Writer)
		enc.SetIndent("", "  ")
		_ = enc.Encode(data)
	default:
		fmt.Fprintln(o.Writer, strings.TrimSpace(fmt.Sprintf("%v", v)))
	}
}
