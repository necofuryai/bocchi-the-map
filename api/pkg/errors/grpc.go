package errors

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError converts a DomainError to a gRPC error with proper status code
func ToGRPCError(ctx context.Context, err error, operation string) error {
	// Log the error with context
	LogError(ctx, err, operation)

	var domainErr *DomainError
	if !stderrors.As(err, &domainErr) {
		// Handle common standard errors
		if err == sql.ErrNoRows {
			return status.Error(codes.NotFound, "resource not found")
		}
		
		// Handle unexpected errors
		return status.Error(codes.Internal, "internal server error")
	}

	return domainErr.ToGRPCError()
}

// HandleGRPCError is a convenience function for handling errors in gRPC services
func HandleGRPCError(ctx context.Context, err error, operation, fallbackMessage string) error {
	if err == nil {
		return nil
	}

	// If it's already a gRPC status error, log and return as-is
	if _, ok := status.FromError(err); ok {
		LogErrorWithMessage(ctx, err, fmt.Sprintf("%s: gRPC error", operation))
		return err
	}

	// Handle common database errors
	if err == sql.ErrNoRows {
		LogError(ctx, NotFound("resource", "unknown"), operation)
		return status.Error(codes.NotFound, "resource not found")
	}

	// Convert domain error to gRPC error
	return ToGRPCError(ctx, err, operation)
}

// Common gRPC error constructors
func GRPCNotFound(ctx context.Context, resource, id string) error {
	err := NotFound(resource, id)
	LogError(ctx, err, fmt.Sprintf("get_%s", resource))
	return err.ToGRPCError()
}

func GRPCInvalidArgument(ctx context.Context, field, reason string) error {
	err := InvalidInput(field, reason)
	LogError(ctx, err, "validate_input")
	return err.ToGRPCError()
}

func GRPCUnauthenticated(ctx context.Context, message string) error {
	err := Unauthorized(message)
	LogError(ctx, err, "authenticate")
	return err.ToGRPCError()
}

func GRPCPermissionDenied(ctx context.Context, resource, action string) error {
	err := Forbidden(resource, action)
	LogError(ctx, err, "authorize")
	return err.ToGRPCError()
}

func GRPCAlreadyExists(ctx context.Context, resource, reason string) error {
	err := Conflict(resource, reason)
	LogError(ctx, err, fmt.Sprintf("create_%s", resource))
	return err.ToGRPCError()
}

func GRPCInternal(ctx context.Context, message string) error {
	err := Internal(message)
	LogError(ctx, err, "internal_operation")
	return err.ToGRPCError()
}

// Database error handlers
func HandleDatabaseError(ctx context.Context, err error, operation string) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return GRPCNotFound(ctx, "resource", "unknown")
	}

	// Wrap database errors
	dbErr := Database(operation, err)
	LogError(ctx, dbErr, operation)
	return dbErr.ToGRPCError()
}

// Background operation error handler that doesn't return but logs
func HandleBackgroundError(ctx context.Context, err error, operation string) {
	if err == nil {
		return
	}

	// Add background operation context
	bgCtx := WithOperation(ctx, fmt.Sprintf("background_%s", operation))
	
	var domainErr *DomainError
	if !stderrors.As(err, &domainErr) {
		domainErr = Wrap(err, ErrTypeInternal, fmt.Sprintf("background operation failed: %s: %s", operation, err.Error()))
	}

	LogError(bgCtx, domainErr, operation)
}