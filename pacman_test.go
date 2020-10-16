package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func init() {
	Mocks["pacman"] = func() {
		fmt.Println("pacman")
		fmt.Println("linux")
		os.Exit(0)
	}

	Mocks["pacman-conf"] = func() {
		fmt.Println("https://mirror.rackspace.com/archlinux/core/os/x86_64")
		fmt.Println("https://mirror.pkgbuild.com/core/os/x86_64")
		os.Exit(0)
	}
}

func TestGetInstalledPackages(t *testing.T) {
	pacman := Pacman{}
	pacman.pacman = os.Args[0]
	pacman.timeout = 1 * time.Second
	pacman.env = []string{"TEST_MOCK=pacman"}

	out, err := pacman.GetInstalledPackages()

	if err != nil {
		t.Error(err, out)
	}
	if out != "pacman\nlinux" {
		t.Error(out)
	}
}

func TestGetServer(t *testing.T) {
	pacman := Pacman{}
	pacman.pacmanConf = os.Args[0]
	pacman.timeout = 1 * time.Second
	pacman.repository = "core"
	pacman.env = []string{"TEST_MOCK=pacman-conf"}

	out, err := pacman.GetServer()

	if err != nil {
		t.Error(err, out)
	}
	if out != "https://mirror.rackspace.com/archlinux/" {
		t.Error(out)
	}
}
