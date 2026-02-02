package project

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.ProjectServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.ProjectServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryProjects(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryProjectsRequest],
) (*connect.Response[altalunev1.QueryProjectsResponse], error) {
	// Authorization: requires project:read OR dashboard:read permission
	// - project:read: admin access to all projects
	// - dashboard:read: regular user access (frontend filters by membership)
	if err := h.auth.CheckAnyPermission(ctx, []string{"project:read", "dashboard:read"}); err != nil {
		return nil, err
	}

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
	// Authorization: requires project:write permission (global, not project-scoped)
	if err := h.auth.CheckPermission(ctx, "project:write"); err != nil {
		return nil, err
	}

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
	// Authorization: requires project:read permission (global, not project-scoped)
	if err := h.auth.CheckPermission(ctx, "project:read"); err != nil {
		return nil, err
	}

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
	// Authorization: requires project:write permission (global, not project-scoped)
	if err := h.auth.CheckPermission(ctx, "project:write"); err != nil {
		return nil, err
	}

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
	// Authorization: requires project:delete permission (global, not project-scoped)
	if err := h.auth.CheckPermission(ctx, "project:delete"); err != nil {
		return nil, err
	}

	response, err := h.svc.DeleteProject(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
