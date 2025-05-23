package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/pacman"
	"pkgstats-cli/internal/system"

	"github.com/spf13/cobra"
)

var (
	dumpJSON = false
	quiet    = false
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a list of your installed packages to the pkgstats project",
	Long:  "Submit a list of your installed packages, your system architecture\nand the mirror you are using to the pkgstats project.\n\nStatistics are available at " + baseURL,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dumpJSON && quiet {
			return errors.New("--quiet and --dump-json cannot be used at the same time")
		}

		if !dumpJSON && !quiet {
			fmt.Fprintln(cmd.OutOrStdout(), "Collecting data...")
		}

		p, err := pacman.Parse(pacmanConf)
		if err != nil {
			return err
		}
		request, err := submit.CreateRequest(p, system.NewSystem())
		if err != nil {
			return err
		}

		if !dumpJSON {
			if !quiet {
				fmt.Fprintln(cmd.OutOrStdout(), "Submitting data...")
			}
			client := submit.NewClient(baseURL)

			err := client.SendRequest(*request)
			if err != nil {
				return err
			}

			if !quiet {
				fmt.Fprintln(cmd.OutOrStdout(), "Data were successfully sent")
			}
		} else {
			formattedRequest, err := json.MarshalIndent(request, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(formattedRequest))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVarP(&dumpJSON, "dump-json", "d", dumpJSON, "Dump information that would be sent as JSON")
	submitCmd.Flags().BoolVarP(&quiet, "quiet", "q", quiet, "Suppress any output except errors")
}
