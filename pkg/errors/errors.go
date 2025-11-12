package errors

import (
	"fmt"
)

// Error codes for Stax
const (
	ErrCodeConfigNotFound        = "STX-001"
	ErrCodeCredentialsNotFound   = "STX-002"
	ErrCodeSSHKeyNotFound        = "STX-003"
	ErrCodeDDEVNotInstalled      = "STX-101"
	ErrCodeDDEVNotConfigured     = "STX-102"
	ErrCodeCommandNotImplemented = "STX-900"
	ErrCodeInvalidConfig         = "STX-004"
	ErrCodeWPEngineAPI           = "STX-201"
)

// Common error types for Stax

// DDEVError represents an error with DDEV operations
type DDEVError struct {
	Message string
	Err     error
}

func (e *DDEVError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("DDEV error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("DDEV error: %s", e.Message)
}

// NewDDEVError creates a new DDEV error
func NewDDEVError(message string, err error) *DDEVError {
	return &DDEVError{Message: message, Err: err}
}

// WPEngineError represents an error with WPEngine operations
type WPEngineError struct {
	Message string
	Err     error
}

func (e *WPEngineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("WPEngine error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("WPEngine error: %s", e.Message)
}

// NewWPEngineError creates a new WPEngine error
func NewWPEngineError(message string, err error) *WPEngineError {
	return &WPEngineError{Message: message, Err: err}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Configuration error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("Configuration error: %s", e.Message)
}

// NewConfigError creates a new configuration error
func NewConfigError(message string, err error) *ConfigError {
	return &ConfigError{Message: message, Err: err}
}

// CredentialsError represents a credentials error
type CredentialsError struct {
	Message string
	Err     error
}

func (e *CredentialsError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Credentials error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("Credentials error: %s", e.Message)
}

// NewCredentialsError creates a new credentials error
func NewCredentialsError(message string, err error) *CredentialsError {
	return &CredentialsError{Message: message, Err: err}
}
