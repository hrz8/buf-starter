package role

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.RoleServiceServer
}

func NewHandler(svc altalunev1.RoleServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) QueryRoles(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryRolesRequest],
) (*connect.Response[altalunev1.QueryRolesResponse], error) {
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
	response, err := h.svc.DeleteRole(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
