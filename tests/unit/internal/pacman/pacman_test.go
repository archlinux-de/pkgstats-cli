package pacman_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"pkgstats-cli/internal/pacman"
)

const (
	SERVER = "https://mirror.rackspace.com/archlinux/"
	MIRROR = SERVER + "core/os/x86_64"
)

func createPacmanDb(t *testing.T, packages []string) string {
	t.Helper()

	dbPath := t.TempDir()
	localDir := filepath.Join(dbPath, "local")
	if err := os.Mkdir(localDir, 0o700); err != nil {
		t.Fatalf("Failed to create local directory: %v", err)
	}

	for _, dir := range packages {
		if err := os.Mkdir(filepath.Join(localDir, dir), 0o700); err != nil {
			t.Fatalf("Failed to create subdirectory %s: %v", dir, err)
		}
	}

	return dbPath
}

func createPacmanConf(t *testing.T, dbPath string, servers []string) string {
	t.Helper()

	pacmanConfFile := filepath.Join(t.TempDir(), "pacman.conf")

	var conf strings.Builder
	if dbPath != "" {
		conf.WriteString(fmt.Sprintf("[options]\nDBPath=%s\n", dbPath))
	}

	if len(servers) > 0 {
		conf.WriteString("[core]\n")
		for _, server := range servers {
			conf.WriteString(fmt.Sprintf("Server=%s\n", server))
		}
	}

	if err := os.WriteFile(pacmanConfFile, []byte(conf.String()), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}

	return pacmanConfFile
}

func TestGetInstalledPackages(t *testing.T) {
	dbPath := createPacmanDb(t, []string{"pacman-1.0-1", "go-2.0-2", "php-fpm-8-8.3-32"})
	pacmanConfFile := createPacmanConf(t, dbPath, []string{MIRROR})

	p, err := pacman.Parse(pacmanConfFile)
	if err != nil {
		t.Fatal(err)
	}
	out, err := p.GetInstalledPackages()
	if err != nil {
		t.Fatal(err, out)
	}

	for _, p := range []string{"pacman", "go", "php-fpm-8"} {
		if !slices.Contains(out, p) {
			t.Errorf("could not find package %s in %v", p, out)
		}
	}
}

func TestGetServer(t *testing.T) {
	pacmanConfFile := createPacmanConf(t, "", []string{MIRROR, "https://geo.mirror.pkgbuild.com/core/os/x86_64"})

	p, err := pacman.Parse(pacmanConfFile)
	if err != nil {
		t.Fatal(err)
	}
	out, err := p.GetServer()
	if err != nil {
		t.Fatal(err, out)
	}
	if out != SERVER {
		t.Error(out)
	}
}

func TestPacmanConfIncludes(t *testing.T) {
	tempDir := t.TempDir()
	pacmanConfFile := filepath.Join(tempDir, "pacman.conf")
	pacmanConfFileInclude1 := filepath.Join(tempDir, "pacman-include1.conf")
	pacmanConfFileInclude2 := filepath.Join(tempDir, "pacman-include2.conf")

	if err := os.WriteFile(pacmanConfFile, fmt.Appendf(nil, "[core]\nInclude=%s\n", pacmanConfFileInclude1), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}
	if err := os.WriteFile(pacmanConfFileInclude1, fmt.Appendf(nil, "Include=%s\nInclude=%s\n", pacmanConfFileInclude2, pacmanConfFileInclude2), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}
	if err := os.WriteFile(pacmanConfFileInclude2, fmt.Appendf(nil, "Server=%s\n", MIRROR), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}

	p, err := pacman.Parse(pacmanConfFile)
	if err != nil {
		t.Fatal(err)
	}
	out, err := p.GetServer()
	if err != nil {
		t.Fatal(err, out)
	}
	if out != SERVER {
		t.Error(out)
	}
}

