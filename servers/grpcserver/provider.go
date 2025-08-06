package grpcserver

import "google.golang.org/grpc"

type Provider interface {
	GRPCServer() *grpc.Server
}
