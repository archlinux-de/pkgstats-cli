package cmd_test

import (
	"strings"
	"testing"

	"pkgstats-cli/internal/system"
)

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
