package request

import (
	"strings"
)

func ValidatePackageName(name string) bool {
	const alpha_numeric = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const special = "@_+"
	const delimiter = ".-"

	if len(name) < 1 || len(name) > 190 {
		return false
	}

	if !strings.ContainsAny(name[0:1], alpha_numeric) && !strings.ContainsAny(name[0:1], special) {
		return false
	}

	for _, char := range name[1:] {
		if !strings.ContainsRune(alpha_numeric, char) && !strings.ContainsRune(special, char) && !strings.ContainsRune(delimiter, char) {
			return false
		}
	}

	return true
}

func ValidatePackageNames(names []string) bool {
	for _, name := range names {
		if !ValidatePackageName(name) {
			return false
		}
	}

	return true
}
