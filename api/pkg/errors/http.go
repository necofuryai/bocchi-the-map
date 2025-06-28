package errors

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

// HTTPErrorDetail represents additional error details for HTTP responses
type HTTPErrorDetail struct {
	Type     string                 `json:"type"`
	Resource string                 `json:"resource,omitempty"`
	Code     string                 `json:"code,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}

// Error implements the error interface for HTTPErrorDetail
func (h *HTTPErrorDetail) Error() string {
	if h.Resource != "" {
		return fmt.Sprintf("%s error for resource %s", h.Type, h.Resource)
	}
	return fmt.Sprintf("%s error", h.Type)
}

// ToHumaError converts a DomainError to a Huma error with proper status code and details
func ToHumaError(ctx context.Context, err error, operation string) error {
	// Log the error with context
	LogError(ctx, err, operation)

	var domainErr *DomainError
	if !stderrors.As(err, &domainErr) {
		// Handle non-domain errors
		return huma.Error500InternalServerError("internal server error")
	}

	// Create error detail
	detail := &HTTPErrorDetail{
		Type:     string(domainErr.Type),
		Resource: domainErr.Resource,
		Code:     domainErr.Code,
		Fields:   domainErr.Fields,
	}

	// Add request ID to error response if available
	if requestID := GetRequestID(ctx); requestID != "" {
		if detail.Fields == nil {
			detail.Fields = make(map[string]interface{})
		}
		detail.Fields["request_id"] = requestID
	}

	// Convert to appropriate HTTP status
	statusCode := domainErr.ToHTTPStatus()
	message := domainErr.Message

	switch statusCode {
	case 400:
		return huma.Error400BadRequest(message, detail)
	case 401:
		return huma.Error401Unauthorized(message, detail)
	case 403:
		return huma.Error403Forbidden(message, detail)
	case 404:
		return huma.Error404NotFound(message, detail)
	case 409:
		return huma.Error409Conflict(message, detail)
	case 410:
		return huma.Error410Gone(message, detail)
	case 422:
		return huma.Error422UnprocessableEntity(message, detail)
	case 429:
		return huma.Error429TooManyRequests(message, detail)
	case 500:
		return huma.Error500InternalServerError(message, detail)
	case 503:
		return huma.Error503ServiceUnavailable(message, detail)
	default:
		return huma.Error500InternalServerError(message, detail)
	}
}

// HandleHTTPError is a convenience function for handling errors in HTTP handlers
func HandleHTTPError(ctx context.Context, err error, operation, fallbackMessage string) error {
	if err == nil {
		return nil
	}

	// If it's already a Huma error, return as-is
	if humaErr, ok := err.(huma.StatusError); ok {
		LogErrorWithMessage(ctx, err, fmt.Sprintf("%s: %s", operation, humaErr.Error()))
		return err
	}

	// Convert domain error to HTTP error
	if Is(err, ErrTypeNotFound) ||
		Is(err, ErrTypeInvalidInput) ||
		Is(err, ErrTypeUnauthorized) ||
		Is(err, ErrTypeForbidden) ||
		Is(err, ErrTypeConflict) {
		return ToHumaError(ctx, err, operation)
	}

	// Handle unexpected errors
	LogError(ctx, err, operation)
	return huma.Error500InternalServerError(fallbackMessage)
}

// Quick error constructors for common HTTP scenarios
func HTTPNotFound(ctx context.Context, resource, id string) error {
	return ToHumaError(ctx, NotFound(resource, id), fmt.Sprintf("get_%s", resource))
}

func HTTPInvalidInput(ctx context.Context, field, reason string) error {
	return ToHumaError(ctx, InvalidInput(field, reason), "validate_input")
}

func HTTPUnauthorized(ctx context.Context, message string) error {
	return ToHumaError(ctx, Unauthorized(message), "authenticate")
}

func HTTPForbidden(ctx context.Context, resource, action string) error {
	return ToHumaError(ctx, Forbidden(resource, action), "authorize")
}

func HTTPConflict(ctx context.Context, resource, reason string) error {
	return ToHumaError(ctx, Conflict(resource, reason), fmt.Sprintf("create_%s", resource))
}

func HTTPInternal(ctx context.Context, message string) error {
	return ToHumaError(ctx, Internal(message), "internal_operation")
}