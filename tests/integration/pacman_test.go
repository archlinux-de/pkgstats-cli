package integration_test

import (
	"net/url"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"pkgstats-cli/internal/pacman"
)

func TestGetInstalledPackagesMatchesPacman(t *testing.T) {
	if _, err := exec.LookPath("pacman"); err != nil {
		const pacmanNotFoundMessage = "test requires pacman to be installed"
		if os.Getenv("INTEGRATION_TEST") == "1" {
			t.Fatal(pacmanNotFoundMessage)
		} else {
			t.Skip(pacmanNotFoundMessage)
		}
	}

	p, err := pacman.Parse("/etc/pacman.conf")
	if err != nil {
		t.Fatal(err)
	}
	parsedPackages, err := p.GetInstalledPackages()
	if err != nil {
		t.Fatal(err, parsedPackages)
	}

	cmd := exec.Command("pacman", "-Qq")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err, out)
	}
	pacmanPackages := strings.Split(strings.TrimSpace(string(out)), "\n")

	slices.Sort(parsedPackages)
	slices.Sort(pacmanPackages)

	if slices.Compare(parsedPackages, pacmanPackages) != 0 {
		t.Fatalf("pacman and pkgstats disagree")
	}
}

func TestGetServerMatchesPacmanConf(t *testing.T) {
	if _, err := exec.LookPath("pacman-conf"); err != nil {
		const pacmanNotFoundMessage = "test requires pacman-conf to be installed"
		if os.Getenv("INTEGRATION_TEST") == "1" {
			t.Fatal(pacmanNotFoundMessage)
		} else {
			t.Skip(pacmanNotFoundMessage)
		}
	}

	p, err := pacman.Parse("/etc/pacman.conf")
	if err != nil {
		t.Fatal(err)
	}
	parsedServer, err := p.GetServer()
	if err != nil {
		t.Fatal(err, parsedServer)
	}
	parsedServerUrl, err := url.Parse(parsedServer)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("pacman-conf", "--repo=core", "Server")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err, out)
	}
	pacmanConfServer := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	pacmanConfServerUrl, err := url.Parse(parsedServer)
	if err != nil {
		t.Fatal(err)
	}

	if parsedServerUrl.Redacted() != pacmanConfServerUrl.Redacted() {
		t.Fatalf("pacman-conf and pkgstats disagree: %s vs %s", pacmanConfServer, parsedServer)
	}
}
