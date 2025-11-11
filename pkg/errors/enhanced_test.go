package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestNewEnhancedError(t *testing.T) {
	err := NewEnhancedError(
		"STX-001",
		"Test error",
		"This is a test error",
		[]Solution{
			{
				Description: "Fix it",
				Command:     "stax fix",
			},
		},
		"https://example.com/docs",
		errors.New("underlying error"),
	)

	if err.Code != "STX-001" {
		t.Errorf("expected code STX-001, got %s", err.Code)
	}

	if err.Message != "Test error" {
		t.Errorf("expected message 'Test error', got %s", err.Message)
	}

	if len(err.Tried) != 0 {
		t.Errorf("expected empty Tried array, got %d items", len(err.Tried))
	}
}

func TestNewEnhancedErrorWithTried(t *testing.T) {
	tried := []string{"Location 1", "Location 2", "Location 3"}
	err := NewEnhancedErrorWithTried(
		"STX-002",
		"Test error with tried",
		"This is a test error with tried locations",
		tried,
		[]Solution{},
		"",
		nil,
	)

	if err.Code != "STX-002" {
		t.Errorf("expected code STX-002, got %s", err.Code)
	}

	if len(err.Tried) != 3 {
		t.Errorf("expected 3 tried locations, got %d", len(err.Tried))
	}

	if err.Tried[0] != "Location 1" {
		t.Errorf("expected first tried location 'Location 1', got %s", err.Tried[0])
	}
}

func TestEnhancedErrorFormatting(t *testing.T) {
	err := NewEnhancedErrorWithTried(
		"STX-003",
		"Configuration not found",
		"The configuration file could not be found",
		[]string{"~/.stax/config.yml", "/etc/stax/config.yml"},
		[]Solution{
			{
				Description: "Create a configuration file",
				Command:     "stax config init",
				Steps: []string{
					"Step 1",
					"Step 2",
				},
			},
		},
		"https://example.com/docs/config",
		errors.New("file not found"),
	)

	errMsg := err.Error()

	// Check that error message contains all components
	if !strings.Contains(errMsg, "STX-003") {
		t.Errorf("error message should contain error code")
	}

	if !strings.Contains(errMsg, "Configuration not found") {
		t.Errorf("error message should contain message")
	}

	if !strings.Contains(errMsg, "The configuration file could not be found") {
		t.Errorf("error message should contain details")
	}

	if !strings.Contains(errMsg, "~/.stax/config.yml") {
		t.Errorf("error message should contain tried location 1")
	}

	if !strings.Contains(errMsg, "/etc/stax/config.yml") {
		t.Errorf("error message should contain tried location 2")
	}

	if !strings.Contains(errMsg, "Create a configuration file") {
		t.Errorf("error message should contain solution description")
	}

	if !strings.Contains(errMsg, "stax config init") {
		t.Errorf("error message should contain solution command")
	}

	if !strings.Contains(errMsg, "https://example.com/docs/config") {
		t.Errorf("error message should contain docs URL")
	}

	if !strings.Contains(errMsg, "file not found") {
		t.Errorf("error message should contain underlying error")
	}
}

func TestEnhancedErrorUnwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := NewEnhancedError(
		"STX-004",
		"Test error",
		"Details",
		[]Solution{},
		"",
		underlyingErr,
	)

	unwrapped := err.Unwrap()
	if unwrapped != underlyingErr {
		t.Errorf("Unwrap() should return underlying error")
	}
}

func TestNewConfigNotFoundError(t *testing.T) {
	err := NewConfigNotFoundError("/path/to/config.yml", errors.New("not found"))

	if err.Code != ErrCodeConfigNotFound {
		t.Errorf("expected code %s, got %s", ErrCodeConfigNotFound, err.Code)
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}

func TestNewCredentialsNotFoundError(t *testing.T) {
	tried := []string{"Keychain", "Environment", "File"}
	err := NewCredentialsNotFoundError(tried, errors.New("not found"))

	if err.Code != ErrCodeCredentialsNotFound {
		t.Errorf("expected code %s, got %s", ErrCodeCredentialsNotFound, err.Code)
	}

	if len(err.Tried) != 3 {
		t.Errorf("expected 3 tried locations, got %d", len(err.Tried))
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}

func TestNewSSHKeyNotFoundError(t *testing.T) {
	tried := []string{"~/.ssh/id_rsa", "~/.ssh/id_ed25519"}
	err := NewSSHKeyNotFoundError("/path/to/key", tried, errors.New("not found"))

	if err.Code != ErrCodeSSHKeyNotFound {
		t.Errorf("expected code %s, got %s", ErrCodeSSHKeyNotFound, err.Code)
	}

	if len(err.Tried) != 2 {
		t.Errorf("expected 2 tried locations, got %d", len(err.Tried))
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}

func TestNewDDEVNotInstalledError(t *testing.T) {
	err := NewDDEVNotInstalledError(errors.New("ddev not found"))

	if err.Code != ErrCodeDDEVNotInstalled {
		t.Errorf("expected code %s, got %s", ErrCodeDDEVNotInstalled, err.Code)
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}

func TestNewDDEVNotConfiguredError(t *testing.T) {
	err := NewDDEVNotConfiguredError("/path/to/project", errors.New("config not found"))

	if err.Code != ErrCodeDDEVNotConfigured {
		t.Errorf("expected code %s, got %s", ErrCodeDDEVNotConfigured, err.Code)
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}

func TestNewCommandNotImplementedError(t *testing.T) {
	err := NewCommandNotImplementedError(
		"test-command",
		"workaround command",
		[]string{"Step 1", "Step 2"},
	)

	if err.Code != ErrCodeCommandNotImplemented {
		t.Errorf("expected code %s, got %s", ErrCodeCommandNotImplemented, err.Code)
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}

func TestNewInvalidConfigError(t *testing.T) {
	err := NewInvalidConfigError("Invalid YAML syntax", errors.New("parse error"))

	if err.Code != ErrCodeInvalidConfig {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidConfig, err.Code)
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}

func TestNewWPEngineAPIError(t *testing.T) {
	err := NewWPEngineAPIError("API rate limit exceeded", errors.New("429 Too Many Requests"))

	if err.Code != ErrCodeWPEngineAPI {
		t.Errorf("expected code %s, got %s", ErrCodeWPEngineAPI, err.Code)
	}

	if len(err.Solutions) == 0 {
		t.Errorf("expected solutions, got none")
	}
}
