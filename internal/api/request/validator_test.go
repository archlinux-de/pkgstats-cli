package request

import (
	"testing"
)

func TestValidPackageName(t *testing.T) {
	result := ValidatePackageName("foo")

	if !result {
		t.Error("name should be valid")
	}
}

func TestInvalidPackageName(t *testing.T) {
	result := ValidatePackageName("üß")

	if result {
		t.Error("name should be invalid")
	}
}

func TestValidPackageNames(t *testing.T) {
	result := ValidatePackageNames([]string{"foo", "bar"})

	if !result {
		t.Error("names should be valid")
	}
}

func TestInvalidPackageNames(t *testing.T) {
	result := ValidatePackageNames([]string{"foo", "bar", "_baz"})

	if result {
		t.Error("names should be invalid")
	}
}
