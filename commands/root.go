package commands

import (
	"errors"
	"fmt"
	"github.com/belitre/gotpl/tpl"
	"github.com/spf13/cobra"
	"os"
	"strings"
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

var valuesFile string

var supportedExtensionFiles = []string{".yaml", ".yml", ".json"}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("incorrect number of arguments")
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		return fmt.Errorf("file %s not found", args[0])
	}
	if _, err := os.Stat(valuesFile); os.IsNotExist(err) {
		return fmt.Errorf("file %s not found", valuesFile)
	}
	if args[0] == valuesFile {
		return errors.New("gotpl \"inception\" is not allowed")
	}
	isSupported := false
	for _, v := range supportedExtensionFiles {
		if strings.HasSuffix(valuesFile, v) {
			isSupported = true
		}
	}
	if !isSupported {
		return fmt.Errorf("file %s has invalid extension, supported extensions for values files are: %s", valuesFile, strings.Join(supportedExtensionFiles, ", "))
	}
	return nil
}

func runCommand(cmd *cobra.Command, args []string) {
	if err := tpl.ParseTemplate(args[0], valuesFile); err != nil {
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
	rootCmd.PersistentFlags().StringVarP(&valuesFile, "file", "f", "", fmt.Sprintf("values file, supports json and yaml files with extensions: %s", strings.Join(supportedExtensionFiles, ", ")))
	rootCmd.MarkPersistentFlagRequired("file")
}
