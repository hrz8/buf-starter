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
	// Validation Errors (600XX)
	CodeInvalidPayload = "60001"

	// Greeting Domain Errors (601XX)
	CodeGreetingUnrecognized = "60101"

	// Example/Employee Domain Errors (602XX)
	CodeEmployeeNotFound      = "60201"
	CodeEmployeeAlreadyExists = "60202"

	// Project Domain Errors (603XX)
	CodeProjectNotFound = "60301"

	// API Key Domain Errors (604XX)
	CodeApiKeyNotFound      = "60401"
	CodeApiKeyAlreadyExists = "60402"

	// User Domain Errors (605XX)
	CodeUserNotFound         = "60500"
	CodeUserAlreadyExists    = "60501"
	CodeUserInvalidEmail     = "60502"
	CodeUserAlreadyActive    = "60503"
	CodeUserAlreadyInactive  = "60504"
	CodeUserCannotDeleteSelf = "60505"

	// Role Domain Errors (606XX)
	CodeRoleNotFound      = "60600"
	CodeRoleAlreadyExists = "60601"
	CodeRoleInvalidName   = "60602"
	CodeRoleInUse         = "60603"
	CodeRoleProtected     = "60604"

	// Permission Domain Errors (607XX)
	CodePermissionNotFound      = "60700"
	CodePermissionAlreadyExists = "60701"
	CodePermissionInvalidName   = "60702"
	CodePermissionInUse         = "60704"
	CodePermissionProtected     = "60705"

	// IAM Mapper Domain Errors (608XX)
	CodeMappingNotFound           = "60800"
	CodeMappingAlreadyExists      = "60801"
	CodeInvalidProjectRole        = "60802"
	CodeCannotRemoveLastOwner     = "60803"
	CodeMappingUserNotFound       = "60804"
	CodeMappingRoleNotFound       = "60805"
	CodeMappingPermissionNotFound = "60806"
	CodeMappingProjectNotFound    = "60807"

	// OAuth Provider Domain Errors (608XX continued)
	CodeOAuthProviderNotFound        = "60810"
	CodeOAuthProviderDuplicateType   = "60811"
	CodeOAuthProviderEncryptionError = "60812"
	CodeOAuthProviderDecryptionError = "60813"

	// OAuth Client Domain Errors (609XX)
	CodeOAuthClientNotFound      = "60900"
	CodeOAuthClientAlreadyExists = "60901"
	CodeInvalidRedirectURI       = "60902"
	CodeOAuthClientSecretInvalid = "60903"

	// Internal Errors (699XX)
	CodeUnexpectedError = "69901"
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

// common
func NewInvalidPayloadError(message string) *AppError {
	return &AppError{
		code:     CodeInvalidPayload,
		message:  message,
		grpcCode: codes.InvalidArgument,
	}
}

func NewUnexpectedError(message string, err error) *AppError {
	code := CodeUnexpectedError

	details := make(map[string]string)
	if err != nil {
		// TODO: should exclude expose runtime error in prod env
		details["underlying_error"] = message + ": " + err.Error()
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

func NewProjectNotFound(projectID string) *AppError {
	code := CodeProjectNotFound
	return &AppError{
		code:     code,
		message:  "Project not found",
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"project_id": projectID,
				},
			},
		},
	}
}

// domain-based
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

// NewEmployeeNotFoundError creates an error for when an employee is not found
func NewEmployeeNotFoundError(publicID string) *AppError {
	code := CodeEmployeeNotFound
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Employee with ID '%s' not found", publicID),
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"employee_id": publicID,
				},
			},
		},
	}
}

// NewAlreadyExistsError creates a new already exists error
func NewAlreadyExistsError(email string) *AppError {
	code := CodeEmployeeAlreadyExists
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("employee with email '%s' already exists", email),
		grpcCode: codes.AlreadyExists,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"email": email,
				},
			},
		},
	}
}

// NewApiKeyNotFoundError creates an error for when an API key is not found
func NewApiKeyNotFoundError(publicID string) *AppError {
	code := CodeApiKeyNotFound
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("API key with ID '%s' not found", publicID),
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"api_key_id": publicID,
				},
			},
		},
	}
}

