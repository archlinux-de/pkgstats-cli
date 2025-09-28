package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"

	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/config"
	"pkgstats-cli/internal/filter"
	"pkgstats-cli/internal/pacman"
	"pkgstats-cli/internal/system"

	"github.com/spf13/cobra"
)

var (
	dumpJSON   = false
	quiet      = false
	configHelp = false
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a list of your installed packages to the pkgstats project",
	Long: "Submit a list of your installed packages, your system architecture\nand the mirror you are using to the pkgstats project.\n\n" +
		"Statistics are available at " + baseURL,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configHelp {
			printConfigHelp(cmd.OutOrStdout())
			return nil
		}

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

		req, err := submit.CreateRequest(p, system.NewSystem())
		if err != nil {
			return err
		}

		c, err := config.Load(pkgstatsConf)
		if err != nil {
			return err
		}
		err = filter.FilterRequest(c, req)
		if err != nil {
			return err
		}

		if dumpJSON {
			return dumpRequest(cmd.OutOrStdout(), req)
		} else {
			return submitRequest(cmd.OutOrStdout(), req)
		}
	},
}

func dumpRequest(writer io.Writer, req *submit.Request) error {
	slices.Sort(req.Pacman.Packages)
	formattedRequest, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(writer, string(formattedRequest))
	return nil
}

func submitRequest(writer io.Writer, req *submit.Request) error {
	if !quiet {
		fmt.Fprintln(writer, "Submitting data...")
	}

	client := submit.NewClient(baseURL)

	if err := client.SendRequest(*req); err != nil {
		return err
	}

	if !quiet {
		fmt.Fprintln(writer, "Data were successfully sent")
	}

	return nil
}

func printConfigHelp(writer io.Writer) {
	fmt.Fprintln(writer, `You can configure blocklists for packages and mirrors by creating or editing /etc/pkgstats.yaml.
Packages can be blocked by name, and mirrors by hostname. Both support glob patterns.
You can verify your blocklist configuration by running "pkgstats submit --dump-json".

Example configuration:

blocklist:
  packages:
    - "secret-*"
    - "*-debug"
  mirrors:
    - "mirror.example.com"
    - "*.lan"`)
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVar(&configHelp, "config-help", configHelp, "Show help for configuring blocklists in /etc/pkgstats.yaml")
	submitCmd.Flags().BoolVarP(&dumpJSON, "dump-json", "d", dumpJSON, "Dump information that would be sent as JSON")
	submitCmd.Flags().BoolVarP(&quiet, "quiet", "q", quiet, "Suppress any output except errors")
}
