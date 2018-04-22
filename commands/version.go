package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// TODO: makefile with cross compile and set version
// Version should be updated each time there is a new release
var (
	Version      = "v0.2-alpha"
	GitCommit    = ""
	GitTreeState = ""
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
		if len(GitTreeState) > 0 {
			fmt.Printf("git tree state: %s\n", GitCommit)
		}
	},
}
