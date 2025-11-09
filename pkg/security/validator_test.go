package security

import (
	"strings"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"valid simple", "myproject", false},
		{"valid with hyphen", "my-project", false},
		{"valid with underscore", "my_project", false},
		{"valid mixed", "my-project_123", false},

		// Invalid cases
		{"empty", "", true},
		{"too long", strings.Repeat("a", 65), true},
		{"starts with hyphen", "-project", true},
		{"contains space", "my project", true},
		{"contains special char", "my@project", true},
		{"contains slash", "my/project", true},
		{"contains backslash", "my\\project", true},
		{"contains semicolon", "my;project", true},
		{"contains pipe", "my|project", true},
		{"contains ampersand", "my&project", true},
		{"contains dollar", "my$project", true},
		{"contains backtick", "my`project", true},
		{"contains parenthesis", "my(project", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateProjectName(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateProjectName(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"valid simple", "example.com", false},
		{"valid subdomain", "www.example.com", false},
		{"valid hyphen", "my-site.example.com", false},
		{"valid numbers", "site123.example.com", false},

		// Invalid cases
		{"empty", "", true},
		{"too long", strings.Repeat("a", 254), true},
		{"localhost reserved", "localhost", true},
		{"broadcasthost reserved", "broadcasthost", true},
		{"contains space", "my site.com", true},
		{"contains underscore", "my_site.com", true},
		{"starts with hyphen", "-example.com", true},
		{"ends with hyphen", "example-.com", true},
		{"contains special char", "my@site.com", true},
		{"double dot", "my..site.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHostname(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateHostname(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateHostname(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"valid http", "http://example.com", false},
		{"valid https", "https://example.com", false},
		{"valid with path", "https://example.com/path", false},
		{"valid with query", "https://example.com?query=value", false},

		// Invalid cases
		{"empty", "", true},
		{"no protocol", "example.com", true},
		{"ftp protocol", "ftp://example.com", true},
		{"contains space", "https://example.com/my file", true},
		{"contains newline", "https://example.com\n", true},
		{"contains tab", "https://example.com\t", true},
		{"contains angle bracket", "https://example.com/<script>", true},
		{"too long", "https://" + strings.Repeat("a", 2050), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateURL(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateURL(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestIsPathTraversal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Path traversal attempts
		{"unix traversal", "../etc/passwd", true},
		{"windows traversal", "..\\windows\\system32", true},
		{"hidden traversal", "valid/../../etc", true},
		{"traversal at start", "/../etc/passwd", true},
		{"traversal at end", "/var/log/..", true},

		// Valid paths
		{"simple path", "valid/path/file.txt", false},
		{"absolute path", "/var/www/html", false},
		{"dotfile", ".hidden", false},
		{"dotdir", "./current", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPathTraversal(tt.input)
			if result != tt.expected {
				t.Errorf("IsPathTraversal(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"simple path", "path/to/file", false},
		{"dotfile", ".hidden", false},

		// Invalid cases
		{"empty", "", true},
		{"traversal", "../etc/passwd", true},
		{"windows traversal", "..\\windows", true},
		{"hidden traversal", "valid/../../etc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SanitizePath(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("SanitizePath(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("SanitizePath(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestValidateTablePrefix(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"standard wp", "wp_", false},
		{"custom", "mysite_", false},
		{"with numbers", "wp123_", false},

		// Invalid cases
		{"empty", "", true},
		{"too long", strings.Repeat("a", 65), true},
		{"contains hyphen", "wp-prefix_", true},
		{"contains space", "wp _", true},
		{"contains semicolon", "wp;_", true},
		{"contains quote", "wp'_", true},
		{"sql injection attempt", "wp_'; DROP TABLE users--", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTablePrefix(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateTablePrefix(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateTablePrefix(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestValidateTableName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"standard", "wp_posts", false},
		{"with hyphen", "wp_my-table", false},
		{"with numbers", "wp_table123", false},

		// Invalid cases
		{"empty", "", true},
		{"too long", strings.Repeat("a", 65), true},
		{"contains space", "wp posts", true},
		{"contains semicolon", "wp_posts;", true},
		{"contains quote", "wp_posts'", true},
		{"sql injection", "wp_posts; DROP TABLE users--", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTableName(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateTableName(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateTableName(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestValidateRsyncPattern(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"simple glob", "*.log", false},
		{"directory", "cache/", false},
		{"dotfile", ".DS_Store", false},
		{"complex glob", "**/*.tmp", false},

		// Invalid cases
		{"empty", "", true},
		{"too long", strings.Repeat("a", 257), true},
		{"contains semicolon", "*.log;rm -rf /", true},
		{"contains pipe", "*.log|malicious", true},
		{"contains ampersand", "*.log&malicious", true},
		{"contains dollar", "*.log$malicious", true},
		{"contains backtick", "*.log`malicious`", true},
		{"command substitution", "$(malicious)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRsyncPattern(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateRsyncPattern(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateRsyncPattern(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestValidateCommand(t *testing.T) {
	allowlist := []string{"ls", "cat", "echo"}

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"allowed ls", "ls", false},
		{"allowed cat", "cat", false},
		{"allowed echo", "echo", false},

		// Invalid cases
		{"empty", "", true},
		{"not in allowlist", "rm", true},
		{"malicious", "rm -rf /", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommand(tt.input, allowlist)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateCommand(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateCommand(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestValidateEnvironment(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		// Valid cases
		{"dev", "dev", false},
		{"development", "development", false},
		{"staging", "staging", false},
		{"stage", "stage", false},
		{"production", "production", false},
		{"prod", "prod", false},
		{"uppercase dev", "DEV", false},
		{"mixed case", "Production", false},

		// Invalid cases
		{"empty", "", true},
		{"invalid", "invalid", true},
		{"test", "test", true},
		{"local", "local", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnvironment(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateEnvironment(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateEnvironment(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

// Fuzzing-style tests for path traversal
func TestPathTraversalFuzzing(t *testing.T) {
	maliciousInputs := []string{
		"../",
		"..\\",
		"/../",
		"\\..\\",
		"....//",
		"..../\\",
		"./../",
		".\\..\\",
		"../../../../../../../etc/passwd",
		"..\\..\\..\\..\\..\\windows\\system32",
	}

	for _, input := range maliciousInputs {
		t.Run(input, func(t *testing.T) {
			// Check if the common patterns are detected
			hasTraversal := strings.Contains(input, "../") || strings.Contains(input, "..\\")
			detected := IsPathTraversal(input)

			if hasTraversal && !detected {
				t.Errorf("IsPathTraversal(%q) should detect traversal", input)
			}

			_, err := SanitizePath(input)
			if err == nil && hasTraversal {
				t.Errorf("SanitizePath(%q) should reject traversal", input)
			}
		})
	}
}

// Command injection tests
func TestCommandInjectionPatterns(t *testing.T) {
	maliciousPatterns := []string{
		"; rm -rf /",
		"| cat /etc/passwd",
		"& malicious",
		"$(whoami)",
		"`id`",
		"$(/bin/sh)",
		"|| echo pwned",
		"&& echo pwned",
	}

	for _, pattern := range maliciousPatterns {
		t.Run(pattern, func(t *testing.T) {
			err := ValidateRsyncPattern(pattern)
			if err == nil {
				t.Errorf("ValidateRsyncPattern(%q) should reject malicious pattern", pattern)
			}
		})
	}
}

// SQL injection tests
func TestSQLInjectionPatterns(t *testing.T) {
	maliciousPatterns := []string{
		"wp_'; DROP TABLE users--",
		"wp_' OR '1'='1",
		"wp_\"; DELETE FROM posts--",
		"wp_'; UPDATE users SET admin=1--",
	}

	for _, pattern := range maliciousPatterns {
		t.Run(pattern, func(t *testing.T) {
			err := ValidateTablePrefix(pattern)
			if err == nil {
				t.Errorf("ValidateTablePrefix(%q) should reject SQL injection", pattern)
			}

			err = ValidateTableName(pattern)
			if err == nil {
				t.Errorf("ValidateTableName(%q) should reject SQL injection", pattern)
			}
		})
	}
}
