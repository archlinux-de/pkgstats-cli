//go:build amd64 || 386

package pacman_test

import (
	"errors"
	"strings"
	"testing"

	"pkgstats-cli/internal/pacman"
)

type mockCommandExecutor struct {
	Output map[string]string
	Err    error
}

func (m mockCommandExecutor) Execute(name string, arg ...string) ([]byte, error) {
	key := name + " " + strings.Join(arg, " ")
	if output, ok := m.Output[key]; ok {
		return []byte(output), m.Err
	}
	return nil, errors.New("command not found")
}

func TestGetInstalledPackages(t *testing.T) {
	pacman := pacman.Pacman{Executor: mockCommandExecutor{
		Output: map[string]string{
			"pacman -Qq": "pacman\nlinux",
		},
	}}

	out, err := pacman.GetInstalledPackages()
	if err != nil {
		t.Error(err, out)
	}
	if strings.Join(out, ",") != "pacman,linux" {
		t.Error(out)
	}
}

func TestGetServer(t *testing.T) {
	pacman := pacman.Pacman{Executor: mockCommandExecutor{
		Output: map[string]string{
			"pacman-conf --repo core Server": "https://mirror.rackspace.com/archlinux/core/os/x86_64\nhttps://geo.mirror.pkgbuild.com/core/os/x86_64",
		},
	}}

	out, err := pacman.GetServer()
	if err != nil {
		t.Error(err, out)
	}
	if out != "https://mirror.rackspace.com/archlinux/" {
		t.Error(out)
	}
}
