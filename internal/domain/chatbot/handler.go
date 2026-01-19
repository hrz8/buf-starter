package chatbot

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.ChatbotServiceServer
}

func NewHandler(svc altalunev1.ChatbotServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetChatbotConfig(
	ctx context.Context,
	req *connect.Request[altalunev1.GetChatbotConfigRequest],
) (*connect.Response[altalunev1.GetChatbotConfigResponse], error) {
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
	response, err := h.svc.UpdateModuleConfig(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
