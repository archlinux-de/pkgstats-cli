package integration_test

import (
	"bytes"
	"net/http/httptest"

	"pkgstats-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func pkgstats(args ...string) (output string, err error) {
	rootCmd := cmd.GetRootCmd()
	defer resetFlags(rootCmd)

	server := httptest.NewServer(NewServer())
	defer server.Close()

	args = append(args, "--base-url", server.URL)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	err = rootCmd.Execute()
	return buf.String(), err
}

func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
	})
	for _, subCmd := range cmd.Commands() {
		resetFlags(subCmd)
	}
}
