package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common cases
var (
	ErrEmptyPath           = errors.New("path cannot be empty")
	ErrPathTraversal       = errors.New("path traversal not allowed")
	ErrPathTooLong         = errors.New("path too long")
	ErrFileNotExists       = errors.New("file does not exist")
	ErrNotReadable         = errors.New("file is not readable")
	ErrNotWritable         = errors.New("directory is not writable")
	ErrNotDirectory        = errors.New("path is not a directory")
	ErrNotFile             = errors.New("path is not a file")
	ErrNotMarkdown         = errors.New("file is not a markdown file")
	ErrFileTooLarge        = errors.New("file too large")
	ErrBinaryContent       = errors.New("file appears to be binary")
	ErrInvalidName         = errors.New("invalid name")
	ErrNameTooLong         = errors.New("name too long")
	ErrTemplateNotFound    = errors.New("template not found")
	ErrUnsupportedTemplate = errors.New("unsupported template")
	ErrInvalidCategory     = errors.New("invalid category")
)

// ValidationError represents a validation error with additional context
type ValidationError struct {
	Type    string
	Field   string
	Value   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s validation failed for %s '%s': %s (%v)",
			e.Type, e.Field, e.Value, e.Message, e.Err)
	}
	return fmt.Sprintf("%s validation failed for %s '%s': %s",
		e.Type, e.Field, e.Value, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(typ, field, value, message string, err error) *ValidationError {
	return &ValidationError{
		Type:    typ,
		Field:   field,
		Value:   value,
		Message: message,
		Err:     err,
	}
}

// OperationError represents an error during an operation with context
type OperationError struct {
	Operation string
	Target    string
	Err       error
}

func (e *OperationError) Error() string {
	return fmt.Sprintf("failed to %s %s: %v", e.Operation, e.Target, e.Err)
}

func (e *OperationError) Unwrap() error {
	return e.Err
}

// NewOperationError creates a new operation error
func NewOperationError(operation, target string, err error) *OperationError {
	return &OperationError{
		Operation: operation,
		Target:    target,
		Err:       err,
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Component string
	Issue     string
	Err       error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s configuration error: %s (%v)", e.Component, e.Issue, e.Err)
	}
	return fmt.Sprintf("%s configuration error: %s", e.Component, e.Issue)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new configuration error
func NewConfigError(component, issue string, err error) *ConfigError {
	return &ConfigError{
		Component: component,
		Issue:     issue,
		Err:       err,
	}
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	return fmt.Errorf(format+": %w", append(args, err)...)
}

// WithContext adds context to an error
func WithContext(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}
