package oauth_provider

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.OAuthProviderServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.OAuthProviderServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryOAuthProviders(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryOAuthProvidersRequest],
) (*connect.Response[altalunev1.QueryOAuthProvidersResponse], error) {
	// Authorization: requires client:read permission (global - OAuth providers are global resources)
	if err := h.auth.CheckPermission(ctx, "client:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.QueryOAuthProviders(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) CreateOAuthProvider(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateOAuthProviderRequest],
) (*connect.Response[altalunev1.CreateOAuthProviderResponse], error) {
	// Authorization: requires client:write permission (global)
	if err := h.auth.CheckPermission(ctx, "client:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.CreateOAuthProvider(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetOAuthProvider(
	ctx context.Context,
	req *connect.Request[altalunev1.GetOAuthProviderRequest],
) (*connect.Response[altalunev1.GetOAuthProviderResponse], error) {
	// Authorization: requires client:read permission (global)
	if err := h.auth.CheckPermission(ctx, "client:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.GetOAuthProvider(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateOAuthProvider(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateOAuthProviderRequest],
) (*connect.Response[altalunev1.UpdateOAuthProviderResponse], error) {
	// Authorization: requires client:write permission (global)
	if err := h.auth.CheckPermission(ctx, "client:write"); err != nil {
		return nil, err
	}

	response, err := h.svc.UpdateOAuthProvider(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteOAuthProvider(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteOAuthProviderRequest],
) (*connect.Response[altalunev1.DeleteOAuthProviderResponse], error) {
	// Authorization: requires client:delete permission (global)
	if err := h.auth.CheckPermission(ctx, "client:delete"); err != nil {
		return nil, err
	}

	response, err := h.svc.DeleteOAuthProvider(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) RevealClientSecret(
	ctx context.Context,
	req *connect.Request[altalunev1.RevealClientSecretRequest],
) (*connect.Response[altalunev1.RevealClientSecretResponse], error) {
	// Authorization: requires client:read permission (global)
	if err := h.auth.CheckPermission(ctx, "client:read"); err != nil {
		return nil, err
	}

	response, err := h.svc.RevealClientSecret(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
