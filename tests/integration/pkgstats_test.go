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

type pkgstatsOptions struct {
	mirror          string
	pkgBlocklist    []string
	mirrorBlocklist []string
}

type PkgstatsOption func(*pkgstatsOptions)

func WithMirror(mirror string) PkgstatsOption {
	return func(o *pkgstatsOptions) {
		o.mirror = mirror
	}
}

func WithPkgBlocklist(list []string) PkgstatsOption {
	return func(o *pkgstatsOptions) {
		o.pkgBlocklist = list
	}
}

func WithMirrorBlocklist(list []string) PkgstatsOption {
	return func(o *pkgstatsOptions) {
		o.mirrorBlocklist = list
	}
}

func pkgstats(t *testing.T, args []string, opts ...PkgstatsOption) (string, error) {
	t.Helper()

	options := &pkgstatsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	rootCmd := cmd.GetRootCmd()
	defer resetFlags(t, rootCmd)

	server := httptest.NewServer(NewServer())
	defer server.Close()

	args = append(args, "--base-url", server.URL)

	if os.Getenv("INTEGRATION_TEST") == "1" {
		t.Log("Using actual pacman in integration test")
	} else {
		args = append(args, "--pacman-conf", createPacmanConf(t, options.mirror))
	}

	args = append(args, "--pkgstats-conf", createPkgstatsConf(t, options.pkgBlocklist, options.mirrorBlocklist))

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
