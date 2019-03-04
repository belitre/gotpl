package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/belitre/gotpl/commands/options"
	"github.com/belitre/gotpl/tpl"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotpl template_file1 templates_dir ...",
	Short: "CLI tool for Golang templates",
	Long: `gotpl - CLI tool for Golang templates
https://github.com/belitre/gotpl
			`,
	Args: validateArgs,
	Run:  runCommand,
}

var opt = &options.Options{}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("incorrect number of arguments")
	}
	for _, f := range args {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return fmt.Errorf("file %s not found", f)
		}
	}
	if len(opt.OutputPath) > 0 {
		info, err := os.Stat(opt.OutputPath)
		if err != nil {
			return fmt.Errorf("Path %s not found: %s", opt.OutputPath, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("Path %s is not a directory", opt.OutputPath)
		}
	}
	return nil
}

func runCommand(cmd *cobra.Command, args []string) {
	if err := tpl.ParseTemplate(args, opt); err != nil {
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
	rootCmd.Flags().StringArrayVarP(&opt.ValueFiles, "values", "f", []string{}, "specify values in a YAML or JSON files")
	rootCmd.Flags().StringArrayVarP(&opt.SetValues, "set", "s", []string{}, "<key>=<value> pairs (take precedence over values in --values files)")
	rootCmd.Flags().BoolVarP(&opt.IsStrict, "strict", "", false, "If strict is enabled, template rendering will fail if a template references a value that was not passed in")
	rootCmd.Flags().StringVarP(&opt.OutputPath, "output", "o", "", "If an output path is provided, instead of stdout, gotpl will generate the output as files in the specified path")
}
