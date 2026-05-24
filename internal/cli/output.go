package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
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
	Field  string
	Writer io.Writer
}

// NewOutput creates a new Output writer.
func NewOutput(format OutputFormat) *Output {
	return &Output{Format: format, Writer: os.Stdout}
}

// Print prints key-value pairs.
func (o *Output) Print(data map[string]interface{}) error {
	if o.Field != "" {
		value, ok := data[o.Field]
		if !ok {
			return fmt.Errorf("unknown field %q", o.Field)
		}
		return writeSelected(o.Writer, value)
	}

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
func (o *Output) PrintValue(v interface{}) error {
	if o.Field != "" && o.Field != "value" {
		return fmt.Errorf("unknown field %q", o.Field)
	}

	if o.Field != "" {
		return writeSelected(o.Writer, v)
	}

	switch o.Format {
	case FormatJSON:
		data := map[string]interface{}{"value": v}
		enc := json.NewEncoder(o.Writer)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	default:
		fmt.Fprintln(o.Writer, strings.TrimSpace(fmt.Sprintf("%v", v)))
		return nil
	}
}

func writeSelected(w io.Writer, value interface{}) error {
	if value == nil {
		_, err := fmt.Fprintln(w)
		return err
	}

	switch selected := value.(type) {
	case string:
		_, err := fmt.Fprintln(w, selected)
		return err
	case []byte:
		_, err := fmt.Fprintln(w, string(selected))
		return err
	case fmt.Stringer:
		_, err := fmt.Fprintln(w, selected.String())
		return err
	}

	rv := reflect.ValueOf(value)
	if rv.IsValid() {
		switch rv.Kind() {
		case reflect.Bool:
			_, err := fmt.Fprintln(w, rv.Bool())
			return err
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_, err := fmt.Fprintln(w, rv.Int())
			return err
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			_, err := fmt.Fprintln(w, rv.Uint())
			return err
		case reflect.Float32, reflect.Float64:
			_, err := fmt.Fprintln(w, rv.Float())
			return err
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			return enc.Encode(value)
		}
	}

	_, err := fmt.Fprintln(w, value)
	return err
}
