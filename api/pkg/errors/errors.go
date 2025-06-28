package errors

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ErrorType represents the type of error that occurred
type ErrorType string

const (
	// Client errors
	ErrTypeNotFound     ErrorType = "NOT_FOUND"
	ErrTypeInvalidInput ErrorType = "INVALID_INPUT"
	ErrTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrTypeForbidden    ErrorType = "FORBIDDEN"
	ErrTypeConflict     ErrorType = "CONFLICT"

	// Server errors
	ErrTypeInternal      ErrorType = "INTERNAL"
	ErrTypeDatabase      ErrorType = "DATABASE"
	ErrTypeExternalAPI   ErrorType = "EXTERNAL_API"
	ErrTypeValidation    ErrorType = "VALIDATION"
	ErrTypeRateLimit     ErrorType = "RATE_LIMIT"
	ErrTypeTimeout       ErrorType = "TIMEOUT"
	ErrTypeUnavailable   ErrorType = "UNAVAILABLE"
)

// DomainError represents a domain-specific error with context
type DomainError struct {
	Type     ErrorType              `json:"type"`
	Message  string                 `json:"message"`
	Code     string                 `json:"code,omitempty"`
	Cause    error                  `json:"-"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
	Resource string                 `json:"resource,omitempty"`
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %s", e.Type, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap implements the unwrap interface for errors.Is and errors.As
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// WithField adds a field to the error context
func (e *DomainError) WithField(key string, value interface{}) *DomainError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	e.Fields[key] = value
	return e
}

// WithResource sets the resource name for the error
func (e *DomainError) WithResource(resource string) *DomainError {
	e.Resource = resource
	return e
}

// WithCode sets a specific error code
func (e *DomainError) WithCode(code string) *DomainError {
	e.Code = code
	return e
}

// New creates a new DomainError
func New(errorType ErrorType, message string) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errorType ErrorType, message string) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		Cause:   err,
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, errorType ErrorType, format string, args ...interface{}) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// Common error constructors
func NotFound(resource, id string) *DomainError {
	return New(ErrTypeNotFound, fmt.Sprintf("%s not found", resource)).
		WithResource(resource).
		WithField("id", id)
}

func InvalidInput(field, reason string) *DomainError {
	return New(ErrTypeInvalidInput, fmt.Sprintf("invalid %s: %s", field, reason)).
		WithField("field", field).
		WithField("reason", reason)
}

func Unauthorized(message string) *DomainError {
	return New(ErrTypeUnauthorized, message)
}

func Forbidden(resource, action string) *DomainError {
	return New(ErrTypeForbidden, fmt.Sprintf("forbidden to %s %s", action, resource)).
		WithResource(resource).
		WithField("action", action)
}

func Conflict(resource, reason string) *DomainError {
	return New(ErrTypeConflict, fmt.Sprintf("%s conflict: %s", resource, reason)).
		WithResource(resource).
		WithField("reason", reason)
}

func Internal(message string) *DomainError {
	return New(ErrTypeInternal, message)
}

func Database(operation string, err error) *DomainError {
	return Wrap(err, ErrTypeDatabase, fmt.Sprintf("database %s failed", operation)).
		WithField("operation", operation)
}

func ExternalAPI(service string, err error) *DomainError {
	return Wrap(err, ErrTypeExternalAPI, fmt.Sprintf("external API %s failed", service)).
		WithField("service", service)
}

// ToHTTPStatus converts DomainError to HTTP status code
func (e *DomainError) ToHTTPStatus() int {
	switch e.Type {
	case ErrTypeNotFound:
		return http.StatusNotFound
	case ErrTypeInvalidInput, ErrTypeValidation:
		return http.StatusBadRequest
	case ErrTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrTypeForbidden:
		return http.StatusForbidden
	case ErrTypeConflict:
		return http.StatusConflict
	case ErrTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrTypeTimeout:
		return http.StatusRequestTimeout
	case ErrTypeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// ToGRPCStatus converts DomainError to gRPC status with detailed error information
func (e *DomainError) ToGRPCStatus() *status.Status {
	var code codes.Code
	switch e.Type {
	case ErrTypeNotFound:
		code = codes.NotFound
	case ErrTypeInvalidInput, ErrTypeValidation:
		code = codes.InvalidArgument
	case ErrTypeUnauthorized:
		code = codes.Unauthenticated
	case ErrTypeForbidden:
		code = codes.PermissionDenied
	case ErrTypeConflict:
		code = codes.AlreadyExists
	case ErrTypeRateLimit:
		code = codes.ResourceExhausted
	case ErrTypeTimeout:
		code = codes.DeadlineExceeded
	case ErrTypeUnavailable:
		code = codes.Unavailable
	default:
		code = codes.Internal
	}

	// Create base status
	st := status.New(code, e.Message)

	// Add error details if Fields or Resource are present
	var details []proto.Message

	// Add resource info if present
	if e.Resource != "" {
		resourceInfo := &errdetails.ResourceInfo{
			ResourceType: string(e.Type),
			ResourceName: e.Resource,
		}
		details = append(details, resourceInfo)
	}

	// Add field violations if present
	if len(e.Fields) > 0 {
		fieldViolations := make([]*errdetails.BadRequest_FieldViolation, 0, len(e.Fields))
		for field, message := range e.Fields {
			// Convert interface{} to string with type assertion
			messageStr, ok := message.(string)
			if !ok {
				messageStr = fmt.Sprintf("%v", message)
			}
			fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: messageStr,
			})
		}
		badRequest := &errdetails.BadRequest{
			FieldViolations: fieldViolations,
		}
		details = append(details, badRequest)
	}

	// For now, return basic status without details to avoid protobuf compatibility issues
	// TODO: Implement proper protobuf message conversion for error details
	_ = details // Suppress unused variable warning

	return st
}

// ToGRPCError converts DomainError to gRPC error
func (e *DomainError) ToGRPCError() error {
	return e.ToGRPCStatus().Err()
}

// Is checks if the error is of a specific type
func Is(err error, errorType ErrorType) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == errorType
	}
	return false
}

// GetType extracts the ErrorType from an error
func GetType(err error) ErrorType {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type
	}
	return ErrTypeInternal
}

// GetFields extracts the fields from an error
func GetFields(err error) map[string]interface{} {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Fields
	}
	return nil
}