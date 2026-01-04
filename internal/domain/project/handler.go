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

func (h *Handler) CreateProject(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateProjectRequest],
) (*connect.Response[altalunev1.CreateProjectResponse], error) {
	response, err := h.svc.CreateProject(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetProject(
	ctx context.Context,
	req *connect.Request[altalunev1.GetProjectRequest],
) (*connect.Response[altalunev1.GetProjectResponse], error) {
	response, err := h.svc.GetProject(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateProject(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateProjectRequest],
) (*connect.Response[altalunev1.UpdateProjectResponse], error) {
	response, err := h.svc.UpdateProject(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteProject(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteProjectRequest],
) (*connect.Response[altalunev1.DeleteProjectResponse], error) {
	response, err := h.svc.DeleteProject(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
