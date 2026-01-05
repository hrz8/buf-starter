package iam_mapper

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	svc altalunev1.IAMMapperServiceServer
}

func NewHandler(svc altalunev1.IAMMapperServiceServer) *Handler {
	return &Handler{svc: svc}
}

// ==================== User-Role Mappings ====================

func (h *Handler) AssignUserRoles(
	ctx context.Context,
	req *connect.Request[altalunev1.AssignUserRolesRequest],
) (*connect.Response[emptypb.Empty], error) {
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
	response, err := h.svc.GetProjectMembers(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