// NewApiKeyAlreadyExistsError creates a new API key already exists error
func NewApiKeyAlreadyExistsError(name string) *AppError {
	code := CodeApiKeyAlreadyExists
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("API key with name '%s' already exists", name),
		grpcCode: codes.AlreadyExists,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"name": name,
				},
			},
		},
	}
}

// User Domain Errors

// NewUserNotFoundError creates an error for when a user is not found
func NewUserNotFoundError(userID string) *AppError {
	code := CodeUserNotFound
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("User with ID '%s' not found", userID),
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"user_id": userID,
				},
			},
		},
	}
}

// NewUserAlreadyExistsError creates a new user already exists error
func NewUserAlreadyExistsError(email string) *AppError {
	code := CodeUserAlreadyExists
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("User with email '%s' already exists", email),
		grpcCode: codes.AlreadyExists,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"email": email,
				},
			},
		},
	}
}

// NewUserInvalidEmailError creates an error for invalid email
func NewUserInvalidEmailError(email string) *AppError {
	code := CodeUserInvalidEmail
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Invalid email format: '%s'", email),
		grpcCode: codes.InvalidArgument,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"email": email,
				},
			},
		},
	}
}

// NewUserAlreadyActiveError creates an error when user is already active
func NewUserAlreadyActiveError(userID string) *AppError {
	code := CodeUserAlreadyActive
	return &AppError{
		code:     code,
		message:  "User is already active",
		grpcCode: codes.FailedPrecondition,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"user_id": userID,
				},
			},
		},
	}
}

// NewUserAlreadyInactiveError creates an error when user is already inactive
func NewUserAlreadyInactiveError(userID string) *AppError {
	code := CodeUserAlreadyInactive
	return &AppError{
		code:     code,
		message:  "User is already inactive",
		grpcCode: codes.FailedPrecondition,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"user_id": userID,
				},
			},
		},
	}
}

// NewUserCannotDeleteSelfError creates an error when user tries to delete themselves
func NewUserCannotDeleteSelfError(userID string) *AppError {
	code := CodeUserCannotDeleteSelf
	return &AppError{
		code:     code,
		message:  "Cannot delete your own user account",
		grpcCode: codes.PermissionDenied,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"user_id": userID,
				},
			},
		},
	}
}

// Role Domain Errors

// NewRoleNotFoundError creates an error for when a role is not found
func NewRoleNotFoundError(roleID string) *AppError {
	code := CodeRoleNotFound
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Role with ID '%s' not found", roleID),
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"role_id": roleID,
				},
			},
		},
	}
}

// NewRoleAlreadyExistsError creates a new role already exists error
func NewRoleAlreadyExistsError(name string) *AppError {
	code := CodeRoleAlreadyExists
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Role with name '%s' already exists", name),
		grpcCode: codes.AlreadyExists,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"name": name,
				},
			},
		},
	}
}

// NewRoleInvalidNameError creates an error for invalid role name
func NewRoleInvalidNameError(name string) *AppError {
	code := CodeRoleInvalidName
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Invalid role name: '%s'", name),
		grpcCode: codes.InvalidArgument,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"name": name,
				},
			},
		},
	}
}

// NewRoleInUseError creates an error when role cannot be deleted because it's in use
func NewRoleInUseError(roleID string) *AppError {
	code := CodeRoleInUse
	return &AppError{
		code:     code,
		message:  "Role is in use and cannot be deleted",
		grpcCode: codes.FailedPrecondition,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"role_id": roleID,
				},
			},
		},
	}
}

// NewRoleProtectedError creates an error when trying to delete a protected role
func NewRoleProtectedError(roleName string) *AppError {
	code := CodeRoleProtected
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Role '%s' is protected and cannot be deleted or modified", roleName),
		grpcCode: codes.PermissionDenied,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"role_name": roleName,
				},
			},
		},
	}
}

// Permission Domain Errors

// NewPermissionNotFoundError creates an error for when a permission is not found
func NewPermissionNotFoundError(permissionID string) *AppError {
	code := CodePermissionNotFound
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Permission with ID '%s' not found", permissionID),
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"permission_id": permissionID,
				},
			},
		},
	}
}

// NewPermissionAlreadyExistsError creates a new permission already exists error
func NewPermissionAlreadyExistsError(name string) *AppError {
	code := CodePermissionAlreadyExists
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Permission with name '%s' already exists", name),
		grpcCode: codes.AlreadyExists,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"name": name,
				},
			},
		},
	}
}

