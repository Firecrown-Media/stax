package errors

import (
	"fmt"
	"os"
)

type StaxError struct {
	Code    int
	Message string
	Cause   error
}

func (e *StaxError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *StaxError) Unwrap() error {
	return e.Cause
}

const (
	ErrCodeGeneral        = 1
	ErrCodeDDEVNotFound   = 2
	ErrCodeProjectNotFound = 3
	ErrCodeInvalidConfig  = 4
	ErrCodePermissions    = 5
)

func New(code int, message string) *StaxError {
	return &StaxError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code int, message string, cause error) *StaxError {
	return &StaxError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func DDEVNotFound() *StaxError {
	return New(ErrCodeDDEVNotFound, "DDEV is not installed or not in PATH. Please install DDEV first: https://ddev.readthedocs.io/en/stable/#installation")
}

func ProjectNotFound(path string) *StaxError {
	return New(ErrCodeProjectNotFound, fmt.Sprintf("No DDEV project found in %s. Run 'stax init' first.", path))
}

func InvalidConfig(message string) *StaxError {
	return New(ErrCodeInvalidConfig, fmt.Sprintf("Invalid configuration: %s", message))
}

func PermissionDenied(path string) *StaxError {
	return New(ErrCodePermissions, fmt.Sprintf("Permission denied accessing %s", path))
}

func HandleError(err error) {
	if err == nil {
		return
	}

	if staxErr, ok := err.(*StaxError); ok {
		fmt.Fprintf(os.Stderr, "Error: %s\n", staxErr.Message)
		if staxErr.Cause != nil && os.Getenv("STAX_VERBOSE") == "true" {
			fmt.Fprintf(os.Stderr, "Cause: %v\n", staxErr.Cause)
		}
		os.Exit(staxErr.Code)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ErrCodeGeneral)
	}
}