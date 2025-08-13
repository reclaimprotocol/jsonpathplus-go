package jsonpathplus

import (
	"fmt"
	"strings"
)

// ErrorType represents different types of JSONPath errors.
type ErrorType int

const (
	// ErrInvalidPath indicates an invalid JSONPath expression.
	ErrInvalidPath ErrorType = iota
	// ErrParseError indicates a parsing error in the JSONPath expression.
	ErrParseError
	// ErrEvaluationError indicates an error during evaluation of the JSONPath.
	ErrEvaluationError
	// ErrInvalidJSON indicates invalid JSON input.
	ErrInvalidJSON
	// ErrInvalidExpression indicates an invalid filter expression.
	ErrInvalidExpression
	// ErrOutOfBounds indicates array index out of bounds.
	ErrOutOfBounds
	// ErrTypeError indicates a type mismatch error.
	ErrTypeError
	// ErrRecursionLimit indicates recursion depth limit exceeded.
	ErrRecursionLimit
)

// JSONPathError represents an error that occurred during JSONPath operations.
type JSONPathError struct {
	Type     ErrorType
	Message  string
	Path     string
	Position int
	Cause    error
}

func (e *JSONPathError) Error() string {
	var parts []string

	switch e.Type {
	case ErrInvalidPath:
		parts = append(parts, "invalid JSONPath")
	case ErrParseError:
		parts = append(parts, "parse error")
	case ErrEvaluationError:
		parts = append(parts, "evaluation error")
	case ErrInvalidJSON:
		parts = append(parts, "invalid JSON")
	case ErrInvalidExpression:
		parts = append(parts, "invalid expression")
	case ErrOutOfBounds:
		parts = append(parts, "index out of bounds")
	case ErrTypeError:
		parts = append(parts, "type error")
	case ErrRecursionLimit:
		parts = append(parts, "recursion limit exceeded")
	}

	if e.Path != "" {
		parts = append(parts, fmt.Sprintf("path: %s", e.Path))
	}

	if e.Position >= 0 {
		parts = append(parts, fmt.Sprintf("position: %d", e.Position))
	}

	if e.Message != "" {
		parts = append(parts, e.Message)
	}

	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Cause))
	}

	return strings.Join(parts, ": ")
}

func (e *JSONPathError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target type.
func (e *JSONPathError) Is(target error) bool {
	if t, ok := target.(*JSONPathError); ok {
		return e.Type == t.Type
	}
	return false
}

// NewError creates a new JSONPathError.
func NewError(errType ErrorType, message string, path string, position int) *JSONPathError {
	return &JSONPathError{
		Type:     errType,
		Message:  message,
		Path:     path,
		Position: position,
	}
}

// WrapError wraps an existing error with JSONPath context.
func WrapError(errType ErrorType, cause error, path string, position int) *JSONPathError {
	return &JSONPathError{
		Type:     errType,
		Cause:    cause,
		Path:     path,
		Position: position,
	}
}

// ValidationError represents validation errors.
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s' with value '%v': %s", e.Field, e.Value, e.Message)
}

// PathLengthError represents errors when path is too long.
type PathLengthError struct {
	Length int
	Limit  int
}

func (e *PathLengthError) Error() string {
	return fmt.Sprintf("path length %d exceeds limit of %d", e.Length, e.Limit)
}

// RecursionLimitError represents errors when recursion depth is exceeded.
type RecursionLimitError struct {
	Depth int
	Limit int
}

func (e *RecursionLimitError) Error() string {
	return fmt.Sprintf("recursion depth %d exceeds limit of %d", e.Depth, e.Limit)
}
