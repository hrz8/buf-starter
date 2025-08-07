package altalune

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
	"google.golang.org/protobuf/types/known/anypb"
)

// Error code constants for consistent error handling across the application
const (
	// Validation Errors
	CodeInvalidPayload = "6001"

	// Resource/Domain Errors
	CodeGreetingUnrecognized = "6101"

	// Internal Errors
	CodeInternalError   = "6901"
	CodeUnexpectedError = "6902"
)

// AppError represents a structured application error
type AppError struct {
	code     string
	message  string
	grpcCode codes.Code
	details  []proto.Message
}

// ToConnectError converts an error into a ConnectRPC-compatible error.
func ToConnectError(err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		var cErr *connect.Error

		switch appErr.grpcCode {
		case codes.InvalidArgument:
			cErr = connect.NewError(connect.CodeInvalidArgument, appErr)
		case codes.NotFound:
			cErr = connect.NewError(connect.CodeNotFound, appErr)
		case codes.AlreadyExists:
			cErr = connect.NewError(connect.CodeAlreadyExists, appErr)
		case codes.PermissionDenied:
			cErr = connect.NewError(connect.CodePermissionDenied, appErr)
		case codes.Unauthenticated:
			cErr = connect.NewError(connect.CodeUnauthenticated, appErr)
		case codes.ResourceExhausted:
			cErr = connect.NewError(connect.CodeResourceExhausted, appErr)
		case codes.FailedPrecondition:
			cErr = connect.NewError(connect.CodeFailedPrecondition, appErr)
		case codes.Unavailable:
			cErr = connect.NewError(connect.CodeUnavailable, appErr)
		case codes.DeadlineExceeded:
			cErr = connect.NewError(connect.CodeDeadlineExceeded, appErr)
		case codes.Internal:
			cErr = connect.NewError(connect.CodeInternal, appErr)
		default:
			cErr = connect.NewError(connect.CodeInternal, appErr)
		}

		for _, d := range appErr.details {
			if errDetail, err := connect.NewErrorDetail(d); err == nil {
				cErr.AddDetail(errDetail)
			}
		}

		return cErr
	}
	return err
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

// GRPCStatus converts to traditional gRPC status (for pure gRPC services)
func (e *AppError) GRPCStatus() *status.Status {
	st := status.New(e.grpcCode, e.Error())

	detailAnys := make([]protoadapt.MessageV1, 0, len(e.details))

	for _, d := range e.details {
		anyDetail, err := anypb.New(d)
		if err != nil {
			continue
		}
		detailAnys = append(detailAnys, anyDetail)
	}

	stWithDetails, err := st.WithDetails(detailAnys...)
	if err != nil {
		return st
	}

	return stWithDetails

}

func NewInvalidPayloadError(message string) *AppError {
	return &AppError{
		code:     CodeInvalidPayload,
		message:  message,
		grpcCode: codes.InvalidArgument,
	}
}

func NewGreetingUnrecognize(greeting string) *AppError {
	code := CodeGreetingUnrecognized
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Greeting to '%s' is not recognized", greeting),
		grpcCode: codes.InvalidArgument,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"name": greeting,
				},
			},
		},
	}
}

func NewInternalError(message string, err error) *AppError {
	code := CodeGreetingUnrecognized

	details := make(map[string]string)
	if err != nil {
		details["underlying_error"] = err.Error()
	}

	return &AppError{
		code:     code,
		message:  message,
		grpcCode: codes.Internal,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: details,
			},
		},
	}
}

func NewUnexpectedError(err error) *AppError {
	code := CodeUnexpectedError

	details := make(map[string]string)
	if err != nil {
		details["underlying_error"] = err.Error()
	}

	return &AppError{
		code:     code,
		message:  "An unexpected error occurred",
		grpcCode: codes.Internal,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: details,
			},
		},
	}
}
