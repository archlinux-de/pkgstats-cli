package cmd

import (
	"fmt"

	"pkgstats-cli/internal/api/request"

	"github.com/spf13/cobra"
)

const (
	minLimit = 1
	maxLimit = 10000
)

var limit = 10

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search packages and list their popularity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if limit < minLimit || limit > maxLimit {
			return fmt.Errorf("valid limit needs to be between %d and %d", minLimit, maxLimit)
		}
		if !request.ValidatePackageName(args[0]) {
			return fmt.Errorf("invalid package name %s", args[0])
		}

		client := request.NewClient(baseURL)

		ppl, err := client.SearchPackages(args[0], limit)
		if err != nil {
			return err
		}

		request.PrintPackagePopularities(cmd.OutOrStdout(), ppl)
		fmt.Fprintln(cmd.OutOrStdout())
		request.PrintSearchURL(cmd.OutOrStdout(), baseURL, args[0])

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&limit, "limit", "l", limit, fmt.Sprintf("Limit the results from %d to %d entries", minLimit, maxLimit))
}
