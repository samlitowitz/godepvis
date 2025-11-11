package main

import (
	"github.com/samlitowitz/godepvis/cmd/godepvis/cmd"
	"os"
)

var (
	Build   string
	Commit  string
	Version string
)

func main() {
	// Setup commands
	rootCmd := cmd.Root()
	versionCmd := cmd.Version(Build, Commit, Version)

	rootCmd.AddCommand(versionCmd)

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