// NewPermissionInvalidNameError creates an error for invalid permission name
func NewPermissionInvalidNameError(name string) *AppError {
	code := CodePermissionInvalidName
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Invalid permission name: '%s' (must match ^[a-zA-Z0-9_:]+$)", name),
		grpcCode: codes.InvalidArgument,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"name": name,
				},
			},
		},
	}
}

// NewPermissionInUseError creates an error when permission cannot be deleted because it's in use
func NewPermissionInUseError(permissionID string) *AppError {
	code := CodePermissionInUse
	return &AppError{
		code:     code,
		message:  "Permission is in use and cannot be deleted",
		grpcCode: codes.FailedPrecondition,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"permission_id": permissionID,
				},
			},
		},
	}
}

// NewPermissionProtectedError creates an error when trying to delete a protected permission
func NewPermissionProtectedError(permissionName string) *AppError {
	code := CodePermissionProtected
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Permission '%s' is protected and cannot be deleted or modified", permissionName),
		grpcCode: codes.PermissionDenied,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"permission_name": permissionName,
				},
			},
		},
	}
}

// IAM Mapper Domain Errors

// NewInvalidProjectRoleError creates an error for invalid project role
func NewInvalidProjectRoleError(role string) *AppError {
	code := CodeInvalidProjectRole
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("Invalid project role: '%s' (must be owner, admin, member, or viewer)", role),
		grpcCode: codes.InvalidArgument,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"role": role,
				},
			},
		},
	}
}

// NewCannotRemoveLastOwnerError creates an error when trying to remove the last owner
func NewCannotRemoveLastOwnerError(projectID string) *AppError {
	code := CodeCannotRemoveLastOwner
	return &AppError{
		code:     code,
		message:  "Cannot remove the last owner from the project",
		grpcCode: codes.FailedPrecondition,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"project_id": projectID,
				},
			},
		},
	}
}

// NewOAuthProviderNotFoundError creates an error for when an OAuth provider is not found
func NewOAuthProviderNotFoundError(providerID string) *AppError {
	code := CodeOAuthProviderNotFound
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("OAuth provider with ID '%s' not found", providerID),
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"provider_id": providerID,
				},
			},
		},
	}
}

// NewOAuthProviderDuplicateTypeError creates an error for duplicate provider type
func NewOAuthProviderDuplicateTypeError(providerType string) *AppError {
	code := CodeOAuthProviderDuplicateType
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("OAuth provider with type '%s' already exists", providerType),
		grpcCode: codes.AlreadyExists,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"provider_type": providerType,
				},
			},
		},
	}
}

// NewOAuthProviderEncryptionError creates an error for client secret encryption failure
func NewOAuthProviderEncryptionError(providerID string) *AppError {
	code := CodeOAuthProviderEncryptionError
	return &AppError{
		code:     code,
		message:  "Failed to encrypt OAuth provider client secret",
		grpcCode: codes.Internal,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"provider_id": providerID,
				},
			},
		},
	}
}

// NewOAuthProviderDecryptionError creates an error for client secret decryption failure
func NewOAuthProviderDecryptionError(providerID string) *AppError {
	code := CodeOAuthProviderDecryptionError
	return &AppError{
		code:     code,
		message:  "Failed to decrypt OAuth provider client secret",
		grpcCode: codes.Internal,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"provider_id": providerID,
				},
			},
		},
	}
}

// ==================== OAuth Client Domain Errors ====================

// NewOAuthClientNotFoundError creates an error for OAuth client not found
func NewOAuthClientNotFoundError(clientID string) *AppError {
	code := CodeOAuthClientNotFound
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("OAuth client with ID '%s' not found", clientID),
		grpcCode: codes.NotFound,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"client_id": clientID,
				},
			},
		},
	}
}

// NewOAuthClientAlreadyExistsError creates an error for duplicate OAuth client name
func NewOAuthClientAlreadyExistsError(name string) *AppError {
	code := CodeOAuthClientAlreadyExists
	return &AppError{
		code:     code,
		message:  fmt.Sprintf("OAuth client with name '%s' already exists", name),
		grpcCode: codes.AlreadyExists,
		details: []proto.Message{
			&altalunev1.ErrorDetail{
				Code: code,
				Meta: map[string]string{
					"name": name,
				},
			},
		},
	}
}
