package cmd

import (
	"fmt"

	"pkgstats-cli/internal/api/request"

	"github.com/spf13/cobra"
)

const (
	minPackages = 1
	maxPackages = 20
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show one or more packages and compare their popularity",
	Args:  cobra.RangeArgs(minPackages, maxPackages),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := request.ValidatePackageNames(args); err != nil {
			return err
		}

		showFunc := func(client *request.Client, args []string) (request.PackagePopularityList, []error) {
			return client.GetPackages(args...)
		}

		if err := handleRequest(cmd.OutOrStdout(), args, showFunc); err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout())
		request.PrintShowURL(cmd.OutOrStdout(), baseURL, args)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
