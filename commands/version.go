package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version should be updated each time there is a new release
var (
	Version   = "v0.5"
	GitCommit = ""
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gotpl",
	Long:  `Print the version number of gotpl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gotpl-%s\n", Version)
		if len(GitCommit) > 0 {
			fmt.Printf("git commit: %s\n", GitCommit)
		}
	},
}
