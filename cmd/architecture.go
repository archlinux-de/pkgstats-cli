package cmd

import (
	"fmt"

	"pkgstats-cli/internal/system"

	"github.com/spf13/cobra"
)

func newArchitectureCmd() *cobra.Command {
	architectureCmd := &cobra.Command{
		Use:     "architecture",
		Aliases: []string{"arch"},
		Short:   "Shows information about CPU and OS architecture",
		Hidden:  true,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := system.NewSystem()

			osArchitecture, err := s.GetArchitecture()
			if err != nil {
				return fmt.Errorf("could not get OS architecture: %w", err)
			}

			systemArchitecture, err := s.GetCpuArchitecture()
			if err != nil {
				return fmt.Errorf("could not get CPU architecture: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "You are using a %s CPU on a %s OS\n", systemArchitecture, osArchitecture)

			return nil
		},
	}

	architectureCmd.AddCommand(newOsArchitectureCmd())
	architectureCmd.AddCommand(newSystemArchitectureCmd())

	return architectureCmd
}

func newOsArchitectureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "os",
		Short: "Shows OS architecture",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := system.NewSystem()

			osArchitecture, err := s.GetArchitecture()
			if err != nil {
				return fmt.Errorf("could not get OS architecture: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), osArchitecture)

			return nil
		},
	}
}

func newSystemArchitectureCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "system",
		Aliases: []string{"cpu"},
		Short:   "Shows CPU architecture",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := system.NewSystem()

			systemArchitecture, err := s.GetCpuArchitecture()
			if err != nil {
				return fmt.Errorf("could not get CPU architecture: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), systemArchitecture)

			return nil
		},
	}
}
