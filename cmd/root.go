package cmd

import (
	"fmt"
	"os"

	"pkgstats-cli/internal/build"

	"github.com/spf13/cobra"
)

const (
	baseUrlParam      = "base-url"
	pacmanConfParam   = "pacman-conf"
	pkgstatsConfParam = "pkgstats-conf"
)

var (
	baseURL      = "https://pkgstats.archlinux.de"
	pacmanConf   = "/etc/pacman.conf"
	pkgstatsConf = ""
	rootCmd      = &cobra.Command{
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

	rootCmd.PersistentFlags().StringVar(&pacmanConf, pacmanConfParam, pacmanConf, "path to pacman.conf")
	if err := rootCmd.PersistentFlags().MarkHidden(pacmanConfParam); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVar(&pkgstatsConf, pkgstatsConfParam, pkgstatsConf, "path to pkgstats config file")
	if err := rootCmd.PersistentFlags().MarkHidden(pkgstatsConfParam); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		os.Exit(1)
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
