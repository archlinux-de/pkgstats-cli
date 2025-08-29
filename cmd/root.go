package cmd

import (
	"fmt"
	"os"

	"pkgstats-cli/internal/build"

	"github.com/spf13/cobra"
)

const (
	baseUrlParam    = "base-url"
	pacmanConfParam = "pacman-conf"
)

var (
	baseURL    = "https://pkgstats.archlinux.de"
	pacmanConf = "/etc/pacman.conf"
	rootCmd    = &cobra.Command{
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
	if err := setupFlags(rootCmd); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		os.Exit(1)
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func setupFlags(cmd *cobra.Command) error {
	cmd.PersistentFlags().StringVar(&baseURL, baseUrlParam, baseURL, "base url of the pkgstats server")
	if err := cmd.PersistentFlags().MarkHidden(baseUrlParam); err != nil {
		return err
	}

	cmd.PersistentFlags().StringVar(&pacmanConf, pacmanConfParam, pacmanConf, "path to pacman.conf")
	if err := cmd.PersistentFlags().MarkHidden(pacmanConfParam); err != nil {
		return err
	}
	return nil
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