func TestPacmanConfIncludesGlob(t *testing.T) {
	tempDir := t.TempDir()

	pacmanConfFile := filepath.Join(tempDir, "pacman.conf")
	pacmanConfFileInclude1 := filepath.Join(tempDir, "pacman-include1.conf")
	pacmanConfFileInclude2 := filepath.Join(tempDir, "pacman-include2.conf")

	if err := os.WriteFile(pacmanConfFile, fmt.Appendf(nil, "[core]\nInclude=%s/pacman-*.conf\n", tempDir), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}
	if err := os.WriteFile(pacmanConfFileInclude1, []byte(""), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}
	if err := os.WriteFile(pacmanConfFileInclude2, fmt.Appendf(nil, "Server=%s\n", MIRROR), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}

	p, err := pacman.Parse(pacmanConfFile)
	if err != nil {
		t.Fatal(err)
	}
	out, err := p.GetServer()
	if err != nil {
		t.Fatal(err, out)
	}
	if out != SERVER {
		t.Error(out)
	}
}

func TestPacmanConfComments(t *testing.T) {
	pacmanConfFile := filepath.Join(t.TempDir(), "pacman.conf")
	if err := os.WriteFile(pacmanConfFile, fmt.Appendf(nil, "[core]\n#Server=invalid\nServer=%s\n", MIRROR), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}

	p, err := pacman.Parse(pacmanConfFile)
	if err != nil {
		t.Fatal(err)
	}
	s, err := p.GetServer()
	if err != nil || s != SERVER {
		t.Fatal(err, s)
	}
}

func TestPacmanConfEmptySections(t *testing.T) {
	pacmanConfFile := filepath.Join(t.TempDir(), "pacman.conf")
	if err := os.WriteFile(pacmanConfFile, fmt.Appendf(nil, "[]\n[core]\nServer=%s\n", MIRROR), 0o600); err != nil {
		t.Fatalf("Failed to create pacman.conf: %v", err)
	}

	p, err := pacman.Parse(pacmanConfFile)
	if err != nil {
		t.Fatal(err)
	}
	s, err := p.GetServer()
	if err != nil || s != SERVER {
		t.Fatal(err, s)
	}
}

func TestNormalizeMirrorUrl(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		// Valid inputs
		{"http://example.com/path/to/mirror", "http://example.com/", false},
		{"ftp://example.com/path/to/mirror", "ftp://example.com/", false},
		{"rsync://example.com/path/to/mirror", "rsync://example.com/", false},
		{"https://example.com/path/to/mirror", "https://example.com/", false},
		{"http://example.com/path/to/mirror/", "http://example.com/", false},
		{"https://example.com:8080/path/to/mirror", "https://example.com:8080/", false},
		{"ftp://example.com/path/to/mirror", "ftp://example.com/", false},
		{"ftp://example.com/path/to", "ftp://example.com/", false},
		{"ftp://example.com/path/", "ftp://example.com/", false},
		{"ftp://example.com/path", "ftp://example.com/", false},
		{"ftp://example.com/", "ftp://example.com/", false},
		{"ftp://example.com", "ftp://example.com/", false},
		{"http://example.com/path/to/mirror/", "http://example.com/", false},
		{"http://example.com/path/to/mirror/extra", "http://example.com/path/", false},
		{"http://example.com", "http://example.com/", false},
		{"http://user:password@example.com", "http://example.com/", false},
		{"http://user@example.com", "http://example.com/", false},
		{"http://example.com:1234/", "http://example.com:1234/", false},
		{"file:///foo", "file:///", false},
		{"file:///mnt/mirror/core/os/x86_64", "file:///mnt/mirror/", false},

		// Invalid inputs
		{"", "", true},
		{"/", "", true},
		{"://example.com/path/to/mirror", "", true},
		{"http://example.com:invalidport/path/to/mirror", "", true},
	}

	for _, test := range tests {
		result, err := pacman.NormalizeMirrorUrl(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("NormalizeMirrorUrl(%q) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if result != test.expected {
			t.Errorf("NormalizeMirrorUrl(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}
