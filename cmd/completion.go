package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:                   "completion",
	Short:                 "Generate completion script",
	DisableFlagsInUseLine: true,
	Hidden:                true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		default:
			err = errors.New("Unknown argument")
		}
		return
	},
}

func init() {
	completionCmd.Use += " [" + strings.Join(completionCmd.ValidArgs, "|") + "]"
	rootCmd.AddCommand(completionCmd)
}
