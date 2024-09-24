package integration_test

import (
	"encoding/json"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"pkgstats-cli/internal/system"
)

func requiresPacman(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("tests require pacman to be installed")
	}
	if _, err := exec.LookPath("pacman-conf"); err != nil {
		t.Skip("tests require pacman-conf to be installed")
	}
}

func TestShowHelp(t *testing.T) {
	output, err := pkgstats("help")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected help output to contain 'Usage:', got %s", output)
	}
}

func TestShowVersion(t *testing.T) {
	output, err := pkgstats("version")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if !strings.Contains(output, "version") {
		t.Errorf("Expected version output to contain 'version', got %s", output)
	}
}

func TestShowInformationToBeSent(t *testing.T) {
	requiresPacman(t)

	system := system.NewSystem()
	osArchitecture, _ := system.GetArchitecture()
	cpuArchitecture, _ := system.GetCpuArchitecture()

	output, err := pkgstats("submit", "--dump-json")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if result["version"] != "3" {
		t.Errorf("Expected version 3, got %v", result["version"])
	}
	if result["system"].(map[string]interface{})["architecture"] != cpuArchitecture {
		t.Errorf("Expected system architecture '%s', got %v", cpuArchitecture, result["system"].(map[string]interface{})["architecture"])
	}
	if result["os"].(map[string]interface{})["architecture"] != osArchitecture {
		t.Errorf("Expected OS architecture '%s', got %v", osArchitecture, result["os"].(map[string]interface{})["architecture"])
	}
	if !strings.HasPrefix(result["pacman"].(map[string]interface{})["mirror"].(string), "https://") {
		t.Errorf("Expected pacman mirror to start with 'https://', got %v", result["pacman"].(map[string]interface{})["mirror"])
	}
	if !slices.Contains(result["pacman"].(map[string]interface{})["packages"].([]interface{}), "pacman-mirrorlist") {
		t.Errorf("Expected pacman packages to contain 'pacman-mirrorlist'")
	}
}

func TestSetQuietMode(t *testing.T) {
	requiresPacman(t)

	output, err := pkgstats("submit", "--quiet")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if len(output) != 0 {
		t.Errorf("Expected no output in quiet mode, got %s", output)
	}
}

func TestSendInformation(t *testing.T) {
	requiresPacman(t)

	output, err := pkgstats("submit")
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

func TestSearchPackages(t *testing.T) {
	output, err := pkgstats("search", "php")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	lines := strings.Split(output, "\n")
	if !strings.Contains(lines[0], "php") {
		t.Errorf("Expected 'php', got %s", lines[0])
	}
	if !strings.Contains(lines[1], "php-fpm") {
		t.Errorf("Expected 'php-fpm', got %s", lines[1])
	}
}

func TestShowPackages(t *testing.T) {
	output, err := pkgstats("show", "php", "pacman")
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	lines := strings.Split(output, "\n")
	if !strings.Contains(lines[0], "php") {
		t.Errorf("Expected 'php', got %s", lines[0])
	}
	if !strings.Contains(lines[1], "pacman") {
		t.Errorf("Expected 'pacman', got %s", lines[1])
	}
}
