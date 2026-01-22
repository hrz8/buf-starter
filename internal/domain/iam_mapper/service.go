package iam_mapper

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/domain/permission"
	"github.com/hrz8/altalune/internal/domain/project"
	"github.com/hrz8/altalune/internal/domain/role"
	"github.com/hrz8/altalune/internal/domain/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Valid project roles (per PROJECT_MEMBERSHIP_GUIDE.md)
const (
	ProjectRoleOwner  = "owner"
	ProjectRoleAdmin  = "admin"
	ProjectRoleMember = "member"
	ProjectRoleUser   = "user"
)

var validProjectRoles = map[string]bool{
	ProjectRoleOwner:  true,
	ProjectRoleAdmin:  true,
	ProjectRoleMember: true,
	ProjectRoleUser:   true,
}

type Service struct {
	altalunev1.UnimplementedIAMMapperServiceServer
	validator      protovalidate.Validator
	log            altalune.Logger
	mapperRepo     Repository
	userRepo       user.Repository
	roleRepo       role.Repository
	permissionRepo permission.Repository
	projectRepo    project.Repositor
}

func NewService(
	v protovalidate.Validator,
	log altalune.Logger,
	mapperRepo Repository,
	userRepo user.Repository,
	roleRepo role.Repository,
	permissionRepo permission.Repository,
	projectRepo project.Repositor,
) *Service {
	return &Service{
		validator:      v,
		log:            log,
		mapperRepo:     mapperRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		projectRepo:    projectRepo,
	}
}

// ==================== User-Role Mappings ====================

