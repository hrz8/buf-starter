package api_key

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.ApiKeyServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.ApiKeyServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryApiKeys(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryApiKeysRequest],
) (*connect.Response[altalunev1.QueryApiKeysResponse], error) {
	// Authorization: requires apikey:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "apikey:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires apikey:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "apikey:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires apikey:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "apikey:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires apikey:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "apikey:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires apikey:delete permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "apikey:delete", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires apikey:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "apikey:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires apikey:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "apikey:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.DeactivateApiKey(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
