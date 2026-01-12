package authserver

import (
	"net/http"

	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/container"
	"github.com/hrz8/altalune/internal/server"
)

type Server struct {
	c   *container.Container
	cfg altalune.Config
	log altalune.Logger
}

func NewServer(c *container.Container) *Server {
	return &Server{
		c:   c,
		cfg: c.GetConfig(),
		log: c.GetLogger(),
	}
}

func (s *Server) Bootstrap() http.Handler {
	mux := s.setupRoutes()
	handler := s.setupMiddleware(mux)
	return handler
}

func (s *Server) setupMiddleware(handler http.Handler) http.Handler {
	handler = server.RecoveryMiddleware(handler, s.log)
	if s.cfg.IsHTTPLoggingEnabled() {
		handler = server.LoggingMiddleware(handler, s.log)
	}
	handler = server.SecurityMiddleware(handler)
	if s.cfg.IsCORSEnabled() {
		handler = server.CORSMiddleware(handler, s.cfg.GetAllowedOrigins())
	}
	return handler
}
