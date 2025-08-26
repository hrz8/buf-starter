package server

import (
	"net/http"

	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/container"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type Server struct {
	c   *container.Container
	cfg altalune.Config
	log altalune.Logger

	httpHandler http.Handler
	grpcServer  *grpc.Server
}

func NewServer(c *container.Container) *Server {
	return &Server{
		c:   c,
		cfg: c.GetConfig(),
		log: c.GetLogger(),
	}
}

func (s *Server) Bootstrap() (http.Handler, http.Handler) {
	mux := s.setupRoutes()
	handler := s.setupMiddleware(mux)

	s.httpHandler = h2c.NewHandler(handler, &http2.Server{})
	s.grpcServer = s.setupGRPCServices()

	return s.httpHandler, s.grpcServer
}
