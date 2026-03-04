package cmd

import (
	"os"

	"pkgstats-cli/internal/build"

	"github.com/spf13/cobra"
)

const (
	baseUrlParam      = "base-url"
	pacmanConfParam   = "pacman-conf"
	pkgstatsConfParam = "pkgstats-conf"
)

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	var (
		baseURL      = "https://pkgstats.archlinux.de"
		pacmanConf   = "/etc/pacman.conf"
		pkgstatsConf = ""
	)

	rootCmd := &cobra.Command{
		Use:          "pkgstats",
		Short:        "pkgstats client",
		Version:      build.Version,
		SilenceUsage: true,
	}

	rootCmd.PersistentFlags().StringVar(&baseURL, baseUrlParam, baseURL, "base url of the pkgstats server")
	_ = rootCmd.PersistentFlags().MarkHidden(baseUrlParam)

	rootCmd.PersistentFlags().StringVar(&pacmanConf, pacmanConfParam, pacmanConf, "path to pacman.conf")
	_ = rootCmd.PersistentFlags().MarkHidden(pacmanConfParam)

	rootCmd.PersistentFlags().StringVarP(&pkgstatsConf, pkgstatsConfParam, "c", pkgstatsConf, "path to pkgstats config file")
	_ = rootCmd.PersistentFlags().MarkHidden(pkgstatsConfParam)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(newSubmitCmd(&baseURL, &pacmanConf, &pkgstatsConf))
	rootCmd.AddCommand(newSearchCmd(&baseURL))
	rootCmd.AddCommand(newShowCmd(&baseURL))
	rootCmd.AddCommand(newArchitectureCmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}
