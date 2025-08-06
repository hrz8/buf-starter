package server

import (
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
)

func (s *Server) setupGRPCServices() {
	greeterv1.RegisterGreeterServiceServer(s.grpcServer, s.c.GetGreeterService())
}
