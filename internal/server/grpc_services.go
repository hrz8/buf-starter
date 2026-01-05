package server

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Server) setupGRPCServices() *grpc.Server {
	grpcServer := grpc.NewServer()

	// Examples
	greeterv1.RegisterGreeterServiceServer(grpcServer, s.c.GetGreeterService())
	altalunev1.RegisterEmployeeServiceServer(grpcServer, s.c.GetEmployeeService())

	// Domains
	altalunev1.RegisterProjectServiceServer(grpcServer, s.c.GetProjectService())
	altalunev1.RegisterApiKeyServiceServer(grpcServer, s.c.GetApiKeyService())

	// IAM Domains
	altalunev1.RegisterUserServiceServer(grpcServer, s.c.GetUserService())
	altalunev1.RegisterRoleServiceServer(grpcServer, s.c.GetRoleService())
	altalunev1.RegisterPermissionServiceServer(grpcServer, s.c.GetPermissionService())
	altalunev1.RegisterIAMMapperServiceServer(grpcServer, s.c.GetIAMMapperService())

	reflection.Register(grpcServer)

	return grpcServer
}
