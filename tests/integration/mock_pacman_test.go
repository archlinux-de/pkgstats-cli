package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"pkgstats-cli/internal/system"
)

func createPacmanConf(t *testing.T) string {
	t.Helper()

	s := system.NewSystem()
	osArchitecture, err := s.GetArchitecture()
	if err != nil {
		t.Fatal(err)
	}
	dbPath := createPacmanDBPath(t)

	pacmanConfFile := filepath.Join(t.TempDir(), "pacman.conf")
	if err := os.WriteFile(pacmanConfFile, fmt.Appendf(nil, "[options]\nDBPath=%s\n[core]\nServer=https://geo.mirror.pkgbuild.com/core/os/%s", dbPath, osArchitecture), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}

	return pacmanConfFile
}

func createPacmanDBPath(t *testing.T) string {
	t.Helper()

	dbPath := t.TempDir()
	localDir := filepath.Join(dbPath, "local")
	if err := os.Mkdir(localDir, 0o700); err != nil {
		t.Fatalf("Failed to create local directory: %v", err)
	}

	subdirs := []string{
		"pacman-1.0-1",
		"pacman-mirrorlist-2.0-2",
	}

	for _, dir := range subdirs {
		subdirPath := filepath.Join(localDir, dir)
		err := os.Mkdir(subdirPath, 0o700)
		if err != nil {
			t.Fatalf("Failed to create subdirectory %s: %v", dir, err)
		}
	}

	return dbPath
}
