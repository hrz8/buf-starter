package permission

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.PermissionServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.PermissionServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryPermissions(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryPermissionsRequest],
) (*connect.Response[altalunev1.QueryPermissionsResponse], error) {
	// Authorization: requires permission:read permission (global)
	if err := h.auth.CheckPermission(ctx, "permission:read"); err != nil {
		return nil, err
	}

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
	// Authorization: requires permission:write permission (global)
	if err := h.auth.CheckPermission(ctx, "permission:write"); err != nil {
		return nil, err
	}

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
	// Authorization: requires permission:read permission (global)
	if err := h.auth.CheckPermission(ctx, "permission:read"); err != nil {
		return nil, err
	}

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
	// Authorization: requires permission:write permission (global)
	if err := h.auth.CheckPermission(ctx, "permission:write"); err != nil {
		return nil, err
	}

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
	// Authorization: requires permission:delete permission (global)
	if err := h.auth.CheckPermission(ctx, "permission:delete"); err != nil {
		return nil, err
	}

	response, err := h.svc.DeletePermission(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