func (s *Service) AssignUserRoles(ctx context.Context, req *altalunev1.AssignUserRolesRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty role IDs
	if len(req.RoleIds) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve user public ID to internal ID
	userID, err := s.userRepo.GetIDByPublicID(ctx, req.UserId)
	if err != nil {
		s.log.Error("user not found for role assignment",
			"error", err,
			"user_public_id", req.UserId,
		)
		return nil, altalune.NewUserNotFoundError(req.UserId)
	}

	// Resolve role public IDs to internal IDs
	roleIDs := make([]int64, len(req.RoleIds))
	for i, publicID := range req.RoleIds {
		roleID, err := s.roleRepo.GetIDByPublicID(ctx, publicID)
		if err != nil {
			s.log.Error("role not found for assignment",
				"error", err,
				"role_public_id", publicID,
			)
			return nil, altalune.NewRoleNotFoundError(publicID)
		}
		roleIDs[i] = roleID
	}

	// Assign roles
	if err := s.mapperRepo.AssignUserRoles(ctx, userID, roleIDs); err != nil {
		s.log.Error("failed to assign user roles",
			"error", err,
			"user_id", userID,
			"role_ids", roleIDs,
		)
		return nil, altalune.NewUnexpectedError("failed to assign user roles: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) RemoveUserRoles(ctx context.Context, req *altalunev1.RemoveUserRolesRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty role IDs
	if len(req.RoleIds) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve user public ID to internal ID
	userID, err := s.userRepo.GetIDByPublicID(ctx, req.UserId)
	if err != nil {
		s.log.Error("user not found for role removal",
			"error", err,
			"user_public_id", req.UserId,
		)
		return nil, altalune.NewUserNotFoundError(req.UserId)
	}

	// Resolve role public IDs to internal IDs
	roleIDs := make([]int64, len(req.RoleIds))
	for i, publicID := range req.RoleIds {
		roleID, err := s.roleRepo.GetIDByPublicID(ctx, publicID)
		if err != nil {
			s.log.Error("role not found for removal",
				"error", err,
				"role_public_id", publicID,
			)
			return nil, altalune.NewRoleNotFoundError(publicID)
		}
		roleIDs[i] = roleID
	}

	// Remove roles
	if err := s.mapperRepo.RemoveUserRoles(ctx, userID, roleIDs); err != nil {
		s.log.Error("failed to remove user roles",
			"error", err,
			"user_id", userID,
			"role_ids", roleIDs,
		)
		return nil, altalune.NewUnexpectedError("failed to remove user roles: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) GetUserRoles(ctx context.Context, req *altalunev1.GetUserRolesRequest) (*altalunev1.GetUserRolesResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Resolve user public ID to internal ID
	userID, err := s.userRepo.GetIDByPublicID(ctx, req.UserId)
	if err != nil {
		s.log.Error("user not found for get roles",
			"error", err,
			"user_public_id", req.UserId,
		)
		return nil, altalune.NewUserNotFoundError(req.UserId)
	}

	// Get roles
	roles, err := s.mapperRepo.GetUserRoles(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user roles",
			"error", err,
			"user_id", userID,
		)
		return nil, altalune.NewUnexpectedError("failed to get user roles: %w", err)
	}

	return &altalunev1.GetUserRolesResponse{
		Roles: RolesToProto(roles),
	}, nil
}

// ==================== Role-Permission Mappings ====================

func (s *Service) AssignRolePermissions(ctx context.Context, req *altalunev1.AssignRolePermissionsRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty permission IDs
	if len(req.PermissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve role public ID to internal ID
	roleID, err := s.roleRepo.GetIDByPublicID(ctx, req.RoleId)
	if err != nil {
		s.log.Error("role not found for permission assignment",
			"error", err,
			"role_public_id", req.RoleId,
		)
		return nil, altalune.NewRoleNotFoundError(req.RoleId)
	}

	// Resolve permission public IDs to internal IDs
	permissionIDs := make([]int64, len(req.PermissionIds))
	for i, publicID := range req.PermissionIds {
		permissionID, err := s.permissionRepo.GetIDByPublicID(ctx, publicID)
		if err != nil {
			s.log.Error("permission not found for assignment",
				"error", err,
				"permission_public_id", publicID,
			)
			return nil, altalune.NewPermissionNotFoundError(publicID)
		}
		permissionIDs[i] = permissionID
	}

	// Assign permissions
	if err := s.mapperRepo.AssignRolePermissions(ctx, roleID, permissionIDs); err != nil {
		s.log.Error("failed to assign role permissions",
			"error", err,
			"role_id", roleID,
			"permission_ids", permissionIDs,
		)
		return nil, altalune.NewUnexpectedError("failed to assign role permissions: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) RemoveRolePermissions(ctx context.Context, req *altalunev1.RemoveRolePermissionsRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty permission IDs
	if len(req.PermissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve role public ID to internal ID
	roleID, err := s.roleRepo.GetIDByPublicID(ctx, req.RoleId)
	if err != nil {
		s.log.Error("role not found for permission removal",
			"error", err,
			"role_public_id", req.RoleId,
		)
		return nil, altalune.NewRoleNotFoundError(req.RoleId)
	}

	// Resolve permission public IDs to internal IDs
	permissionIDs := make([]int64, len(req.PermissionIds))
	for i, publicID := range req.PermissionIds {
		permissionID, err := s.permissionRepo.GetIDByPublicID(ctx, publicID)
		if err != nil {
			s.log.Error("permission not found for removal",
				"error", err,
				"permission_public_id", publicID,
			)
			return nil, altalune.NewPermissionNotFoundError(publicID)
		}
		permissionIDs[i] = permissionID
	}

	// Remove permissions
	if err := s.mapperRepo.RemoveRolePermissions(ctx, roleID, permissionIDs); err != nil {
		s.log.Error("failed to remove role permissions",
			"error", err,
			"role_id", roleID,
			"permission_ids", permissionIDs,
		)
		return nil, altalune.NewUnexpectedError("failed to remove role permissions: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) GetRolePermissions(ctx context.Context, req *altalunev1.GetRolePermissionsRequest) (*altalunev1.GetRolePermissionsResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Resolve role public ID to internal ID
	roleID, err := s.roleRepo.GetIDByPublicID(ctx, req.RoleId)
	if err != nil {
		s.log.Error("role not found for get permissions",
			"error", err,
			"role_public_id", req.RoleId,
		)
		return nil, altalune.NewRoleNotFoundError(req.RoleId)
	}

	// Get permissions
	permissions, err := s.mapperRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		s.log.Error("failed to get role permissions",
			"error", err,
			"role_id", roleID,
		)
		return nil, altalune.NewUnexpectedError("failed to get role permissions: %w", err)
	}

	return &altalunev1.GetRolePermissionsResponse{
		Permissions: PermissionsToProto(permissions),
	}, nil
}

// ==================== User-Permission Mappings ====================

func (s *Service) AssignUserPermissions(ctx context.Context, req *altalunev1.AssignUserPermissionsRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty permission IDs
	if len(req.PermissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve user public ID to internal ID
	userID, err := s.userRepo.GetIDByPublicID(ctx, req.UserId)
	if err != nil {
		s.log.Error("user not found for permission assignment",
			"error", err,
			"user_public_id", req.UserId,
		)
		return nil, altalune.NewUserNotFoundError(req.UserId)
	}

	// Resolve permission public IDs to internal IDs
	permissionIDs := make([]int64, len(req.PermissionIds))
	for i, publicID := range req.PermissionIds {
		permissionID, err := s.permissionRepo.GetIDByPublicID(ctx, publicID)
		if err != nil {
			s.log.Error("permission not found for assignment",
				"error", err,
				"permission_public_id", publicID,
			)
			return nil, altalune.NewPermissionNotFoundError(publicID)
		}
		permissionIDs[i] = permissionID
	}

	// Assign permissions
	if err := s.mapperRepo.AssignUserPermissions(ctx, userID, permissionIDs); err != nil {
		s.log.Error("failed to assign user permissions",
			"error", err,
			"user_id", userID,
			"permission_ids", permissionIDs,
		)
		return nil, altalune.NewUnexpectedError("failed to assign user permissions: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) RemoveUserPermissions(ctx context.Context, req *altalunev1.RemoveUserPermissionsRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty permission IDs
	if len(req.PermissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve user public ID to internal ID
	userID, err := s.userRepo.GetIDByPublicID(ctx, req.UserId)
	if err != nil {
		s.log.Error("user not found for permission removal",
			"error", err,
			"user_public_id", req.UserId,
		)
		return nil, altalune.NewUserNotFoundError(req.UserId)
	}

	// Resolve permission public IDs to internal IDs
	permissionIDs := make([]int64, len(req.PermissionIds))
	for i, publicID := range req.PermissionIds {
		permissionID, err := s.permissionRepo.GetIDByPublicID(ctx, publicID)
		if err != nil {
			s.log.Error("permission not found for removal",
				"error", err,
				"permission_public_id", publicID,
			)
			return nil, altalune.NewPermissionNotFoundError(publicID)
		}
		permissionIDs[i] = permissionID
	}

	// Remove permissions
	if err := s.mapperRepo.RemoveUserPermissions(ctx, userID, permissionIDs); err != nil {
		s.log.Error("failed to remove user permissions",
			"error", err,
			"user_id", userID,
			"permission_ids", permissionIDs,
		)
		return nil, altalune.NewUnexpectedError("failed to remove user permissions: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) GetUserPermissions(ctx context.Context, req *altalunev1.GetUserPermissionsRequest) (*altalunev1.GetUserPermissionsResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Resolve user public ID to internal ID
	userID, err := s.userRepo.GetIDByPublicID(ctx, req.UserId)
	if err != nil {
		s.log.Error("user not found for get permissions",
			"error", err,
			"user_public_id", req.UserId,
		)
		return nil, altalune.NewUserNotFoundError(req.UserId)
	}

	// Get permissions
	permissions, err := s.mapperRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user permissions",
			"error", err,
			"user_id", userID,
		)
		return nil, altalune.NewUnexpectedError("failed to get user permissions: %w", err)
	}

	return &altalunev1.GetUserPermissionsResponse{
		Permissions: PermissionsToProto(permissions),
	}, nil
}

// ==================== Project Members ====================

func (s *Service) AssignProjectMembers(ctx context.Context, req *altalunev1.AssignProjectMembersRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty members
	if len(req.Members) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve project public ID to internal ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		s.log.Error("project not found for member assignment",
			"error", err,
			"project_public_id", req.ProjectId,
		)
		return nil, altalune.NewProjectNotFound(req.ProjectId)
	}

	// Validate roles and resolve user IDs
	members := make([]ProjectMemberInput, len(req.Members))
	for i, member := range req.Members {
		// Validate project role
		if !validProjectRoles[member.Role] {
			return nil, altalune.NewInvalidPayloadError(
				fmt.Sprintf("invalid project role: %s (must be one of: owner, admin, member, user)", member.Role),
			)
		}

		// Resolve user public ID to internal ID
		userID, err := s.userRepo.GetIDByPublicID(ctx, member.UserId)
		if err != nil {
			s.log.Error("user not found for project member assignment",
				"error", err,
				"user_public_id", member.UserId,
			)
			return nil, altalune.NewUserNotFoundError(member.UserId)
		}

		members[i] = ProjectMemberInput{
			UserID: userID,
			Role:   member.Role,
		}
	}

	// Assign members
	if err := s.mapperRepo.AssignProjectMembers(ctx, projectID, members); err != nil {
		s.log.Error("failed to assign project members",
			"error", err,
			"project_id", projectID,
		)
		return nil, altalune.NewUnexpectedError("failed to assign project members: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) RemoveProjectMembers(ctx context.Context, req *altalunev1.RemoveProjectMembersRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Handle empty user IDs
	if len(req.UserIds) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Resolve project public ID to internal ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		s.log.Error("project not found for member removal",
			"error", err,
			"project_public_id", req.ProjectId,
		)
		return nil, altalune.NewProjectNotFound(req.ProjectId)
	}

	// Resolve user public IDs to internal IDs
	userIDs := make([]int64, len(req.UserIds))
	for i, publicID := range req.UserIds {
		userID, err := s.userRepo.GetIDByPublicID(ctx, publicID)
		if err != nil {
			s.log.Error("user not found for member removal",
				"error", err,
				"user_public_id", publicID,
			)
			return nil, altalune.NewUserNotFoundError(publicID)
		}
		userIDs[i] = userID
	}

	// Remove members
	if err := s.mapperRepo.RemoveProjectMembers(ctx, projectID, userIDs); err != nil {
		if err == ErrCannotRemoveLastOwner {
			return nil, altalune.NewInvalidPayloadError("cannot remove last owner from project")
		}
		s.log.Error("failed to remove project members",
			"error", err,
			"project_id", projectID,
			"user_ids", userIDs,
		)
		return nil, altalune.NewUnexpectedError("failed to remove project members: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) GetProjectMembers(ctx context.Context, req *altalunev1.GetProjectMembersRequest) (*altalunev1.GetProjectMembersResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Resolve project public ID to internal ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		s.log.Error("project not found for get members",
			"error", err,
			"project_public_id", req.ProjectId,
		)
		return nil, altalune.NewProjectNotFound(req.ProjectId)
	}

	// Get members
	members, err := s.mapperRepo.GetProjectMembers(ctx, projectID)
	if err != nil {
		s.log.Error("failed to get project members",
			"error", err,
			"project_id", projectID,
		)
		return nil, altalune.NewUnexpectedError("failed to get project members: %w", err)
	}

	return &altalunev1.GetProjectMembersResponse{
		Members: ProjectMembersToProto(members),
	}, nil
}

func (s *Service) GetUserProjects(ctx context.Context, req *altalunev1.GetUserProjectsRequest) (*altalunev1.GetUserProjectsResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	userID, err := s.userRepo.GetIDByPublicID(ctx, req.UserId)
	if err != nil {
		s.log.Error("user not found for get projects",
			"error", err,
			"user_public_id", req.UserId,
		)
		return nil, altalune.NewUserNotFoundError(req.UserId)
	}

	projects, err := s.mapperRepo.GetUserProjects(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user projects",
			"error", err,
			"user_id", userID,
		)
		return nil, altalune.NewUnexpectedError("failed to get user projects: %w", err)
	}

	return &altalunev1.GetUserProjectsResponse{
		Projects: UserProjectsToProto(projects),
	}, nil
}
