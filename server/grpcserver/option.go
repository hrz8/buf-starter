package grpcserver

import (
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type options struct {
	port           int
	cleanupTimeout time.Duration
}

type Option func(*Server)

func defaultOptions() *options {
	return &options{
		port:           3001,
		cleanupTimeout: 30 * time.Second,
	}
}

func WithHandler(gs http.Handler) Option {
	return func(s *Server) {
		if gs != nil {
			if g, ok := gs.(*grpc.Server); ok {
				s.grpcServer = g
			}
		}
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.opts.port = port
	}
}

func WithCleanupTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.opts.cleanupTimeout = timeout
	}
}
