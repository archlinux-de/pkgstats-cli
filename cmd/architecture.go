package cmd

import (
	"fmt"

	"pkgstats-cli/internal/system"

	"github.com/spf13/cobra"
)

var architectureCmd = &cobra.Command{
	Use:     "architecture",
	Aliases: []string{"arch"},
	Short:   "Shows information about CPU and OS architecture",
	Hidden:  true,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		system := system.NewSystem()
		osArchitecture, _ := system.GetArchitecture()
		systemArchitecture, _ := system.GetCpuArchitecture()
		fmt.Fprintf(cmd.OutOrStdout(), "You are using a %s CPU on a %s OS\n", systemArchitecture, osArchitecture)
	},
}

var osArchitectureCommand = &cobra.Command{
	Use:   "os",
	Short: "Shows OS architecture",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		system := system.NewSystem()
		osArchitecture, _ := system.GetArchitecture()
		fmt.Fprintln(cmd.OutOrStdout(), osArchitecture)
	},
}

var systemArchitectureCommand = &cobra.Command{
	Use:     "system",
	Aliases: []string{"cpu"},
	Short:   "Shows CPU architecture",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		system := system.NewSystem()
		systemArchitecture, _ := system.GetCpuArchitecture()
		fmt.Fprintln(cmd.OutOrStdout(), systemArchitecture)
	},
}

func init() {
	rootCmd.AddCommand(architectureCmd)
	architectureCmd.AddCommand(osArchitectureCommand)
	architectureCmd.AddCommand(systemArchitectureCommand)
}
