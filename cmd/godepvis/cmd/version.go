package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Version(build, commit, version string) *cobra.Command {
	return &cobra.Command{
		Use:           "version",
		Short:         "Print version information and quit",
		Long:          "Print version information and quit",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(*cobra.Command, []string) {
			fmt.Println("godepvis")
			fmt.Printf(" Version:       %s\n", version)
			fmt.Printf(" Git commit:    %s\n", commit)
			fmt.Printf(" Build:         %s\n", build)
		},
	}
}
