package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const baseUrlParam = "base-url"

var baseURL = "https://pkgstats.archlinux.de"
var rootCmd = &cobra.Command{
	Use:   "pkgstats",
	Short: "pkgstats client",
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, baseUrlParam, baseURL, "base url of the pkgstats server")
	if err := rootCmd.PersistentFlags().MarkHidden(baseUrlParam); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
