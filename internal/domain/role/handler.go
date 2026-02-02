package role

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.RoleServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.RoleServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryRoles(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryRolesRequest],
) (*connect.Response[altalunev1.QueryRolesResponse], error) {
	// Authorization: requires role:read permission (global)
	if err := h.auth.CheckPermission(ctx, "role:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.QueryRoles(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) CreateRole(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateRoleRequest],
) (*connect.Response[altalunev1.CreateRoleResponse], error) {
	// Authorization: requires role:write permission (global)
	if err := h.auth.CheckPermission(ctx, "role:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.CreateRole(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetRole(
	ctx context.Context,
	req *connect.Request[altalunev1.GetRoleRequest],
) (*connect.Response[altalunev1.GetRoleResponse], error) {
	// Authorization: requires role:read permission (global)
	if err := h.auth.CheckPermission(ctx, "role:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.GetRole(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateRole(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateRoleRequest],
) (*connect.Response[altalunev1.UpdateRoleResponse], error) {
	// Authorization: requires role:write permission (global)
	if err := h.auth.CheckPermission(ctx, "role:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.UpdateRole(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteRole(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteRoleRequest],
) (*connect.Response[altalunev1.DeleteRoleResponse], error) {
	// Authorization: requires role:delete permission (global)
	if err := h.auth.CheckPermission(ctx, "role:delete"); err != nil {
		return nil, err
	}

	response, err := h.svc.DeleteRole(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
