package iam_mapper

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/domain/permission"
	"github.com/hrz8/altalune/internal/domain/role"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RolesToProto converts a slice of role domain models to protobuf messages
func RolesToProto(roles []*role.Role) []*altalunev1.Role {
	if roles == nil {
		return []*altalunev1.Role{}
	}

	protoRoles := make([]*altalunev1.Role, len(roles))
	for i, r := range roles {
		protoRoles[i] = r.ToRoleProto()
	}
	return protoRoles
}

// PermissionsToProto converts a slice of permission domain models to protobuf messages
func PermissionsToProto(permissions []*permission.Permission) []*altalunev1.Permission {
	if permissions == nil {
		return []*altalunev1.Permission{}
	}

	protoPermissions := make([]*altalunev1.Permission, len(permissions))
	for i, p := range permissions {
		protoPermissions[i] = p.ToPermissionProto()
	}
	return protoPermissions
}

// ProjectMembersToProto converts a slice of project member domain models to protobuf messages
func ProjectMembersToProto(members []*ProjectMemberWithUser) []*altalunev1.ProjectMemberWithUser {
	if members == nil {
		return []*altalunev1.ProjectMemberWithUser{}
	}

	protoMembers := make([]*altalunev1.ProjectMemberWithUser, len(members))
	for i, m := range members {
		protoMembers[i] = &altalunev1.ProjectMemberWithUser{
			User:      m.User.ToUserProto(),
			Role:      m.Role,
			CreatedAt: timestamppb.New(m.CreatedAt),
		}
	}
	return protoMembers
}

// UserProjectsToProto converts a slice of user project membership domain models to protobuf messages
func UserProjectsToProto(projects []*UserProjectMembership) []*altalunev1.UserProjectMembership {
	if projects == nil {
		return []*altalunev1.UserProjectMembership{}
	}

	protoProjects := make([]*altalunev1.UserProjectMembership, len(projects))
	for i, p := range projects {
		protoProjects[i] = &altalunev1.UserProjectMembership{
			ProjectId:   p.ProjectID,
			ProjectName: p.ProjectName,
			Role:        p.Role,
			JoinedAt:    timestamppb.New(p.JoinedAt),
		}
	}
	return protoProjects
}
