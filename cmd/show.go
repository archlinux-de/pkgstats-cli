package cmd

import (
	"fmt"
	"os"
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
		client := request.NewClient(baseURL)

		ppl, err := client.GetPackages(args...)
		if err != nil {
			return err
		}

		request.PrintPackagePopularities(os.Stdout, ppl)
		fmt.Println()
		request.PrintShowURL(os.Stdout, baseURL, args)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
