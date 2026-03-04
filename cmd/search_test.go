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
