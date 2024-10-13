package integration_test

import (
	"bytes"
	"net/http/httptest"
	"os"
	"testing"

	"pkgstats-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func pkgstats(t *testing.T, args ...string) (string, error) {
	t.Helper()

	rootCmd := cmd.GetRootCmd()
	defer resetFlags(t, rootCmd)

	server := httptest.NewServer(NewServer())
	defer server.Close()

	args = append(args, "--base-url", server.URL)

	if os.Getenv("INTEGRATION_TEST") != "1" {
		args = append(args, "--pacman-conf", createPacmanConf(t))
	} else {
		t.Log("Using actual pacman in integration test")
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		return "", err
	}
	return buf.String(), nil
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
