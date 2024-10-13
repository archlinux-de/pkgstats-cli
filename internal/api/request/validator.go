package request

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	minLength      = 1
	maxLength      = 190
	alphaNumeric   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special        = "@_+"
	delimiter      = ".-"
	validFirstChar = alphaNumeric + special
	validChars     = validFirstChar + delimiter
)

func ValidatePackageName(name string) error {
	length := len(name)
	if length < minLength || length > maxLength {
		return fmt.Errorf("length must be between %d and %d characters, %d were given", minLength, maxLength, length)
	}

	firstChar, _ := utf8.DecodeRuneInString(name)
	if !strings.ContainsRune(validFirstChar, firstChar) {
		return fmt.Errorf("invalid character '%c'", firstChar)
	}

	for _, char := range name[1:] {
		if !strings.ContainsRune(validChars, char) {
			return fmt.Errorf("invalid character '%c'", char)
		}
	}

	return nil
}

func ValidatePackageNames(names []string) error {
	for _, name := range names {
		if err := ValidatePackageName(name); err != nil {
			return fmt.Errorf("'%s' is invalid: %w", name, err)
		}
	}

	return nil
}
