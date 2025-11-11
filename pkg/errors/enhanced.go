package errors

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Error codes
const (
	ErrCodeConfigNotFound        = "STX-001"
	ErrCodeCredentialsNotFound   = "STX-002"
	ErrCodeSSHKeyNotFound        = "STX-003"
	ErrCodeDDEVNotInstalled      = "STX-004"
	ErrCodeDDEVNotConfigured     = "STX-005"
	ErrCodeCommandNotImplemented = "STX-006"
	ErrCodeInvalidConfig         = "STX-007"
	ErrCodeWPEngineAPI           = "STX-008"
)

// Solution represents a proposed solution to an error
type Solution struct {
	// Description of the solution
	Description string

	// Command to execute (if applicable)
	Command string

	// Steps to manually resolve the issue
	Steps []string
}

// EnhancedError provides detailed error information with actionable solutions
type EnhancedError struct {
	// Error code (e.g., STX-001)
	Code string

	// Short error message
	Message string

	// Detailed explanation
	Details string

	// Proposed solutions
	Solutions []Solution

	// URL to documentation
	DocsURL string

	// Underlying error (if any)
	Err error
}

// Error implements the error interface
func (e *EnhancedError) Error() string {
	return formatEnhancedError(e)
}

// Unwrap returns the underlying error
func (e *EnhancedError) Unwrap() error {
	return e.Err
}

// formatEnhancedError formats an enhanced error for display
func formatEnhancedError(e *EnhancedError) string {
	var sb strings.Builder

	// Error header with code
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	white := color.New(color.FgWhite)

	sb.WriteString("\n")
	sb.WriteString(red.Sprintf("Error [%s]: %s\n", e.Code, e.Message))

	// Details
	if e.Details != "" {
		sb.WriteString("\n")
		sb.WriteString(white.Sprint(e.Details))
		sb.WriteString("\n")
	}

	// Underlying error
	if e.Err != nil {
		sb.WriteString("\n")
		sb.WriteString(yellow.Sprintf("Cause: %v\n", e.Err))
	}

	// Solutions
	if len(e.Solutions) > 0 {
		sb.WriteString("\n")
		sb.WriteString(cyan.Sprint("Suggested Solutions:\n"))
		for i, sol := range e.Solutions {
			sb.WriteString(fmt.Sprintf("\n%d. %s\n", i+1, sol.Description))

			if sol.Command != "" {
				sb.WriteString(fmt.Sprintf("   Run: %s\n", cyan.Sprint(sol.Command)))
			}

			if len(sol.Steps) > 0 {
				sb.WriteString("   Steps:\n")
				for _, step := range sol.Steps {
					sb.WriteString(fmt.Sprintf("   - %s\n", step))
				}
			}
		}
	}

	// Documentation link
	if e.DocsURL != "" {
		sb.WriteString("\n")
		sb.WriteString(cyan.Sprintf("Documentation: %s\n", e.DocsURL))
	}

	sb.WriteString("\n")

	return sb.String()
}

// NewEnhancedError creates a new enhanced error
func NewEnhancedError(code, message, details string, solutions []Solution, docsURL string, err error) *EnhancedError {
	return &EnhancedError{
		Code:      code,
		Message:   message,
		Details:   details,
		Solutions: solutions,
		DocsURL:   docsURL,
		Err:       err,
	}
}

// NewWithSolution creates a simple error with a message, details, and one solution
func NewWithSolution(message, details string, solution Solution) *EnhancedError {
	return &EnhancedError{
		Code:      "",
		Message:   message,
		Details:   details,
		Solutions: []Solution{solution},
		DocsURL:   "",
		Err:       nil,
	}
}
