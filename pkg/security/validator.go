package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidateProjectName ensures safe project names
// Only allows alphanumeric characters, hyphens, and underscores
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) > 64 {
		return fmt.Errorf("project name too long (max 64 characters)")
	}

	// Allow only alphanumeric, hyphens, and underscores
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
		return fmt.Errorf("project name contains invalid characters (only alphanumeric, hyphens, and underscores allowed)")
	}

	// Prevent names that start with hyphen (could be confused with flags)
	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("project name cannot start with hyphen")
	}

	return nil
}

// ValidateHostname ensures safe hostnames following RFC 1123
func ValidateHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	if len(hostname) > 253 {
		return fmt.Errorf("hostname too long (max 253 characters)")
	}

	// RFC 1123 hostname validation
	hostnameRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)

	if !hostnameRegex.MatchString(hostname) {
		return fmt.Errorf("invalid hostname format (must comply with RFC 1123)")
	}

	// Reject reserved names
	reserved := []string{"localhost", "broadcasthost"}
	for _, r := range reserved {
		if strings.EqualFold(hostname, r) {
			return fmt.Errorf("hostname cannot be reserved name: %s", r)
		}
	}

	return nil
}

// ValidateURL ensures safe URLs
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Must start with http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	// Check for unsafe characters
	unsafeChars := []string{" ", "\n", "\r", "\t", "<", ">", "\"", "{", "}", "|", "\\", "^", "`"}
	for _, char := range unsafeChars {
		if strings.Contains(url, char) {
			return fmt.Errorf("URL contains unsafe character: %q", char)
		}
	}

	// Basic length check
	if len(url) > 2048 {
		return fmt.Errorf("URL too long (max 2048 characters)")
	}

	return nil
}

// SanitizePath prevents path traversal attacks
func SanitizePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Check for path traversal patterns
	if IsPathTraversal(path) {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}

	// Clean the path
	cleaned := filepath.Clean(path)

	// Ensure cleaned path doesn't escape
	if IsPathTraversal(cleaned) {
		return "", fmt.Errorf("path contains traversal after cleaning: %s", cleaned)
	}

	return cleaned, nil
}

// ValidateFilePath ensures file path is within allowed directory
func ValidateFilePath(path, allowedDir string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if allowedDir == "" {
		return fmt.Errorf("allowed directory cannot be empty")
	}

	// Check for path traversal
	if IsPathTraversal(path) {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	// Convert both to absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	absAllowedDir, err := filepath.Abs(allowedDir)
	if err != nil {
		return fmt.Errorf("failed to resolve allowed directory: %w", err)
	}

	// Ensure path is within allowed directory
	relPath, err := filepath.Rel(absAllowedDir, absPath)
	if err != nil {
		return fmt.Errorf("path is outside allowed directory: %w", err)
	}

	// Check if relative path tries to escape
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("path is outside allowed directory: %s", path)
	}

	return nil
}

// ValidateCommand ensures command is in allowlist
func ValidateCommand(cmd string, allowlist []string) error {
	if cmd == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Check if command is in allowlist
	for _, allowed := range allowlist {
		if cmd == allowed {
			return nil
		}
	}

	return fmt.Errorf("command not in allowlist: %s", cmd)
}

// IsPathTraversal checks for path traversal patterns
func IsPathTraversal(path string) bool {
	// Check for various path traversal patterns
	patterns := []string{
		"../",
		"..\\",
		"/..",
		"\\..",
	}

	for _, pattern := range patterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	return false
}

// ValidateTablePrefix validates database table prefix format
func ValidateTablePrefix(prefix string) error {
	if prefix == "" {
		return fmt.Errorf("table prefix cannot be empty")
	}

	// Allow only alphanumeric and underscores
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(prefix) {
		return fmt.Errorf("invalid table prefix format (only alphanumeric and underscores allowed)")
	}

	if len(prefix) > 64 {
		return fmt.Errorf("table prefix too long (max 64 characters)")
	}

	return nil
}

// ValidateTableName validates database table name
func ValidateTableName(tableName string) error {
	if tableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	// Allow only alphanumeric, underscores, and hyphens
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(tableName) {
		return fmt.Errorf("invalid table name format")
	}

	if len(tableName) > 64 {
		return fmt.Errorf("table name too long (max 64 characters)")
	}

	return nil
}

// ValidateRsyncPattern validates rsync include/exclude patterns
func ValidateRsyncPattern(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}

	// Reject dangerous shell metacharacters
	dangerous := []string{";", "|", "&", "$", "`", "(", ")", "<", ">", "\n", "\r"}
	for _, char := range dangerous {
		if strings.Contains(pattern, char) {
			return fmt.Errorf("pattern contains dangerous character: %q", char)
		}
	}

	// Check pattern length
	if len(pattern) > 256 {
		return fmt.Errorf("pattern too long (max 256 characters)")
	}

	return nil
}

// ValidateEnvironment validates environment name (dev, staging, production)
func ValidateEnvironment(env string) error {
	validEnvs := []string{"dev", "development", "staging", "stage", "production", "prod"}

	env = strings.ToLower(strings.TrimSpace(env))

	for _, valid := range validEnvs {
		if env == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid environment: %s (allowed: dev, staging, production)", env)
}
