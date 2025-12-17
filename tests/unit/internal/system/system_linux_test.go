package system_test

import (
	"os"
	"runtime"
	"slices"
	"strings"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetMachine(t *testing.T) {
	arch, err := system.NewSystem().GetArchitecture()

	expectedArch := func(arch string) bool { return arch == runtime.GOARCH }
	switch runtime.GOARCH {
	case "amd64":
		expectedArch = func(arch string) bool { return strings.HasPrefix(arch, system.X86_64) }
	case "386":
		expectedArch = func(arch string) bool { return slices.Contains([]string{system.I586, system.I686}, arch) }
	case "arm":
		expectedArch = func(arch string) bool { return strings.HasPrefix(arch, "armv") }
	case "arm64":
		expectedArch = func(arch string) bool { return arch == system.AARCH64 }
	case "loong64":
		expectedArch = func(arch string) bool { return arch == "loongarch64" }
	}

	if err != nil {
		t.Error(err)
	}
	if !expectedArch(arch) {
		t.Error(arch)
	}
}

func createOsReleaseFile(t *testing.T, content string) string {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "os-release-test-")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Fatalf("failed to remove temp file: %v", err)
		}
	})

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}

func TestGetOSID(t *testing.T) {
	testCases := []struct {
		name         string
		content      string
		expectedOSID string
	}{
		{
			name: "should return ID from simple file",
			content: `
NAME="Test OS"
VERSION="1.0"
ID=testos
ID_LIKE=anotheros
`,
			expectedOSID: "testos",
		},
		{
			name: "should return ID from quoted syntac with whitespaces",
			content: `
NAME="Test OS"
VERSION="1.0"
 ID = "testos" 
ID_LIKE=anotheros
`,
			expectedOSID: "testos",
		},
		{
			name: "should return ID with whitespaces",
			content: `
NAME="Test OS"
VERSION="1.0"
 ID = testos 
ID_LIKE=anotheros
`,
			expectedOSID: "testos",
		},
		{
			name: "should return ID from quoted syntac with whitespaces",
			content: `
NAME="Test OS"
VERSION="1.0"
ID='testos'
ID_LIKE=anotheros
`,
			expectedOSID: "testos",
		},
		{
			name: "should return the last ID when duplicates exist",
			content: `
ID=firstid
NAME="Test OS"
ID=secondid
VERSION="1.0"
ID=lastid
`,
			expectedOSID: "lastid",
		},
		{
			name:         "should return default string for empty file",
			content:      "",
			expectedOSID: runtime.GOOS,
		},
		{
			name: "should return default string for file with no ID",
			content: `
NAME="Test OS"
VERSION="1.0"
`,
			expectedOSID: runtime.GOOS,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osReleaseFile := createOsReleaseFile(t, tc.content)
			sys := system.NewSystem()
			osID, err := sys.GetOSID(osReleaseFile)
			if err != nil {
				t.Fatalf("unexpected error getting OSID: %v", err)
			}

			if osID != tc.expectedOSID {
				t.Errorf("expected OSID %q, got %q", tc.expectedOSID, osID)
			}
		})
	}
}

func TestGetOSIDMultipleFiles(t *testing.T) {
	t.Run("should return default GOOS when no path is provided", func(t *testing.T) {
		sys := system.NewSystem()
		osID, err := sys.GetOSID()
		if err != nil {
			t.Fatalf("unexpected error getting OSID: %v", err)
		}

		_, err1 := os.Stat("/etc/os-release")
		_, err2 := os.Stat("/usr/lib/os-release")

		if os.IsNotExist(err1) && os.IsNotExist(err2) {
			// Neither file exists, so we should get the default GOOS
			if osID != runtime.GOOS {
				t.Errorf("expected OSID %q, got %q", runtime.GOOS, osID)
			}
		} else {
			// At least one file exists, so we should get a non-empty string
			if osID == "" {
				t.Error("expected a non-empty OSID, but it was empty")
			}
		}
	})

	t.Run("should return ID from the second file if the first does not exist", func(t *testing.T) {
		osReleaseFile := createOsReleaseFile(t, "ID=second")
		sys := system.NewSystem()
		osID, err := sys.GetOSID("/non/existent/file", osReleaseFile)
		if err != nil {
			t.Fatalf("unexpected error getting OSID: %v", err)
		}
		if osID != "second" {
			t.Errorf("expected OSID %q, got %q", "second", osID)
		}
	})

	t.Run("should return default GOOS if the first file has no ID", func(t *testing.T) {
		firstFile := createOsReleaseFile(t, "NAME=No ID here")
		secondFile := createOsReleaseFile(t, "ID=second")
		sys := system.NewSystem()
		osID, err := sys.GetOSID(firstFile, secondFile)
		if err != nil {
			t.Fatalf("unexpected error getting OSID: %v", err)
		}
		if osID != runtime.GOOS {
			t.Errorf("expected OSID %q, got %q", runtime.GOOS, osID)
		}
	})

	t.Run("should return ID from the first file if both exist", func(t *testing.T) {
		firstFile := createOsReleaseFile(t, "ID=first")
		secondFile := createOsReleaseFile(t, "ID=second")
		sys := system.NewSystem()
		osID, err := sys.GetOSID(firstFile, secondFile)
		if err != nil {
			t.Fatalf("unexpected error getting OSID: %v", err)
		}
		if osID != "first" {
			t.Errorf("expected OSID %q, got %q", "first", osID)
		}
	})
}
