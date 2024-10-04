package integration_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/system"
)

func TestShowHelp(t *testing.T) {
	output, err := pkgstats(t, "help")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected help output to contain 'Usage:', got %s", output)
	}
}

func TestShowVersion(t *testing.T) {
	output, err := pkgstats(t, "version")
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

	output, err := pkgstats(t, "submit", "--dump-json")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	var request submit.Request
	if err := json.Unmarshal([]byte(output), &request); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if request.Version != "3" {
		t.Errorf("Expected version 3, got %v", request.Version)
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
}

func TestSetQuietMode(t *testing.T) {
	output, err := pkgstats(t, "submit", "--quiet")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if len(output) != 0 {
		t.Errorf("Expected no output in quiet mode, got %s", output)
	}
}

func TestSendInformation(t *testing.T) {
	output, err := pkgstats(t, "submit")
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
	output, err := pkgstats(t, "search", "php")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	linesContainsPackageStatistic(t, strings.Split(output, "\n"), []string{"php", "php-fpm"})
}

func TestShowPackages(t *testing.T) {
	output, err := pkgstats(t, "show", "php", "pacman")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	linesContainsPackageStatistic(t, strings.Split(output, "\n"), []string{"php", "pacman"})
}

func TestSohwOsArchitecture(t *testing.T) {
	output, err := pkgstats(t, "arch", "os")
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

func TestSohwCpuArchitecture(t *testing.T) {
	output, err := pkgstats(t, "arch", "cpu")
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
