package starter

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error code constants for consistent error handling across the application
const (
	// Validation Errors (INVALID_ARGUMENT)
	CodeInvalidPayload = "INVALID_PAYLOAD"

	// Resource Errors
	CodeGreetingUnrecognize = "GREETING_UNRECOGNIZE"

	// Internal Errors (INTERNAL)
	CodeInternalError   = "INTERNAL_ERROR"
	CodeUnexpectedError = "UNEXPECTED_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Code     string
	Message  string
	GRPCCode codes.Code
	Details  map[string]any
}

// ToConnectError converts an error into a ConnectRPC-compatible error.
func ToConnectError(err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.httpErr()
	}
	return err
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// httpErr converts the AppError to a gRPC status error compatible with both gRPC and ConnectRPC
func (e *AppError) httpErr() error {
	switch e.GRPCCode {
	case codes.InvalidArgument:
		return connect.NewError(connect.CodeInvalidArgument, e)
	case codes.NotFound:
		return connect.NewError(connect.CodeNotFound, e)
	case codes.AlreadyExists:
		return connect.NewError(connect.CodeAlreadyExists, e)
	case codes.PermissionDenied:
		return connect.NewError(connect.CodePermissionDenied, e)
	case codes.Unauthenticated:
		return connect.NewError(connect.CodeUnauthenticated, e)
	case codes.ResourceExhausted:
		return connect.NewError(connect.CodeResourceExhausted, e)
	case codes.FailedPrecondition:
		return connect.NewError(connect.CodeFailedPrecondition, e)
	case codes.Unavailable:
		return connect.NewError(connect.CodeUnavailable, e)
	case codes.DeadlineExceeded:
		return connect.NewError(connect.CodeDeadlineExceeded, e)
	case codes.Internal:
		return connect.NewError(connect.CodeInternal, e)
	default:
		return connect.NewError(connect.CodeInternal, e)
	}
}

// GRPCStatus converts to traditional gRPC status (for pure gRPC services)
func (e *AppError) GRPCStatus() *status.Status {
	return status.New(e.GRPCCode, e.Error())
}

// Validation Errors
func NewInvalidPayloadError(message string) *AppError {
	return &AppError{
		Code:     CodeInvalidPayload,
		Message:  message,
		GRPCCode: codes.InvalidArgument,
	}
}

// Resource Errors
func NewGreetingUnrecognize(greeting string) *AppError {
	return &AppError{
		Code:     CodeGreetingUnrecognize,
		Message:  fmt.Sprintf("Greeting to '%s' is not recognized", greeting),
		GRPCCode: codes.InvalidArgument,
		Details: map[string]any{
			"name": greeting,
		},
	}
}

// Internal Errors
func NewInternalError(message string, err error) *AppError {
	details := make(map[string]any)
	if err != nil {
		details["underlying_error"] = err.Error()
	}

	return &AppError{
		Code:     CodeInternalError,
		Message:  message,
		GRPCCode: codes.Internal,
		Details:  details,
	}
}

func NewUnexpectedError(err error) *AppError {
	details := make(map[string]any)
	if err != nil {
		details["underlying_error"] = err.Error()
	}

	return &AppError{
		Code:     CodeUnexpectedError,
		Message:  "An unexpected error occurred",
		GRPCCode: codes.Internal,
		Details:  details,
	}
}

// Helper function to check if an error is a specific type
func IsErrorCode(err error, code string) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// Helper function to extract error code from any error
func GetErrorCode(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return CodeUnexpectedError
}
