package project

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.ProjectServiceServer
}

func NewHandler(svc altalunev1.ProjectServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) QueryProjects(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryProjectsRequest],
) (*connect.Response[altalunev1.QueryProjectsResponse], error) {
	response, err := h.svc.QueryProjects(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
