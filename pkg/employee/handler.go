package employee

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

type Handler struct {
	svc altalunev1.EmployeeServiceServer
}

func NewHandler(svc altalunev1.EmployeeServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) QueryEmployees(
	ctx context.Context,
	req *connect.Request[altalunev1.QueryEmployeesRequest],
) (*connect.Response[altalunev1.QueryEmployeesResponse], error) {
	response, err := h.svc.QueryEmployees(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
