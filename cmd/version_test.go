package cmd_test

import (
	"strings"
	"testing"
)

func TestShowVersion(t *testing.T) {
	output, err := pkgstats(t, []string{"version"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	if !strings.Contains(output, "version") {
		t.Errorf("Expected version output to contain 'version', got %s", output)
	}
}
