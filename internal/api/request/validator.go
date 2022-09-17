package request

import "regexp"

func ValidatePackageName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9@:\.+_-]{0,190}$`).MatchString(name)
}

func ValidatePackageNames(names []string) bool {
	for _, name := range names {
		if !ValidatePackageName(name) {
			return false
		}
	}

	return true
}
