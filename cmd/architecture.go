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
	RunE: func(cmd *cobra.Command, args []string) error {
		s := system.NewSystem()

		osArchitecture, err := s.GetArchitecture()
		if err != nil {
			return fmt.Errorf("could not get OS architecture: %v", err)
		}

		systemArchitecture, err := s.GetCpuArchitecture()
		if err != nil {
			return fmt.Errorf("could not get CPU architecture: %v", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "You are using a %s CPU on a %s OS\n", systemArchitecture, osArchitecture)

		return nil
	},
}

var osArchitectureCommand = &cobra.Command{
	Use:   "os",
	Short: "Shows OS architecture",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := system.NewSystem()

		osArchitecture, err := s.GetArchitecture()
		if err != nil {
			return fmt.Errorf("could not get OS architecture: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), osArchitecture)

		return nil
	},
}

var systemArchitectureCommand = &cobra.Command{
	Use:     "system",
	Aliases: []string{"cpu"},
	Short:   "Shows CPU architecture",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := system.NewSystem()

		systemArchitecture, err := s.GetCpuArchitecture()
		if err != nil {
			return fmt.Errorf("could not get CPU architecture: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), systemArchitecture)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(architectureCmd)
	architectureCmd.AddCommand(osArchitectureCommand)
	architectureCmd.AddCommand(systemArchitectureCommand)
}
