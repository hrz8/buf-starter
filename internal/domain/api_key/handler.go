package api_key

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.ApiKeyServiceServer
}

func NewHandler(svc altalunev1.ApiKeyServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) QueryApiKeys(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryApiKeysRequest],
) (*connect.Response[altalunev1.QueryApiKeysResponse], error) {
	response, err := h.svc.QueryApiKeys(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) CreateApiKey(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateApiKeyRequest],
) (*connect.Response[altalunev1.CreateApiKeyResponse], error) {
	response, err := h.svc.CreateApiKey(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetApiKey(
	ctx context.Context,
	req *connect.Request[altalunev1.GetApiKeyRequest],
) (*connect.Response[altalunev1.GetApiKeyResponse], error) {
	response, err := h.svc.GetApiKey(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateApiKey(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateApiKeyRequest],
) (*connect.Response[altalunev1.UpdateApiKeyResponse], error) {
	response, err := h.svc.UpdateApiKey(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteApiKey(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteApiKeyRequest],
) (*connect.Response[altalunev1.DeleteApiKeyResponse], error) {
	response, err := h.svc.DeleteApiKey(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) ActivateApiKey(
	ctx context.Context,
	req *connect.Request[altalunev1.ActivateApiKeyRequest],
) (*connect.Response[altalunev1.ActivateApiKeyResponse], error) {
	response, err := h.svc.ActivateApiKey(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeactivateApiKey(
	ctx context.Context,
	req *connect.Request[altalunev1.DeactivateApiKeyRequest],
) (*connect.Response[altalunev1.DeactivateApiKeyResponse], error) {
	response, err := h.svc.DeactivateApiKey(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}