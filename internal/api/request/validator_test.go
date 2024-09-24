package request

import (
	"testing"
)

var (
	validPackageNames   = []string{"pacman", "@pacman", "_pacman", "+pacman", "pacman-contrib", "pacman_foo", "pacman@7", "pacman+bar", "pacman.bar"}
	invalidPackageNames = []string{"-pacman", ".pacman", "ö", "", "pacman-ö"}
)

func TestValidPackageName(t *testing.T) {
	for _, name := range validPackageNames {
		result := ValidatePackageName(name)

		if !result {
			t.Error("name should be valid:", name)
		}
	}
}

func TestInvalidPackageName(t *testing.T) {
	for _, name := range invalidPackageNames {
		result := ValidatePackageName(name)

		if result {
			t.Error("name should be invalid")
		}
	}
}

func TestValidPackageNames(t *testing.T) {
	for _, name := range validPackageNames {
		result := ValidatePackageNames(append([]string{"foo"}, name))

		if !result {
			t.Error("names should be valid")
		}
	}
}

func TestInvalidPackageNames(t *testing.T) {
	for _, name := range invalidPackageNames {
		result := ValidatePackageNames(append([]string{"foo"}, name))

		if result {
			t.Error("names should be invalid")
		}
	}
}
