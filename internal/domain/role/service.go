package role

import (
	"context"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/shared/query"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	altalunev1.UnimplementedRoleServiceServer
	validator protovalidate.Validator
	log       altalune.Logger
	roleRepo  Repository
}

func NewService(v protovalidate.Validator, log altalune.Logger, roleRepo Repository) *Service {
	return &Service{
		validator: v,
		log:       log,
		roleRepo:  roleRepo,
	}
}

func (s *Service) QueryRoles(ctx context.Context, req *altalunev1.QueryRolesRequest) (*altalunev1.QueryRolesResponse, error) {
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

	// Query roles from repository
	result, err := s.roleRepo.Query(ctx, queryParams)
	if err != nil {
		s.log.Error("failed to query roles",
			"error", err,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query roles: %w", err)
	}

	// Convert domain result to proto response
	if result == nil {
		return &altalunev1.QueryRolesResponse{
			Data: []*altalunev1.Role{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
		}, nil
	}

	return &altalunev1.QueryRolesResponse{
		Data: mapRolesToProto(result.Data),
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   mapFiltersToProto(result.Filters),
		},
	}, nil
}

func (s *Service) CreateRole(ctx context.Context, req *altalunev1.CreateRoleRequest) (*altalunev1.CreateRoleResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Trim name for consistency
	name := strings.TrimSpace(req.Name)

	// Check if role with same name already exists
	existingRole, err := s.roleRepo.GetByName(ctx, name)
	if err != nil && err != ErrRoleNotFound {
		s.log.Error("failed to check existing role",
			"error", err,
			"name", name,
		)
		return nil, altalune.NewUnexpectedError("failed to check existing role: %w", err)
	}

	if existingRole != nil {
		return nil, altalune.NewRoleAlreadyExistsError(name)
	}

	result, err := s.roleRepo.Create(ctx, &CreateRoleInput{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
	})
	if err != nil {
		if err == ErrRoleAlreadyExists {
			return nil, altalune.NewRoleAlreadyExistsError(name)
		}
		s.log.Error("failed to create role",
			"error", err,
			"name", name,
		)
		return nil, altalune.NewUnexpectedError("failed to create role: %w", err)
	}

	return &altalunev1.CreateRoleResponse{
		Role: &altalunev1.Role{
			Id:          result.PublicID,
			Name:        result.Name,
			Description: result.Description,
			CreatedAt:   timestamppb.New(result.CreatedAt),
			UpdatedAt:   timestamppb.New(result.UpdatedAt),
		},
		Message: "Role created successfully",
	}, nil
}

func (s *Service) GetRole(ctx context.Context, req *altalunev1.GetRoleRequest) (*altalunev1.GetRoleResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	role, err := s.roleRepo.GetByID(ctx, req.Id)
	if err != nil {
		if err == ErrRoleNotFound {
			return nil, altalune.NewRoleNotFoundError(req.Id)
		}
		s.log.Error("failed to get role", "error", err, "role_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to get role", err)
	}

	return &altalunev1.GetRoleResponse{
		Role: role.ToRoleProto(),
	}, nil
}

func (s *Service) UpdateRole(ctx context.Context, req *altalunev1.UpdateRoleRequest) (*altalunev1.UpdateRoleResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Get internal ID
	internalID, err := s.roleRepo.GetIDByPublicID(ctx, req.Id)
	if err != nil {
		if err == ErrRoleNotFound {
			return nil, altalune.NewRoleNotFoundError(req.Id)
		}
		return nil, altalune.NewUnexpectedError("failed to resolve role ID", err)
	}

	// Trim name for consistency
	name := strings.TrimSpace(req.Name)

	input := &UpdateRoleInput{
		ID:          internalID,
		PublicID:    req.Id,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
	}

	result, err := s.roleRepo.Update(ctx, input)
	if err != nil {
		if err == ErrRoleNotFound {
			return nil, altalune.NewRoleNotFoundError(req.Id)
		}
		if err == ErrRoleAlreadyExists {
			return nil, altalune.NewRoleAlreadyExistsError(name)
		}
		s.log.Error("failed to update role", "error", err, "role_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to update role", err)
	}

	s.log.Info("role updated successfully", "role_id", req.Id, "name", name)

	return &altalunev1.UpdateRoleResponse{
		Role:    result.ToRole().ToRoleProto(),
		Message: "Role updated successfully",
	}, nil
}

func (s *Service) DeleteRole(ctx context.Context, req *altalunev1.DeleteRoleRequest) (*altalunev1.DeleteRoleResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// First get the role to check if it's protected
	role, err := s.roleRepo.GetByID(ctx, req.Id)
	if err != nil {
		if err == ErrRoleNotFound {
			return nil, altalune.NewRoleNotFoundError(req.Id)
		}
		s.log.Error("failed to get role for deletion check", "error", err, "role_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to get role", err)
	}

	// Check if role is protected
	if role.Name == "superadmin" {
		s.log.Warn("attempt to delete protected role", "role_name", role.Name, "role_id", req.Id)
		return nil, altalune.NewRoleProtectedError(role.Name)
	}

	err = s.roleRepo.Delete(ctx, req.Id)
	if err != nil {
		if err == ErrRoleNotFound {
			return nil, altalune.NewRoleNotFoundError(req.Id)
		}
		if err == ErrRoleInUse {
			return nil, altalune.NewRoleInUseError(req.Id)
		}
		if err == ErrRoleProtected {
			return nil, altalune.NewRoleProtectedError(role.Name)
		}
		s.log.Error("failed to delete role", "error", err, "role_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to delete role", err)
	}

	s.log.Info("role deleted successfully", "role_id", req.Id)

	return &altalunev1.DeleteRoleResponse{
		Message: "Role deleted successfully",
	}, nil
}
