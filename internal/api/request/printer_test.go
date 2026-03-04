package request_test

import (
	"bytes"
	"regexp"
	"testing"

	"pkgstats-cli/internal/api/request"
)

func TestPrintPackagePopularities(t *testing.T) {
	var buffer bytes.Buffer
	request.PrintPackagePopularities(&buffer, request.PackagePopularityList{Total: 1, Count: 1, PackagePopularities: []request.PackagePopularity{{Name: "foo", Popularity: 12.34}}})

	output := buffer.String()

	if !regexp.MustCompile(`^foo\s+12.34\s+1 of 1 results\s+$`).MatchString(output) {
		t.Errorf("Unexpected output %s", output)
	}
}

func TestPrintSearchURL(t *testing.T) {
	var buffer bytes.Buffer
	request.PrintSearchURL(&buffer, "/foo", "bar")

	output := buffer.String()

	if !regexp.MustCompile(`\s+/foo/packages#query=bar\s+`).MatchString(output) {
		t.Errorf("Unexpected output %s", output)
	}
}

func TestPrintShowURL(t *testing.T) {
	var buffer bytes.Buffer
	request.PrintShowURL(&buffer, "/foo", []string{"bar", "baz"})

	output := buffer.String()

	if !regexp.MustCompile(`\s+/foo/compare/packages#packages=bar,baz\s+`).MatchString(output) {
		t.Errorf("Unexpected output %s", output)
	}
}
