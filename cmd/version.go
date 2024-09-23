package cmd

import (
	"fmt"
	"pkgstats-cli/internal/build"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the pkgstats client version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "pkgstats, version", build.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
