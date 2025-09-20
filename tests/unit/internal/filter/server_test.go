package filter_test

import (
	"testing"

	"pkgstats-cli/internal/filter"
)

func TestFilterMirrorUrl(t *testing.T) {
	tests := []struct {
		input    string
		filter   []string
		expected bool
		hasError bool
	}{
		// Filtered inputs
		{"http://example.com/path/to/mirror", []string{"example.com"}, true, false},
		{"ftp://example.com/path/to/mirror", []string{"example.com"}, true, false},
		{"ftp://example.com:1234/path/to/mirror", []string{"example.com"}, true, false},
		{"ftp://user:pw@example.com:1234/path/to/mirror", []string{"example.com"}, true, false},
		{"http://example.com/path/to/mirror", []string{"http://example.com/"}, true, false},
		{"https://example.com/path/to/mirror", []string{"http://example.com/"}, true, false},
		{"http://example.com/path/to/mirror", []string{"http://example.com/"}, true, false},
		{"http://example.com/path/to/mirror", []string{"example.com"}, true, false},
		{"http://www.example.com/path/to/mirror", []string{"*.example.com"}, true, false},
		{"http://www.example.de/path/to/mirror", []string{"*.example.*"}, true, false},
		{"example.com", []string{"*"}, true, false},
		{"example.com", []string{"*"}, true, false},
		{"http://example.com/path/to/mirror", []string{"example.*"}, true, false},
		{"http://example.com/path/to/mirror", []string{"example*"}, true, false},

		// Unfiltered inputs
		{"http://example.com/path/to/mirror", []string{"*.example.*"}, false, false},
		{"http://example.com/path/to/mirror", []string{"http://example.com/path/*"}, false, false},
		{"http://example.com/path/to/mirror", []string{"http://example.com/*"}, false, false},
		{"http://example.com/path/to/mirror", []string{"example"}, false, false},
		{"http://example.com/path/to/mirror", []string{"https://example.com/*"}, false, false},
		{"http://example.com/path/to/mirror", []string{"https://example.com/path/*"}, false, false},
		{"http://another.com/path/to/mirror", []string{"example.com"}, false, false},

		// Edge cases
		{"http://example.com/path/to/mirror", []string{}, false, false},
		{"", []string{"*"}, false, false},
		{"http://example.com/path/to/mirror", []string{"["}, false, true},

		// Invalid inputs
		{"://example.com/path/to/mirror", []string{"*"}, false, true},
		{"http://example.com:invalidport/path/to/mirror", []string{"*"}, false, true},
	}

	for _, test := range tests {
		result, err := filter.IsFilteredMirrorUrl(test.filter, test.input)
		if (err != nil) != test.hasError {
			t.Errorf("IsFilteredMirrorUrl(%q, %q) error = %v, wantErr %v", test.filter, test.input, err, test.hasError)
			continue
		}
		if result != test.expected {
			t.Errorf("IsFilteredMirrorUrl(%q, %q) = %v, want %v", test.filter, test.input, result, test.expected)
		}
	}
}
