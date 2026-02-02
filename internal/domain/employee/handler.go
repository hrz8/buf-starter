package employee

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
	svc  altalunev1.EmployeeServiceServer
	auth *auth.Authorizer
}

func NewHandler(svc altalunev1.EmployeeServiceServer, authorizer *auth.Authorizer) *Handler {
	return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryEmployees(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryEmployeesRequest],
) (*connect.Response[altalunev1.QueryEmployeesResponse], error) {
	// Authorization: requires employee:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.QueryEmployees(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) CreateEmployee(
	ctx context.Context,
	req *connect.Request[altalunev1.CreateEmployeeRequest],
) (*connect.Response[altalunev1.CreateEmployeeResponse], error) {
	// Authorization: requires employee:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "employee:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.CreateEmployee(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetEmployee(
	ctx context.Context,
	req *connect.Request[altalunev1.GetEmployeeRequest],
) (*connect.Response[altalunev1.GetEmployeeResponse], error) {
	// Authorization: requires employee:read permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.GetEmployee(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) UpdateEmployee(
	ctx context.Context,
	req *connect.Request[altalunev1.UpdateEmployeeRequest],
) (*connect.Response[altalunev1.UpdateEmployeeResponse], error) {
	// Authorization: requires employee:write permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "employee:write", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.UpdateEmployee(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteEmployee(
	ctx context.Context,
	req *connect.Request[altalunev1.DeleteEmployeeRequest],
) (*connect.Response[altalunev1.DeleteEmployeeResponse], error) {
	// Authorization: requires employee:delete permission and project membership
	if err := h.auth.CheckProjectAccess(ctx, "employee:delete", req.Msg.ProjectId); err != nil {
		return nil, err
	}

	response, err := h.svc.DeleteEmployee(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
