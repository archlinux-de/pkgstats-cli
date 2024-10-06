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
		{"http://user:password@example.com", "http://user:xxxxx@example.com/", false},
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
