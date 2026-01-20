package chatbot_node

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.ChatbotNodeServiceServer
}

func NewHandler(svc altalunev1.ChatbotNodeServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListNodes(
	ctx context.Context,
	req *connect.Request[altalunev1.ListNodesRequest],
) (*connect.Response[altalunev1.ListNodesResponse], error) {
	response, err := h.svc.ListNodes(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) CreateNode(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateNodeRequest],
) (*connect.Response[altalunev1.CreateNodeResponse], error) {
	response, err := h.svc.CreateNode(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetNode(
	ctx context.Context,
	req *connect.Request[altalunev1.GetNodeRequest],
) (*connect.Response[altalunev1.GetNodeResponse], error) {
	response, err := h.svc.GetNode(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateNode(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateNodeRequest],
) (*connect.Response[altalunev1.UpdateNodeResponse], error) {
	response, err := h.svc.UpdateNode(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteNode(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteNodeRequest],
) (*connect.Response[altalunev1.DeleteNodeResponse], error) {
	response, err := h.svc.DeleteNode(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
