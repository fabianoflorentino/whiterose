package errors

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// ErrorType represents the type of domain error
type ErrorType string

const (
	// ValidationError indicates invalid input or data
	ValidationError ErrorType = "VALIDATION_ERROR"
	// BusinessRuleError indicates violation of business rules
	BusinessRuleError ErrorType = "BUSINESS_RULE_ERROR"
	// NotFoundError indicates requested resource was not found
	NotFoundError ErrorType = "NOT_FOUND_ERROR"
	// ConflictError indicates resource conflict
	ConflictError ErrorType = "CONFLICT_ERROR"
)

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Is checks if the error is of a specific type
func (e *DomainError) Is(target error) bool {
	if other, ok := target.(*DomainError); ok {
		return e.Type == other.Type
	}
	return false
}

// NewValidationError creates a new validation error
func NewValidationError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ValidationError,
		Message: message,
		Cause:   cause,
	}
}

// NewBusinessRuleError creates a new business rule error
func NewBusinessRuleError(message string) *DomainError {
	return &DomainError{
		Type:    BusinessRuleError,
		Message: message,
		Cause:   nil,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *DomainError {
	return &DomainError{
		Type:    NotFoundError,
		Message: message,
		Cause:   nil,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *DomainError {
	return &DomainError{
		Type:    ConflictError,
		Message: message,
		Cause:   nil,
	}
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ValidationError
	}
	return false
}

// IsBusinessRuleError checks if error is a business rule error
func IsBusinessRuleError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == BusinessRuleError
	}
	return false
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == NotFoundError
	}
	return false
}

// IsConflictError checks if error is a conflict error
func IsConflictError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ConflictError
	}
	return false
}
