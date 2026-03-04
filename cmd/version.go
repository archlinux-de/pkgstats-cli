package cmd

import (
	"fmt"

	"pkgstats-cli/internal/build"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "version",
		Short:  "Show the pkgstats client version",
		Args:   cobra.NoArgs,
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "pkgstats, version", build.Version)
		},
	}
}
