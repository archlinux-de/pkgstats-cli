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

func TestShowRequiresArguments(t *testing.T) {
	_, err := pkgstats(t, []string{"show"})
	if err == nil {
		t.Fatal("Expected error when no arguments provided")
	}
}

func TestShowRejectsInvalidPackageName(t *testing.T) {
	_, err := pkgstats(t, []string{"show", "ö"})
	if err == nil {
		t.Fatal("Expected error for invalid package name")
	}
}

func TestShowSinglePackage(t *testing.T) {
	output, err := pkgstats(t, []string{"show", "php"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	linesContainsPackageStatistic(t, strings.Split(output, "\n"), []string{"php"})
}
