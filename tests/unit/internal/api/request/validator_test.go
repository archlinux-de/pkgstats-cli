package request_test

import (
	"testing"

	"pkgstats-cli/internal/api/request"
)

var (
	validPackageNames   = []string{"pacman", "@pacman", "_pacman", "+pacman", "pacman-contrib", "pacman_foo", "pacman@7", "pacman+bar", "pacman.bar"}
	invalidPackageNames = []string{"-pacman", ".pacman", "ö", "", "pacman-ö"}
)

func TestValidPackageName(t *testing.T) {
	for _, name := range validPackageNames {
		result := request.ValidatePackageName(name)

		if !result {
			t.Error("name should be valid:", name)
		}
	}
}

func TestInvalidPackageName(t *testing.T) {
	for _, name := range invalidPackageNames {
		result := request.ValidatePackageName(name)

		if result {
			t.Error("name should be invalid")
		}
	}
}

func TestValidPackageNames(t *testing.T) {
	for _, name := range validPackageNames {
		result := request.ValidatePackageNames(append([]string{"foo"}, name))

		if !result {
			t.Error("names should be valid")
		}
	}
}

func TestInvalidPackageNames(t *testing.T) {
	for _, name := range invalidPackageNames {
		result := request.ValidatePackageNames(append([]string{"foo"}, name))

		if result {
			t.Error("names should be invalid")
		}
	}
}
