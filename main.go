package main

import (
	"os"

	"github.com/airfocusio/kustomization-generator/cmd"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	if err := cmd.Execute(cmd.FullVersion{Version: version, Commit: commit, Date: date, BuiltBy: builtBy}); err != nil {
		os.Exit(1)
	}
}
