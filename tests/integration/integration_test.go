package integration_test

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/system"
)

func TestShowHelp(t *testing.T) {
	output, err := pkgstats(t, []string{"help"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected help output to contain 'Usage:', got %s", output)
	}
}

func TestShowVersion(t *testing.T) {
	output, err := pkgstats(t, []string{"version"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if !strings.Contains(output, "version") {
		t.Errorf("Expected version output to contain 'version', got %s", output)
	}
}

func TestShowInformationToBeSent(t *testing.T) {
	s := system.NewSystem()
	osArchitecture, err := s.GetArchitecture()
	if err != nil {
		t.Fatal(err)
	}
	cpuArchitecture, err := s.GetCpuArchitecture()
	if err != nil {
		t.Fatal(err)
	}

	output, err := pkgstats(t, []string{"submit", "--dump-json"}, WithPkgBlocklist([]string{"secret-*"}), WithMirrorBlocklist([]string{"secret.com"}))
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	var request submit.Request
	if err := json.Unmarshal([]byte(output), &request); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if request.Version != submit.Version {
		t.Errorf("Expected version %s, got %v", submit.Version, request.Version)
	}
	if request.System.Architecture != cpuArchitecture {
		t.Errorf("Expected system architecture '%s', got %v", cpuArchitecture, request.System.Architecture)
	}
	if request.OS.Architecture != osArchitecture {
		t.Errorf("Expected OS architecture '%s', got %v", osArchitecture, request.OS.Architecture)
	}
	if !strings.HasPrefix(request.Pacman.Mirror, "https://") {
		t.Errorf("Expected pacman mirror to start with 'https://', got %v", request.Pacman.Mirror)
	}
	if !slices.Contains(request.Pacman.Packages, "pacman-mirrorlist") {
		t.Errorf("Expected pacman packages to contain 'pacman-mirrorlist'")
	}
	if slices.Contains(request.Pacman.Packages, "secret-package") {
		t.Errorf("Expected pacman packages to not contain 'secret-package'")
	}
}

func TestMirrorFiltering(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "1" {
		t.Skip("Skipping mirror filtering test in integration mode")
	}

	output, err := pkgstats(t, []string{"submit", "--dump-json"}, WithMirror("http://my.secret.mirror/"), WithMirrorBlocklist([]string{"my.secret.mirror"}))
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	var request submit.Request
	if err := json.Unmarshal([]byte(output), &request); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if request.Pacman.Mirror != "" {
		t.Errorf("Expected pacman mirror to be empty, got %v", request.Pacman.Mirror)
	}
}

func TestMultipleBlocklists(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "1" {
		t.Skip("Skipping multiple blocklists test in integration mode")
	}

	output, err := pkgstats(t, []string{"submit", "--dump-json"},
		WithPkgBlocklist([]string{"secret-*", "*-dev"}),
		WithMirrorBlocklist([]string{"private.mirror.com", "*.internal.net"}),
		WithMirror("http://private.mirror.com/archlinux/$repo/os/$arch"),
	)
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	var request submit.Request
	if err := json.Unmarshal([]byte(output), &request); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if request.Pacman.Mirror != "" {
		t.Errorf("Expected pacman mirror to be empty, got %v", request.Pacman.Mirror)
	}
	if slices.Contains(request.Pacman.Packages, "secret-package") {
		t.Errorf("Expected pacman packages to not contain 'secret-package'")
	}
	if slices.Contains(request.Pacman.Packages, "my-app-dev") {
		t.Errorf("Expected pacman packages to not contain 'my-app-dev'")
	}
}

func TestComplexGlobPatterns(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "1" {
		t.Skip("Skipping complex glob patterns test in integration mode")
	}

	output, err := pkgstats(t, []string{"submit", "--dump-json"},
		WithPkgBlocklist([]string{"app-[0-9]*.pkg"}),
		WithMirrorBlocklist([]string{"*.example.co", "*.example.net"}),
		WithMirror("http://mirror.example.co/archlinux/$repo/os/$arch"),
	)
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	var request submit.Request
	if err := json.Unmarshal([]byte(output), &request); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if request.Pacman.Mirror != "" {
		t.Errorf("Expected pacman mirror to be empty, got %v", request.Pacman.Mirror)
	}
	if slices.Contains(request.Pacman.Packages, "app-123.pkg") {
		t.Errorf("Expected pacman packages to not contain 'app-123.pkg'")
	}
}

func TestSetQuietMode(t *testing.T) {
	output, err := pkgstats(t, []string{"submit", "--quiet"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if len(output) != 0 {
		t.Errorf("Expected no output in quiet mode, got %s", output)
	}
}

func TestSetQuietModeAndDumpCannotBeUsedTogether(t *testing.T) {
	if _, err := pkgstats(t, []string{"submit", "--dump-json", "--quiet"}); err == nil {
		t.Fatal("Command should fail")
	}
}

func TestSendInformation(t *testing.T) {
	output, err := pkgstats(t, []string{"submit"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	lines := strings.Split(output, "\n")
	if !strings.Contains(lines[0], "Collecting data") {
		t.Errorf("Expected 'Collecting data', got %s", lines[0])
	}
	if !strings.Contains(lines[1], "Submitting data") {
		t.Errorf("Expected 'Submitting data', got %s", lines[1])
	}
	if !strings.Contains(lines[2], "Data were successfully sent") {
		t.Errorf("Expected 'Data were successfully sent', got %s", lines[2])
	}
}

func linesContainsPackageStatistic(t *testing.T, lines []string, packages []string) {
	t.Helper()
	for _, e := range packages {
		matcher := regexp.MustCompile(fmt.Sprintf(`^%s\s+\d+\.\d+$`, regexp.QuoteMeta(e)))
		if !slices.ContainsFunc(lines, func(v string) bool { return matcher.MatchString(v) }) {
			t.Errorf("Expected to find '%s' in %v", e, lines)
		}
	}
}

func TestSearchPackages(t *testing.T) {
	output, err := pkgstats(t, []string{"search", "php"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	linesContainsPackageStatistic(t, strings.Split(output, "\n"), []string{"php", "php-fpm"})
}

func TestShowPackages(t *testing.T) {
	output, err := pkgstats(t, []string{"show", "php", "pacman"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	linesContainsPackageStatistic(t, strings.Split(output, "\n"), []string{"php", "pacman"})
}

func TestShowArchitecture(t *testing.T) {
	output, err := pkgstats(t, []string{"architecture"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	s := system.NewSystem()
	osArchitecture, err := s.GetArchitecture()
	if err != nil {
		t.Fatal(err)
	}
	cpuArchitecture, err := s.GetCpuArchitecture()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, osArchitecture) || !strings.Contains(output, cpuArchitecture) {
		t.Fatalf("Expected OS and CPU architecture %s and %s, but got %s", osArchitecture, cpuArchitecture, strings.TrimSpace(output))
	}
}

func TestShowOsArchitecture(t *testing.T) {
	output, err := pkgstats(t, []string{"architecture", "os"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	s := system.NewSystem()
	osArchitecture, err := s.GetArchitecture()
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(output) != osArchitecture {
		t.Fatalf("Expected OS architecture %s, but got %s", osArchitecture, strings.TrimSpace(output))
	}
}

func TestShowSystemArchitecture(t *testing.T) {
	output, err := pkgstats(t, []string{"architecture", "system"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	s := system.NewSystem()
	cpuArchitecture, err := s.GetCpuArchitecture()
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(output) != cpuArchitecture {
		t.Fatalf("Expected CPU architecture %s, but got %s", cpuArchitecture, strings.TrimSpace(output))
	}
}
