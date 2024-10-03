package integration_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"pkgstats-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func pkgstats(t *testing.T, args ...string) (output string, err error) {
	t.Helper()

	requiresPacman(t)

	rootCmd := cmd.GetRootCmd()
	defer resetFlags(t, rootCmd)

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

func resetFlags(t *testing.T, cmd *cobra.Command) {
	t.Helper()

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if err := flag.Value.Set(flag.DefValue); err != nil {
			t.Fatal(err)
		}
	})
	for _, subCmd := range cmd.Commands() {
		resetFlags(t, subCmd)
	}
}
