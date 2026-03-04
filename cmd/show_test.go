package cmd_test

import (
	"strings"
	"testing"
)

func TestShowPackages(t *testing.T) {
	output, err := pkgstats(t, []string{"show", "php", "pacman"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	linesContainsPackageStatistic(t, strings.Split(output, "\n"), []string{"php", "pacman"})
}
