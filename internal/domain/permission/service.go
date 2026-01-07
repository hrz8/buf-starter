package permission

import (
	"context"
	"regexp"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/shared/query"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// permissionNamePattern validates permission names: alphanumeric, underscores, and colons only
var permissionNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_:]+$`)

type Service struct {
	altalunev1.UnimplementedPermissionServiceServer
	validator      protovalidate.Validator
	log            altalune.Logger
	permissionRepo Repository
}

func NewService(v protovalidate.Validator, log altalune.Logger, permissionRepo Repository) *Service {
	return &Service{
		validator:      v,
		log:            log,
		permissionRepo: permissionRepo,
	}
}

func (s *Service) QueryPermissions(ctx context.Context, req *altalunev1.QueryPermissionsRequest) (*altalunev1.QueryPermissionsResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Set default query if not provided
	if req.Query == nil {
		req.Query = &altalunev1.QueryRequest{
			Pagination: &altalunev1.Pagination{
				Page:     1,
				PageSize: 10,
			},
		}
	}

	// Convert proto request to domain query params
	queryParams := query.DefaultQueryParams(req.Query)

	// Query permissions from repository
	result, err := s.permissionRepo.Query(ctx, queryParams)
	if err != nil {
		s.log.Error("failed to query permissions",
			"error", err,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query permissions: %w", err)
	}

	// Convert domain result to proto response
	if result == nil {
		return &altalunev1.QueryPermissionsResponse{
			Data: []*altalunev1.Permission{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
		}, nil
	}

	return &altalunev1.QueryPermissionsResponse{
		Data: mapPermissionsToProto(result.Data),
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   mapFiltersToProto(result.Filters),
		},
	}, nil
}

func (s *Service) CreatePermission(ctx context.Context, req *altalunev1.CreatePermissionRequest) (*altalunev1.CreatePermissionResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Trim and validate name
	name := strings.TrimSpace(req.Name)
	if !permissionNamePattern.MatchString(name) {
		return nil, altalune.NewPermissionInvalidNameError(name)
	}

	// Check if permission with same name already exists
	existingPermission, err := s.permissionRepo.GetByName(ctx, name)
	if err != nil && err != ErrPermissionNotFound {
		s.log.Error("failed to check existing permission",
			"error", err,
			"name", name,
		)
		return nil, altalune.NewUnexpectedError("failed to check existing permission: %w", err)
	}

	if existingPermission != nil {
		return nil, altalune.NewPermissionAlreadyExistsError(name)
	}

	result, err := s.permissionRepo.Create(ctx, &CreatePermissionInput{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
	})
	if err != nil {
		if err == ErrPermissionAlreadyExists {
			return nil, altalune.NewPermissionAlreadyExistsError(name)
		}
		s.log.Error("failed to create permission",
			"error", err,
			"name", name,
		)
		return nil, altalune.NewUnexpectedError("failed to create permission: %w", err)
	}

	return &altalunev1.CreatePermissionResponse{
		Permission: &altalunev1.Permission{
			Id:          result.PublicID,
			Name:        result.Name,
			Description: result.Description,
			CreatedAt:   timestamppb.New(result.CreatedAt),
			UpdatedAt:   timestamppb.New(result.UpdatedAt),
		},
		Message: "Permission created successfully",
	}, nil
}

func (s *Service) GetPermission(ctx context.Context, req *altalunev1.GetPermissionRequest) (*altalunev1.GetPermissionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	permission, err := s.permissionRepo.GetByID(ctx, req.Id)
	if err != nil {
		if err == ErrPermissionNotFound {
			return nil, altalune.NewPermissionNotFoundError(req.Id)
		}
		s.log.Error("failed to get permission", "error", err, "permission_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to get permission", err)
	}

	return &altalunev1.GetPermissionResponse{
		Permission: permission.ToPermissionProto(),
	}, nil
}

func (s *Service) UpdatePermission(ctx context.Context, req *altalunev1.UpdatePermissionRequest) (*altalunev1.UpdatePermissionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Get internal ID
	internalID, err := s.permissionRepo.GetIDByPublicID(ctx, req.Id)
	if err != nil {
		if err == ErrPermissionNotFound {
			return nil, altalune.NewPermissionNotFoundError(req.Id)
		}
		return nil, altalune.NewUnexpectedError("failed to resolve permission ID", err)
	}

	// Trim and validate name
	name := strings.TrimSpace(req.Name)
	if !permissionNamePattern.MatchString(name) {
		return nil, altalune.NewPermissionInvalidNameError(name)
	}

	input := &UpdatePermissionInput{
		ID:          internalID,
		PublicID:    req.Id,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
	}

	result, err := s.permissionRepo.Update(ctx, input)
	if err != nil {
		if err == ErrPermissionNotFound {
			return nil, altalune.NewPermissionNotFoundError(req.Id)
		}
		if err == ErrPermissionAlreadyExists {
			return nil, altalune.NewPermissionAlreadyExistsError(name)
		}
		s.log.Error("failed to update permission", "error", err, "permission_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to update permission", err)
	}

	s.log.Info("permission updated successfully", "permission_id", req.Id, "name", name)

	return &altalunev1.UpdatePermissionResponse{
		Permission: result.ToPermission().ToPermissionProto(),
		Message:    "Permission updated successfully",
	}, nil
}

func (s *Service) DeletePermission(ctx context.Context, req *altalunev1.DeletePermissionRequest) (*altalunev1.DeletePermissionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// First get the permission to check if it's protected
	permission, err := s.permissionRepo.GetByID(ctx, req.Id)
	if err != nil {
		if err == ErrPermissionNotFound {
			return nil, altalune.NewPermissionNotFoundError(req.Id)
		}
		s.log.Error("failed to get permission for deletion check", "error", err, "permission_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to get permission", err)
	}

	// Check if permission is protected
	if permission.Name == "root" {
		s.log.Warn("attempt to delete protected permission", "permission_name", permission.Name, "permission_id", req.Id)
		return nil, altalune.NewPermissionProtectedError(permission.Name)
	}

	err = s.permissionRepo.Delete(ctx, req.Id)
	if err != nil {
		if err == ErrPermissionNotFound {
			return nil, altalune.NewPermissionNotFoundError(req.Id)
		}
		if err == ErrPermissionInUse {
			return nil, altalune.NewPermissionInUseError(req.Id)
		}
		if err == ErrPermissionProtected {
			return nil, altalune.NewPermissionProtectedError(permission.Name)
		}
		s.log.Error("failed to delete permission", "error", err, "permission_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to delete permission", err)
	}

	s.log.Info("permission deleted successfully", "permission_id", req.Id)

	return &altalunev1.DeletePermissionResponse{
		Message: "Permission deleted successfully",
	}, nil
}
