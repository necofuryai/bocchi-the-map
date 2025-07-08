package handlers

import (
	"github.com/danielgtaylor/huma/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpcToHTTPError converts gRPC errors to appropriate HTTP error responses
func grpcToHTTPError(err error, defaultMessage string) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return huma.Error500InternalServerError(defaultMessage)
	}

	switch st.Code() {
	case codes.NotFound:
		return huma.Error404NotFound(st.Message())
	case codes.InvalidArgument:
		return huma.Error400BadRequest(st.Message())
	case codes.AlreadyExists:
		return huma.Error409Conflict(st.Message())
	case codes.PermissionDenied:
		return huma.Error403Forbidden(st.Message())
	case codes.Unauthenticated:
		return huma.Error401Unauthorized(st.Message())
	case codes.FailedPrecondition:
		return huma.Error412PreconditionFailed(st.Message())
	case codes.OutOfRange:
		return huma.Error400BadRequest(st.Message())
	case codes.Unimplemented:
		return huma.Error501NotImplemented(st.Message())
	case codes.Unavailable:
		return huma.Error503ServiceUnavailable(st.Message())
	case codes.DeadlineExceeded:
		return huma.Error503ServiceUnavailable(st.Message())
	default:
		return huma.Error500InternalServerError(defaultMessage)
	}
}