package server

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
)

func (s *Server) setupGRPCServices() {
	greeterv1.RegisterGreeterServiceServer(s.grpcServer, s.c.GetGreeterService())
	altalunev1.RegisterEmployeeServiceServer(s.grpcServer, s.c.GetEmployeeService())

}
