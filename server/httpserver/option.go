package httpserver

import (
	"net/http"
	"time"
)

type Option func(*Server)

type options struct {
	port           int
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
	cleanupTimeout time.Duration
}

func defaultOptions() *options {
	return &options{
		port:           3100,
		readTimeout:    15 * time.Second,
		writeTimeout:   15 * time.Second,
		idleTimeout:    60 * time.Second,
		cleanupTimeout: 10 * time.Second,
	}
}

func WithHandler(h http.Handler) Option {
	return func(s *Server) {
		if h != nil {
			s.httpHandler = h
		}
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.opts.port = port
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.opts.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.opts.writeTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.opts.idleTimeout = timeout
	}
}

func WithCleanupTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.opts.cleanupTimeout = timeout
	}
}
