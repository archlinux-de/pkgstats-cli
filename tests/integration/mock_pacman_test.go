package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pkgstats-cli/internal/system"
)

func createPacmanConf(t *testing.T, mirror string) string {
	t.Helper()

	s := system.NewSystem()
	osArchitecture, err := s.GetArchitecture()
	if err != nil {
		t.Fatal(err)
	}
	dbPath := createPacmanDBPath(t)

	if mirror == "" {
		mirror = "https://geo.mirror.pkgbuild.com/core/os/"
	}

	pacmanConfFile := filepath.Join(t.TempDir(), "pacman.conf")
	if err := os.WriteFile(pacmanConfFile, fmt.Appendf(nil, "[options]\nDBPath=%s\n[core]\nServer=%s%s", dbPath, mirror, osArchitecture), 0o600); err != nil {
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
		"secret-package-1.0.0-1",
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

func createPkgstatsConf(t *testing.T, packages []string, mirrors []string) string {
	t.Helper()

	var content strings.Builder
	content.WriteString("blocklist:\n")
	content.WriteString("  packages:\n")
	for _, pkg := range packages {
		content.WriteString(fmt.Sprintf("    - \"%s\"\n", pkg))
	}
	content.WriteString("  mirrors:\n")
	for _, mirror := range mirrors {
		content.WriteString(fmt.Sprintf("    - \"%s\"\n", mirror))
	}

	pkgstatsConfFile := filepath.Join(t.TempDir(), "pkgstats.yaml")
	if err := os.WriteFile(pkgstatsConfFile, []byte(content.String()), 0o600); err != nil {
		t.Fatalf("Failed to create pkgstats.conf: %v", err)
	}

	return pkgstatsConfFile
}
