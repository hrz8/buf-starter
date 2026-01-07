package oauth_provider

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.OAuthProviderServiceServer
}

func NewHandler(svc altalunev1.OAuthProviderServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) QueryOAuthProviders(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryOAuthProvidersRequest],
) (*connect.Response[altalunev1.QueryOAuthProvidersResponse], error) {
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
	response, err := h.svc.RevealClientSecret(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
