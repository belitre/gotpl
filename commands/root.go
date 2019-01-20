package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/belitre/gotpl/tpl"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotpl template_file",
	Short: "CLI tool for Golang templates",
	Long: `gotpl - CLI tool for Golang templates
https://github.com/belitre/gotpl
			`,
	Args: validateArgs,
	Run:  runCommand,
}

var valueFiles []string
var setValues []string
var isStrict bool

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("incorrect number of arguments")
	}
	for _, f := range args {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return fmt.Errorf("file %s not found", f)
		}
	}
	return nil
}

func runCommand(cmd *cobra.Command, args []string) {
	if err := tpl.ParseTemplate(args, valueFiles, setValues, isStrict); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArrayVarP(&valueFiles, "values", "f", []string{}, "specify values in a YAML or JSON files")
	rootCmd.Flags().StringArrayVarP(&setValues, "set", "s", []string{}, "<key>=<value> pairs (take precedence over values in --values files)")
	rootCmd.Flags().BoolVarP(&isStrict, "strict", "", false, "If strict is enabled, template rendering will fail if a template references a value that was not passed in")
}
