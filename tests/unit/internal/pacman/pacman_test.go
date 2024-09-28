//go:build amd64 || 386

package pacman_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"pkgstats-cli/internal/pacman"
)

func init() {
	Mocks["pacman"] = func() {
		fmt.Println("pacman")
		fmt.Println("linux")
		os.Exit(0)
	}

	Mocks["pacman-conf"] = func() {
		fmt.Println("https://mirror.rackspace.com/archlinux/core/os/x86_64")
		fmt.Println("https://geo.mirror.pkgbuild.com/core/os/x86_64")
		os.Exit(0)
	}
}

var Mocks = make(map[string]func())

func TestMain(m *testing.M) {
	mockName := os.Getenv("TEST_MOCK")
	if mockName != "" {
		mock, ok := Mocks[mockName]
		if ok {
			mock()
		}
	}

	os.Exit(m.Run())
}

func TestGetInstalledPackages(t *testing.T) {
	pacman := pacman.Pacman{
		Pacman:  os.Args[0],
		Timeout: 1 * time.Second,
		Env:     []string{"TEST_MOCK=pacman"},
	}

	out, err := pacman.GetInstalledPackages()
	if err != nil {
		t.Error(err, out)
	}
	if strings.Join(out, ",") != "pacman,linux" {
		t.Error(out)
	}
}

func TestGetServer(t *testing.T) {
	pacman := pacman.Pacman{
		PacmanConf: os.Args[0],
		Timeout:    1 * time.Second,
		Repository: "core",
		Env:        []string{"TEST_MOCK=pacman-conf"},
	}

	out, err := pacman.GetServer()
	if err != nil {
		t.Error(err, out)
	}
	if out != "https://mirror.rackspace.com/archlinux/" {
		t.Error(out)
	}
}