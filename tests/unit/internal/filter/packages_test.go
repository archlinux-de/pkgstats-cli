package filter_test

import (
	"testing"

	"pkgstats-cli/internal/filter"
)

func TestIsFilteredPackage(t *testing.T) {
	tests := []struct {
		pkg      string
		filters  []string
		expected bool
		hasError bool
	}{
		// Filtered inputs
		{"secret-package", []string{"secret-*"}, true, false},
		{"my-secret-package", []string{"secret-*"}, false, false},
		{"super-secret-package", []string{"*secret*"}, true, false},
		{"debug-info", []string{"*-debug", "debug-*"}, true, false},
		{"package", []string{"package"}, true, false},

		// Unfiltered inputs
		{"my-package", []string{"secret-*"}, false, false},
		{"another-package", []string{"*debug*"}, false, false},
		{"package", []string{"pkg"}, false, false},

		// Edge cases
		{"package", []string{}, false, false},
		{"", []string{"*"}, true, false},
		{"package", []string{"["}, false, true},
	}

	for _, test := range tests {
		result, err := filter.IsFilteredPackage(test.filters, test.pkg)
		if (err != nil) != test.hasError {
			t.Errorf("IsFilteredPackage(%q, %q) error = %v, wantErr %v", test.filters, test.pkg, err, test.hasError)
			continue
		}
		if result != test.expected {
			t.Errorf("IsFilteredPackage(%q, %q) = %v, want %v", test.filters, test.pkg, result, test.expected)
		}
	}
}
