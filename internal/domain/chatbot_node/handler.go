package chatbot_node

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.ChatbotNodeServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.ChatbotNodeServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) ListNodes(
	ctx context.Context,
	req *connect.Request[altalunev1.ListNodesRequest],
) (*connect.Response[altalunev1.ListNodesResponse], error) {
	// Authorization: requires chatbot:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "chatbot:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires chatbot:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "chatbot:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires chatbot:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "chatbot:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires chatbot:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "chatbot:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

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
	// Authorization: requires chatbot:delete permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "chatbot:delete", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.DeleteNode(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
