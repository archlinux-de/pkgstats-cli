package cmd_test

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"
)

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

func TestSearchRequiresArgument(t *testing.T) {
	_, err := pkgstats(t, []string{"search"})
	if err == nil {
		t.Fatal("Expected error when no argument provided")
	}
}

func TestSearchRejectsInvalidPackageName(t *testing.T) {
	_, err := pkgstats(t, []string{"search", "ö"})
	if err == nil {
		t.Fatal("Expected error for invalid package name")
	}
}

func TestSearchRejectsInvalidLimit(t *testing.T) {
	_, err := pkgstats(t, []string{"search", "--limit", "0", "php"})
	if err == nil {
		t.Fatal("Expected error for limit below minimum")
	}
}

func TestSearchWithCustomLimit(t *testing.T) {
	output, err := pkgstats(t, []string{"search", "--limit", "5", "php"})
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	linesContainsPackageStatistic(t, strings.Split(output, "\n"), []string{"php"})
}
