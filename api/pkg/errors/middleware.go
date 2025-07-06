package errors

import (
	"context"
	"fmt"
	"strings"

	"bocchi/api/pkg/logger"
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

// isSensitiveField checks if a field contains sensitive information
func isSensitiveField(fieldKey string) bool {
	sensitiveFields := []string{
		"password", "pass", "pwd", "secret", "key", "token", "auth",
		"credential", "cert", "private", "api_key", "access_token",
		"refresh_token", "session", "cookie", "authorization",
		"email", "phone", "ssn", "credit_card", "card_number",
		"cvv", "pin", "address", "location", "gps", "coordinates",
	}
	
	fieldLower := strings.ToLower(fieldKey)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldLower, sensitive) {
			return true
		}
	}
	return false
}

// LogError logs an error with context information
func LogError(ctx context.Context, err error, operation string) {
	fields := map[string]interface{}{
		"operation": operation,
	}

	// Add request context if available
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			fields["request_id"] = id
		}
	}
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok {
			fields["user_id"] = id
		}
	}

	// Add error-specific fields
	if domainErr, ok := err.(*DomainError); ok {
		fields["error_type"] = string(domainErr.Type)
		fields["error_code"] = domainErr.Code
		fields["resource"] = domainErr.Resource
		
		// Merge error fields with sensitive data redaction
		for k, v := range domainErr.Fields {
			if isSensitiveField(k) {
				fields[fmt.Sprintf("error_%s", k)] = "[REDACTED]"
			} else {
				fields[fmt.Sprintf("error_%s", k)] = v
			}
		}
	}

	logger.ErrorWithContextAndFields(ctx, fmt.Sprintf("Operation failed: %s", operation), err, fields)
}

// LogErrorWithMessage logs an error with a custom message and context
func LogErrorWithMessage(ctx context.Context, err error, message string) {
	fields := map[string]interface{}{}

	// Add request context if available
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			fields["request_id"] = id
		}
	}
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok {
			fields["user_id"] = id
		}
	}
	if operation := ctx.Value(OperationKey); operation != nil {
		if op, ok := operation.(string); ok {
			fields["operation"] = op
		}
	}

	// Add error-specific fields
	if domainErr, ok := err.(*DomainError); ok {
		fields["error_type"] = string(domainErr.Type)
		fields["error_code"] = domainErr.Code
		fields["resource"] = domainErr.Resource
		
		// Merge error fields with sensitive data redaction
		for k, v := range domainErr.Fields {
			if isSensitiveField(k) {
				fields[fmt.Sprintf("error_%s", k)] = "[REDACTED]"
			} else {
				fields[fmt.Sprintf("error_%s", k)] = v
			}
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