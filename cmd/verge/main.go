package main

import (
	"example.com/verge/internal/cli"
	"os"
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
		os.Exit(1)
	}
}
