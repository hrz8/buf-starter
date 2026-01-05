package permission

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.PermissionServiceServer
}

func NewHandler(svc altalunev1.PermissionServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) QueryPermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryPermissionsRequest],
) (*connect.Response[altalunev1.QueryPermissionsResponse], error) {
	response, err := h.svc.QueryPermissions(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) CreatePermission(
	ctx context.Context,
	req *connect.Request[altalunev1.CreatePermissionRequest],
) (*connect.Response[altalunev1.CreatePermissionResponse], error) {
	response, err := h.svc.CreatePermission(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetPermission(
	ctx context.Context,
	req *connect.Request[altalunev1.GetPermissionRequest],
) (*connect.Response[altalunev1.GetPermissionResponse], error) {
	response, err := h.svc.GetPermission(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdatePermission(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdatePermissionRequest],
) (*connect.Response[altalunev1.UpdatePermissionResponse], error) {
	response, err := h.svc.UpdatePermission(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeletePermission(
	ctx context.Context,
	req *connect.Request[altalunev1.DeletePermissionRequest],
) (*connect.Response[altalunev1.DeletePermissionResponse], error) {
	response, err := h.svc.DeletePermission(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
