package greeter

import (
	"context"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
)

type Handler struct {
	svc greeterv1.GreeterServiceServer
}

func NewHandler(svc greeterv1.GreeterServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SayHello(
	ctx context.Context,
	req *connect.Request[greeterv1.SayHelloRequest],
) (*connect.Response[greeterv1.SayHelloResponse], error) {
	response, err := h.svc.SayHello(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

func (h *Handler) GetAllowedNames(
	ctx context.Context,
	req *connect.Request[greeterv1.GetAllowedNamesRequest],
) (*connect.Response[greeterv1.GetAllowedNamesResponse], error) {
	response, err := h.svc.GetAllowedNames(ctx, req.Msg)
	if err != nil {
		return nil, altalune.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}
