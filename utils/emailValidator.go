package utils

import (
	"fmt"
	"regexp"
)

// ValidateEmail checks if an email address has a valid format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if len(email) > 254 {
		return fmt.Errorf("email is too long (max 254 characters)")
	}

	// RFC 5322 simplified regex for email validation
	// This is not a perfect validation but covers most common cases
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email format is invalid")
	}

	return nil
}

// NormalizeEmail returns the email in lowercase
func NormalizeEmail(email string) string {
	return toLower(email)
}

// Helper function for lowercase conversion
func toLower(s string) string {
	result := make([]byte, len(s))
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			result[i] = byte(c - 'A' + 'a')
		} else {
			result[i] = byte(c)
		}
	}
	return string(result)
}
