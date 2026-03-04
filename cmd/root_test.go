package cmd_test

import (
	"strings"
	"testing"
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
