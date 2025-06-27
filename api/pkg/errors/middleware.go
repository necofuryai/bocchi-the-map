package errors

import (
	"context"
	"fmt"

	"github.com/sakai/bocchi-the-map/api/pkg/logger"
)

// ContextKey represents keys for context values
type ContextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// OperationKey is the context key for operation name
	OperationKey ContextKey = "operation"
)

// LogError logs an error with context information
func LogError(ctx context.Context, err error, operation string) {
	fields := map[string]interface{}{
		"operation": operation,
	}

	// Add request context if available
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		fields["request_id"] = requestID
	}
	if userID := ctx.Value(UserIDKey); userID != nil {
		fields["user_id"] = userID
	}

	// Add error-specific fields
	if domainErr, ok := err.(*DomainError); ok {
		fields["error_type"] = string(domainErr.Type)
		fields["error_code"] = domainErr.Code
		fields["resource"] = domainErr.Resource
		
		// Merge error fields
		for k, v := range domainErr.Fields {
			fields[fmt.Sprintf("error_%s", k)] = v
		}
	}

	logger.ErrorWithContextAndFields(ctx, fmt.Sprintf("Operation failed: %s", operation), err, fields)
}

// LogErrorWithMessage logs an error with a custom message and context
func LogErrorWithMessage(ctx context.Context, err error, message string) {
	fields := map[string]interface{}{}

	// Add request context if available
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		fields["request_id"] = requestID
	}
	if userID := ctx.Value(UserIDKey); userID != nil {
		fields["user_id"] = userID
	}
	if operation := ctx.Value(OperationKey); operation != nil {
		fields["operation"] = operation
	}

	// Add error-specific fields
	if domainErr, ok := err.(*DomainError); ok {
		fields["error_type"] = string(domainErr.Type)
		fields["error_code"] = domainErr.Code
		fields["resource"] = domainErr.Resource
		
		// Merge error fields
		for k, v := range domainErr.Fields {
			fields[fmt.Sprintf("error_%s", k)] = v
		}
	}

	logger.ErrorWithContextAndFields(ctx, message, err, fields)
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithOperation adds operation name to context
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetOperation extracts operation name from context
func GetOperation(ctx context.Context) string {
	if operation := ctx.Value(OperationKey); operation != nil {
		if op, ok := operation.(string); ok {
			return op
		}
	}
	return ""
}