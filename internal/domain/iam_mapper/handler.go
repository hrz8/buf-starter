package iam_mapper

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	svc  altalunev1.IAMMapperServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.IAMMapperServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

// ==================== User-Role Mappings ====================

func (h *Handler) AssignUserRoles(
	ctx context.Context,
	req *connect.Request[altalunev1.AssignUserRolesRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires iam:write permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.AssignUserRoles(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) RemoveUserRoles(
	ctx context.Context,
	req *connect.Request[altalunev1.RemoveUserRolesRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires iam:write permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.RemoveUserRoles(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetUserRoles(
	ctx context.Context,
	req *connect.Request[altalunev1.GetUserRolesRequest],
) (*connect.Response[altalunev1.GetUserRolesResponse], error) {
	// Authorization: requires iam:read permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.GetUserRoles(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// ==================== Role-Permission Mappings ====================

func (h *Handler) AssignRolePermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.AssignRolePermissionsRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires iam:write permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.AssignRolePermissions(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) RemoveRolePermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.RemoveRolePermissionsRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires iam:write permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.RemoveRolePermissions(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetRolePermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.GetRolePermissionsRequest],
) (*connect.Response[altalunev1.GetRolePermissionsResponse], error) {
	// Authorization: requires iam:read permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.GetRolePermissions(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// ==================== User-Permission Mappings ====================

func (h *Handler) AssignUserPermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.AssignUserPermissionsRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires iam:write permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.AssignUserPermissions(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) RemoveUserPermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.RemoveUserPermissionsRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires iam:write permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.RemoveUserPermissions(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetUserPermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.GetUserPermissionsRequest],
) (*connect.Response[altalunev1.GetUserPermissionsResponse], error) {
	// Authorization: requires iam:read permission (global)
	if err := h.auth.CheckPermission(ctx, "iam:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.GetUserPermissions(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// ==================== Project Members ====================

func (h *Handler) AssignProjectMembers(
	ctx context.Context,
	req *connect.Request[altalunev1.AssignProjectMembersRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires member:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "member:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.AssignProjectMembers(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) RemoveProjectMembers(
	ctx context.Context,
	req *connect.Request[altalunev1.RemoveProjectMembersRequest],
) (*connect.Response[emptypb.Empty], error) {
	// Authorization: requires member:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "member:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.RemoveProjectMembers(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetProjectMembers(
	ctx context.Context,
	req *connect.Request[altalunev1.GetProjectMembersRequest],
) (*connect.Response[altalunev1.GetProjectMembersResponse], error) {
	// Authorization: requires member:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "member:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.GetProjectMembers(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetUserProjects(
	ctx context.Context,
	req *connect.Request[altalunev1.GetUserProjectsRequest],
) (*connect.Response[altalunev1.GetUserProjectsResponse], error) {
	// Authorization: requires authentication (user can see their own projects)
	if err := h.auth.CheckAuthenticated(ctx); err != nil {
		return nil, err
	}

	response, err := h.svc.GetUserProjects(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
