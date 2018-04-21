package commands

import (
	"errors"
	"fmt"
	"github.com/belitre/gotpl/tpl"
	"github.com/spf13/cobra"
	"os"
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
	if len(args) != 1 {
		return errors.New("incorrect number of arguments")
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		return fmt.Errorf("file %s not found", args[0])
	}
	return nil
}

func runCommand(cmd *cobra.Command, args []string) {
	if err := tpl.ParseTemplate(args[0], valueFiles, setValues, isStrict); err != nil {
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
	rootCmd.PersistentFlags().StringArrayVarP(&valueFiles, "values", "f", []string{}, "specify values in a YAML or JSON files")
	rootCmd.MarkPersistentFlagRequired("file")
	rootCmd.PersistentFlags().StringArrayVarP(&setValues, "set", "s", []string{}, "<key>=<value> pairs (take precedence over values in --values files)")
	rootCmd.PersistentFlags().BoolVarP(&isStrict, "strict", "", false, "If strict is enabled, template rendering will fail if a template references a value that was not passed in")
}
