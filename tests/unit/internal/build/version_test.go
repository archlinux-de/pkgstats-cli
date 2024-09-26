package build_test

import (
	"testing"

	"pkgstats-cli/internal/build"
)

func TestVersion(t *testing.T) {
	if build.Version != "dev" {
		t.Error("Version is not defined")
	}
}
