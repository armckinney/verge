package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionParseCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantOut   string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "Standard Semver (defaults to json)",
			args:    []string{"1.2.3"},
			wantOut: `"version_type": "semver"`,
		},
		{
			name:    "Standard Semver with explicit format text",
			args:    []string{"1.2.3", "--format", "text"},
			wantOut: "1.2.3\n",
		},
		{
			name:    "Standard Semver with global --json flag",
			args:    []string{"1.2.3", "--json"},
			wantOut: `"version_type": "semver"`,
		},
		{
			name:    "VSemver with global field",
			args:    []string{"v1.2.3-dev.4", "--field", "major"},
			wantOut: "1\n",
		},
		{
			name:    "VSemver with nested floating.major",
			args:    []string{"v1.2.3-dev.4", "--field", "floating.major", "--format", "text"},
			wantOut: "v1\n",
		},
		{
			name:    "VSemver with nested floating.minor",
			args:    []string{"v1.2.3-dev.4", "--field", "floating.minor", "--format", "text"},
			wantOut: "v1.2\n",
		},
		{
			name:    "VSemver with nested floating.prerelease",
			args:    []string{"v1.2.3-dev.4", "--field", "floating.prerelease", "--format", "text"},
			wantOut: "v1.2.3-dev\n",
		},
		{
			name:    "Stable vsemver with nested floating.prerelease",
			args:    []string{"v1.2.3", "--field", "floating.prerelease", "--format", "text"},
			wantOut: "v1.2.3-dev\n",
		},
		{
			name:    "Stable pep440 with nested floating.prerelease",
			args:    []string{"1.2.3", "--type", "pep440", "--field", "floating.prerelease", "--format", "text"},
			wantOut: "1.2.3dev\n",
		},
		{
			name:    "PEP440 with json format",
			args:    []string{"1.2.3a5", "--format", "json"},
			wantOut: `"major": 1`,
		},
		{
			name:      "Invalid version format",
			args:      []string{"invalid-version"},
			wantErr:   true,
			errSubstr: `unable to parse version "invalid-version" with any supported parser`,
		},
		{
			name:      "Explicit type invalid",
			args:      []string{"1.2.3", "--type", "invalid"},
			wantErr:   true,
			errSubstr: `invalid version type: "invalid"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset globalFlags before each run
			globalFlags.field = ""
			globalFlags.format = "text"
			globalFlags.json = false

			root := &cobra.Command{Use: "verge"}
			root.PersistentFlags().StringVarP(&globalFlags.format, "format", "f", "text", "")
			root.PersistentFlags().BoolVar(&globalFlags.json, "json", false, "")
			root.PersistentFlags().StringVar(&globalFlags.field, "field", "", "")

			cmd := versionParseCmd()
			root.AddCommand(cmd)

			// Setup arguments for the root command execution
			args := append([]string{"parse"}, tt.args...)
			root.SetArgs(args)

			// Capture os.Stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := root.Execute()

			w.Close()
			var outBuf bytes.Buffer
			_, _ = io.Copy(&outBuf, r)
			os.Stdout = oldStdout

			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() err = %v, wantErr = %v", err, tt.wantErr)
			}

			if err != nil {
				if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("expected err containing %q, got %q", tt.errSubstr, err.Error())
				}
			} else {
				got := outBuf.String()
				if !strings.Contains(got, tt.wantOut) {
					t.Errorf("expected output containing %q, got %q", tt.wantOut, got)
				}
			}
		})
	}
}
