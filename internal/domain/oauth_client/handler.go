package oauth_client

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.OAuthClientServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.OAuthClientServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

// CreateOAuthClient handles OAuth client creation requests
func (h *Handler) CreateOAuthClient(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateOAuthClientRequest],
) (*connect.Response[altalunev1.CreateOAuthClientResponse], error) {
	// Authorization: requires client:write permission (global - no project_id)
	if err := h.auth.CheckPermission(ctx, "client:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.CreateOAuthClient(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// QueryOAuthClients handles OAuth client query requests
func (h *Handler) QueryOAuthClients(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryOAuthClientsRequest],
) (*connect.Response[altalunev1.QueryOAuthClientsResponse], error) {
	// Authorization: requires client:read permission (global - no project_id)
	if err := h.auth.CheckPermission(ctx, "client:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.QueryOAuthClients(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// GetOAuthClient handles OAuth client retrieval by ID
func (h *Handler) GetOAuthClient(
	ctx context.Context,
	req *connect.Request[altalunev1.GetOAuthClientRequest],
) (*connect.Response[altalunev1.GetOAuthClientResponse], error) {
	// Authorization: requires client:read permission (global - no project_id)
	if err := h.auth.CheckPermission(ctx, "client:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.GetOAuthClient(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// UpdateOAuthClient handles OAuth client update requests
func (h *Handler) UpdateOAuthClient(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateOAuthClientRequest],
) (*connect.Response[altalunev1.UpdateOAuthClientResponse], error) {
	// Authorization: requires client:write permission (global - no project_id)
	if err := h.auth.CheckPermission(ctx, "client:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.UpdateOAuthClient(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// DeleteOAuthClient handles OAuth client deletion requests
func (h *Handler) DeleteOAuthClient(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteOAuthClientRequest],
) (*connect.Response[altalunev1.DeleteOAuthClientResponse], error) {
	// Authorization: requires client:delete permission (global - no project_id)
	if err := h.auth.CheckPermission(ctx, "client:delete"); err != nil {
		return nil, err
	}

	response, err := h.svc.DeleteOAuthClient(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

// RevealOAuthClientSecret handles OAuth client secret reveal requests
func (h *Handler) RevealOAuthClientSecret(
	ctx context.Context,
	req *connect.Request[altalunev1.RevealOAuthClientSecretRequest],
) (*connect.Response[altalunev1.RevealOAuthClientSecretResponse], error) {
	// Authorization: requires client:read permission (global - no project_id)
	if err := h.auth.CheckPermission(ctx, "client:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.RevealOAuthClientSecret(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
