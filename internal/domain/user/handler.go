package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.UserServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.UserServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryUsers(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryUsersRequest],
) (*connect.Response[altalunev1.QueryUsersResponse], error) {
	// Authorization: requires user:read permission (global)
	if err := h.auth.CheckPermission(ctx, "user:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.QueryUsers(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) CreateUser(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateUserRequest],
) (*connect.Response[altalunev1.CreateUserResponse], error) {
	// Authorization: requires user:write permission (global)
	if err := h.auth.CheckPermission(ctx, "user:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.CreateUser(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetUser(
	ctx context.Context,
	req *connect.Request[altalunev1.GetUserRequest],
) (*connect.Response[altalunev1.GetUserResponse], error) {
	// Authorization: requires user:read permission (global)
	if err := h.auth.CheckPermission(ctx, "user:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.GetUser(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateUser(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateUserRequest],
) (*connect.Response[altalunev1.UpdateUserResponse], error) {
	// Authorization: requires user:write permission (global)
	if err := h.auth.CheckPermission(ctx, "user:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.UpdateUser(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteUser(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteUserRequest],
) (*connect.Response[altalunev1.DeleteUserResponse], error) {
	// Authorization: requires user:delete permission (global)
	if err := h.auth.CheckPermission(ctx, "user:delete"); err != nil {
		return nil, err
	}

	response, err := h.svc.DeleteUser(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) ActivateUser(
	ctx context.Context,
	req *connect.Request[altalunev1.ActivateUserRequest],
) (*connect.Response[altalunev1.ActivateUserResponse], error) {
	// Authorization: requires user:write permission (global)
	if err := h.auth.CheckPermission(ctx, "user:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.ActivateUser(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeactivateUser(
	ctx context.Context,
	req *connect.Request[altalunev1.DeactivateUserRequest],
) (*connect.Response[altalunev1.DeactivateUserResponse], error) {
	// Authorization: requires user:write permission (global)
	if err := h.auth.CheckPermission(ctx, "user:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.DeactivateUser(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
