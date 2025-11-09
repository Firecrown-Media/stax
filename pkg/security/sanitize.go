package security

import (
	"regexp"
	"strings"
)

// SanitizeForShell escapes shell metacharacters using whitelist approach
// Returns sanitized string or empty if input contains unsafe characters
func SanitizeForShell(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	// Whitelist approach - allow only safe characters
	// Safe: alphanumeric, forward slash, underscore, period, hyphen
	safePattern := regexp.MustCompile(`^[a-zA-Z0-9/_.-]+$`)

	if !safePattern.MatchString(input) {
		return "", &ValidationError{
			Field:   "shell_input",
			Message: "input contains unsafe characters for shell execution",
		}
	}

	return input, nil
}

// SanitizeCommandArgs sanitizes command arguments
// Returns sanitized args or error if unsafe characters detected
func SanitizeCommandArgs(args []string) ([]string, error) {
	sanitized := make([]string, 0, len(args))

	for _, arg := range args {
		// Check for shell metacharacters
		if containsShellMetachars(arg) {
			return nil, &ValidationError{
				Field:   "command_arg",
				Message: "argument contains shell metacharacters",
			}
		}

		sanitized = append(sanitized, arg)
	}

	return sanitized, nil
}

// RemoveSensitiveData strips credentials from logs/errors
func RemoveSensitiveData(message string) string {
	if message == "" {
		return ""
	}

	// Define credential patterns to redact (order matters - most specific first)
	patterns := []struct {
		pattern     string
		replacement string
	}{
		// Authorization header (most specific first)
		{`Authorization:\s*Bearer\s+[^\s]+`, "Authorization: ***"},
		{`Authorization:\s*[^\s]+`, "Authorization: ***"},

		// Database connection strings
		{`mysql://[^\s@]+:[^\s@]+@`, "mysql://***:***@"},
		{`postgres://[^\s@]+:[^\s@]+@`, "postgres://***:***@"},

		// SSH key content (PEM format)
		{`-----BEGIN[A-Z\s]+PRIVATE KEY-----[^-]+-----END[A-Z\s]+PRIVATE KEY-----`, "-----BEGIN PRIVATE KEY----- [REDACTED] -----END PRIVATE KEY-----"},

		// Password patterns
		{`password["\s:=]+[^\s"&]+`, "password=***"},
		{`pass["\s:=]+[^\s"&]+`, "pass=***"},
		{`passwd["\s:=]+[^\s"&]+`, "passwd=***"},

		// Token patterns
		{`token["\s:=]+[^\s"&]+`, "token=***"},
		{`api_token["\s:=]+[^\s"&]+`, "api_token=***"},
		{`access_token["\s:=]+[^\s"&]+`, "access_token=***"},

		// API key patterns
		{`api_key["\s:=]+[^\s"&]+`, "api_key=***"},
		{`apikey["\s:=]+[^\s"&]+`, "apikey=***"},
		{`key["\s:=]+[^\s"&]+`, "key=***"},

		// Generic bearer token
		{`Bearer\s+[^\s]+`, "Bearer ***"},
	}

	result := message

	for _, p := range patterns {
		re := regexp.MustCompile(`(?i)` + p.pattern)
		result = re.ReplaceAllString(result, p.replacement)
	}

	return result
}

// containsShellMetachars checks if string contains shell metacharacters
func containsShellMetachars(s string) bool {
	// Shell metacharacters that could be used for injection
	metachars := []string{
		";",  // Command separator
		"|",  // Pipe
		"&",  // Background/AND
		"$",  // Variable expansion
		"`",  // Command substitution
		"(",  // Subshell
		")",  // Subshell
		"<",  // Redirect input
		">",  // Redirect output
		"\n", // Newline
		"\r", // Carriage return
		"*",  // Glob
		"?",  // Glob
		"[",  // Glob
		"]",  // Glob
		"{",  // Brace expansion
		"}",  // Brace expansion
		"\\", // Escape
		"'",  // Quote
		"\"", // Quote
	}

	for _, char := range metachars {
		if strings.Contains(s, char) {
			return true
		}
	}

	return false
}

// SanitizeLogMessage removes sensitive data and limits message length
func SanitizeLogMessage(message string, maxLength int) string {
	// Remove sensitive data first
	sanitized := RemoveSensitiveData(message)

	// Limit length
	if maxLength > 0 && len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength] + "... [truncated]"
	}

	return sanitized
}

// SanitizeErrorForUser creates user-friendly error message without sensitive details
func SanitizeErrorForUser(err error) string {
	if err == nil {
		return ""
	}

	message := err.Error()

	// Remove sensitive data
	message = RemoveSensitiveData(message)

	// Remove internal paths
	pathPattern := regexp.MustCompile(`/[^\s]*stax[^\s]*`)
	message = pathPattern.ReplaceAllString(message, "[internal path]")

	// Remove IP addresses (optional - might be useful for debugging)
	// ipPattern := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	// message = ipPattern.ReplaceAllString(message, "[IP address]")

	return message
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// SanitizeFilename removes unsafe characters from filenames
func SanitizeFilename(filename string) (string, error) {
	if filename == "" {
		return "", &ValidationError{
			Field:   "filename",
			Message: "filename cannot be empty",
		}
	}

	// Remove path separators
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	// Allow only safe characters in filenames
	safePattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

	if !safePattern.MatchString(filename) {
		return "", &ValidationError{
			Field:   "filename",
			Message: "filename contains unsafe characters",
		}
	}

	// Check for dangerous filenames
	dangerous := []string{"..", ".", "~"}
	for _, d := range dangerous {
		if filename == d {
			return "", &ValidationError{
				Field:   "filename",
				Message: "filename is reserved: " + d,
			}
		}
	}

	// Limit length
	if len(filename) > 255 {
		return "", &ValidationError{
			Field:   "filename",
			Message: "filename too long (max 255 characters)",
		}
	}

	return filename, nil
}

// SanitizeWPCLIArgs sanitizes WP-CLI arguments
func SanitizeWPCLIArgs(args []string) ([]string, error) {
	sanitized := make([]string, 0, len(args))

	for _, arg := range args {
		// Check for dangerous patterns
		if strings.Contains(arg, "$(") || strings.Contains(arg, "`") {
			return nil, &ValidationError{
				Field:   "wp_cli_arg",
				Message: "argument contains command substitution",
			}
		}

		// Allow WP-CLI specific characters but block shell injection
		if containsShellMetachars(arg) {
			// Exception: hyphens for flags, equals for key=value
			cleanArg := strings.ReplaceAll(arg, "-", "")
			cleanArg = strings.ReplaceAll(cleanArg, "=", "")
			cleanArg = strings.ReplaceAll(cleanArg, "/", "")
			cleanArg = strings.ReplaceAll(cleanArg, ".", "")
			cleanArg = strings.ReplaceAll(cleanArg, "_", "")
			cleanArg = strings.ReplaceAll(cleanArg, ":", "")
			cleanArg = strings.ReplaceAll(cleanArg, ",", "")

			if containsShellMetachars(cleanArg) {
				return nil, &ValidationError{
					Field:   "wp_cli_arg",
					Message: "argument contains shell metacharacters",
				}
			}
		}

		sanitized = append(sanitized, arg)
	}

	return sanitized, nil
}
