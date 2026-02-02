package chatbot

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.ChatbotServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.ChatbotServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) GetChatbotConfig(
	ctx context.Context,
	req *connect.Request[altalunev1.GetChatbotConfigRequest],
) (*connect.Response[altalunev1.GetChatbotConfigResponse], error) {
	// Authorization: requires chatbot:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "chatbot:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.GetChatbotConfig(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateModuleConfig(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateModuleConfigRequest],
) (*connect.Response[altalunev1.UpdateModuleConfigResponse], error) {
	// Authorization: requires chatbot:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "chatbot:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.UpdateModuleConfig(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
