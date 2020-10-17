package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"pkgstats-cli/internal/api/submit"

	"github.com/spf13/cobra"
)

var (
	dumpJSON = false
	quiet    = false
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a list of your installed packages to the pkgstats project",
	Long:  fmt.Sprintf("Submit a list of your installed packages, your system architecture\nand the mirror you are using to the pkgstats project.\n\nStatistics are available at %s", baseURL),
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dumpJSON && quiet {
			return errors.New("--quiet and --dump-json cannot be used at the same time")
		}

		if !dumpJSON && !quiet {
			fmt.Println("Collecting data...")
		}

		request, err := submit.CreateRequest()
		if err != nil {
			return err
		}

		if !dumpJSON {
			if !quiet {
				fmt.Println("Submitting data...")
			}
			client := submit.NewClient(baseURL)

			err := client.SendRequest(*request)
			if err != nil {
				return err
			}

			if !quiet {
				fmt.Println("Data were successfully sent")
			}
		} else {
			formatedRequest, err := json.MarshalIndent(request, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(formatedRequest))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVarP(&dumpJSON, "dump-json", "d", dumpJSON, "Dump information that would be sent as JSON")
	submitCmd.Flags().BoolVarP(&quiet, "quiet", "q", quiet, "Suppress any output except errors")
}
