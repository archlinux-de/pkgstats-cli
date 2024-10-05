package cmd

import (
	"fmt"
	"os"

	"pkgstats-cli/internal/build"

	"github.com/spf13/cobra"
)

const baseUrlParam = "base-url"

var (
	baseURL = "https://pkgstats.archlinux.de"
	rootCmd = &cobra.Command{
		Use:          "pkgstats",
		Short:        "pkgstats client",
		Version:      build.Version,
		SilenceUsage: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, baseUrlParam, baseURL, "base url of the pkgstats server")
	if err := rootCmd.PersistentFlags().MarkHidden(baseUrlParam); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		os.Exit(1)
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
