package grpcserver

import "time"

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
