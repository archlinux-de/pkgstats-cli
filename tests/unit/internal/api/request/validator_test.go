package request_test

import (
	"testing"

	"pkgstats-cli/internal/api/request"
)

var (
	validPackageNames   = []string{"pacman", "@pacman", "_pacman", "+pacman", "pacman-contrib", "pacman_foo", "pacman@7", "pacman+bar", "pacman.bar"}
	invalidPackageNames = []string{"-pacman", ".pacman", "ö", "", "pacman-ö", "🐱", "🐱nip"}
)

func TestValidPackageName(t *testing.T) {
	for _, name := range validPackageNames {
		if err := request.ValidatePackageName(name); err != nil {
			t.Error("name should be valid:", name, err)
		}
	}
}

func TestInvalidPackageName(t *testing.T) {
	for _, name := range invalidPackageNames {
		if err := request.ValidatePackageName(name); err == nil {
			t.Error("name should be invalid")
		}
	}
}

func TestValidPackageNames(t *testing.T) {
	for _, name := range validPackageNames {
		if err := request.ValidatePackageNames(append([]string{"foo"}, name)); err != nil {
			t.Error("names should be valid", name, err)
		}
	}
}

func TestInvalidPackageNames(t *testing.T) {
	for _, name := range invalidPackageNames {
		if err := request.ValidatePackageNames(append([]string{"foo"}, name)); err == nil {
			t.Error("names should be invalid")
		}
	}
}
