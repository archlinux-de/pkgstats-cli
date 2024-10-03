package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"pkgstats-cli/internal/system"
)

func requiresPacman(t *testing.T) {
	t.Helper()

	checkAndMockCommand(t, "pacman", "pacman", "pacman-mirrorlist")
	checkAndMockCommand(t, "pacman-conf", func() string {
		s := system.NewSystem()
		osArchitecture, err := s.GetArchitecture()
		if err != nil {
			t.Fatal(err)
		}
		return fmt.Sprintf("https://geo.mirror.pkgbuild.com/core/os/%s", osArchitecture)
	}())
}

func checkAndMockCommand(t *testing.T, name string, output ...string) {
	t.Helper()

	if _, err := exec.LookPath(name); err != nil {
		if _, err := exec.LookPath("/bin/sh"); err != nil {
			t.Skipf("%s was not found and could not be mocked", name)
		}

		t.Logf("Using mocked version of %s", name)
		mockCommand(t, name, output...)
	}
}

func mockCommand(t *testing.T, name string, output ...string) {
	t.Helper()
	tempDir := t.TempDir()
	mockPath := filepath.Join(tempDir, name)

	file, err := os.Create(mockPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "#!/bin/sh\necho '%s'", strings.Join(output, "\n")); err != nil {
		t.Fatal(err)
	}

	if err := os.Chmod(mockPath, 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", fmt.Sprintf("%s%c%s", tempDir, os.PathListSeparator, os.Getenv("PATH")))
}
