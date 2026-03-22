package validator

import (
	"fmt"
	"strings"
)

// Validate checks if a namespace name follows Kubernetes naming rules
//
// Rules:
// - Must be lowercase
// - Only alphanumeric characters and hyphens
// - Length: 1-63 characters
// - Cannot start or end with hyphen
//
// Returns:
//   - bool: true if valid, false if invalid
//   - string: empty if valid, error message if invalid
func Validate(name string) (bool, string) {
	// Rule 1: Check length
	if len(name) == 0 {
		return false, "namespace name cannot be empty"
	}
	if len(name) > 63 {
		return false, fmt.Sprintf("namespace name too long (%d characters, max 63)", len(name))
	}

	// Rule 2: Cannot start with hyphen
	if name[0] == '-' {
		return false, "namespace name cannot start with a hyphen"
	}

	// Rule 3: Cannot end with hyphen
	if name[len(name)-1] == '-' {
		return false, "namespace name cannot end with a hyphen"
	}

	// Rule 4: Must be lowercase
	if name != strings.ToLower(name) {
		return false, "namespace name must be lowercase (no uppercase letters allowed)"
	}

	// Rule 5: Only lowercase letters, numbers, and hyphens
	for i, char := range name {
		if !isValidChar(char) {
			return false, fmt.Sprintf("invalid character '%c' at position %d (only a-z, 0-9, and - allowed)", char, i)
		}
	}

	return true, ""
}

// isValidChar checks if a character is allowed in a namespace name
func isValidChar(char rune) bool {
	isLowercaseLetter := char >= 'a' && char <= 'z'
	isDigit := char >= '0' && char <= '9'
	isHyphen := char == '-'

	return isLowercaseLetter || isDigit || isHyphen
}

// ValidateWithSuggestion validates a name and provides helpful suggestions
func ValidateWithSuggestion(name string) (bool, string, string) {
	valid, reason := Validate(name)
	if valid {
		return true, "", ""
	}

	// Generate suggestion based on error type
	suggestion := generateSuggestion(name, reason)
	return false, reason, suggestion
}

// generateSuggestion creates a helpful suggestion based on the validation error
func generateSuggestion(name, reason string) string {
	// If uppercase, suggest lowercase version
	if strings.Contains(reason, "lowercase") {
		return fmt.Sprintf("Try: %s", strings.ToLower(name))
	}

	// If contains underscore, suggest replacing with hyphen
	if strings.Contains(reason, "_") {
		suggested := strings.ReplaceAll(name, "_", "-")
		return fmt.Sprintf("Try: %s", suggested)
	}

	// If starts/ends with hyphen, suggest trimming
	if strings.Contains(reason, "start") || strings.Contains(reason, "end") {
		suggested := strings.Trim(name, "-")
		return fmt.Sprintf("Try: %s", suggested)
	}

	// Generic suggestion
	return "Use only lowercase letters, numbers, and hyphens"
}

// ReservedNamespaces lists Kubernetes system namespaces that shouldn't be created
var ReservedNamespaces = []string{
	"default",
	"kube-system",
	"kube-public",
	"kube-node-lease",
}

// IsReserved checks if a namespace name is reserved by Kubernetes
func IsReserved(name string) bool {
	for _, reserved := range ReservedNamespaces {
		if name == reserved {
			return true
		}
	}
	return false
}
