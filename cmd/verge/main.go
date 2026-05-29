package main

import (
	"errors"
	"os"

	"example.com/verge/internal/cli"
)

// Version info injected by goreleaser ldflags
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	cli.SetVersionInfo(Version, Commit, Date)
	if err := cli.Execute(); err != nil {
		var cliErr *cli.CLIError
		if errors.As(err, &cliErr) {
			os.Exit(cliErr.Code)
		}
		os.Exit(1)
	}
}
