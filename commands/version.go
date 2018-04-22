package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// TODO: makefile with cross compile and set version
// version, should be updated each time there is a new release
var version = "v0.2-alpha"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gotpl",
	Long:  `Print the version number of gotpl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gotpl %s\n", version)
	},
}
