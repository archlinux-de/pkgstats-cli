package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

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
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", baseURL, "base url of the pkgstats server")
	rootCmd.PersistentFlags().MarkHidden("base-url")
}
