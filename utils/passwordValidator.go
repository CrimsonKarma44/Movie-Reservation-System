package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidatePassword checks password strength according to OWASP recommendations
// Requirements:
// - Minimum 12 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one digit
// - At least one special character from: !@#$%^&*
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters long (provided: %d)", len(password))
	}

	// Check for uppercase letters
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for lowercase letters
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check for digits
	if !strings.ContainsAny(password, "0123456789") {
		return fmt.Errorf("password must contain at least one digit")
	}

	// Check for special characters
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	if !strings.ContainsAny(password, specialChars) {
		return fmt.Errorf("password must contain at least one special character from: !@#$%%^&*()")
	}

	// Check for common weak patterns
	if isCommonPassword(password) {
		return fmt.Errorf("password is too common, please choose a more unique password")
	}

	return nil
}

// isCommonPassword checks if password matches common patterns
func isCommonPassword(password string) bool {
	commonPatterns := []string{
		"123456",
		"password",
		"qwerty",
		"abc123",
		"111111",
		"000000",
		"admin",
		"letmein",
		"welcome",
		"monkey",
		"dragon",
		"master",
		"sunshine",
	}

	lower := strings.ToLower(password)
	for _, pattern := range commonPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}

// GetPasswordStrength returns the strength level of a password
func GetPasswordStrength(password string) string {
	if len(password) < 8 {
		return "Very Weak"
	}

	strength := 0

	// Length points
	if len(password) >= 12 {
		strength++
	}
	if len(password) >= 16 {
		strength++
	}

	// Character variety
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasDigit := strings.ContainsAny(password, "0123456789")
	hasSpecial := false

	for _, r := range password {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			hasSpecial = true
			break
		}
	}

	if hasUpper {
		strength++
	}
	if hasLower {
		strength++
	}
	if hasDigit {
		strength++
	}
	if hasSpecial {
		strength++
	}

	switch strength {
	case 0, 1:
		return "Very Weak"
	case 2:
		return "Weak"
	case 3:
		return "Fair"
	case 4:
		return "Good"
	case 5, 6, 7:
		return "Strong"
	default:
		return "Very Strong"
	}
}
