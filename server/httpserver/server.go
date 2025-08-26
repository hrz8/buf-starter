package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	httpServer  *http.Server
	httpHandler http.Handler

	cancelCtx context.CancelFunc
	notify    chan error

	opts *options
}

func NewHTTPServer(opts ...Option) *Server {
	s := &Server{
		notify: make(chan error, 1),
		opts:   defaultOptions(),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.httpHandler == nil {
		s.httpHandler = http.NewServeMux()
	}

	return s
}

func (s *Server) Start() {
	baseCtx, cancel := context.WithCancel(context.Background())

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.opts.port),
		Handler:      s.httpHandler,
		ReadTimeout:  s.opts.readTimeout,
		WriteTimeout: s.opts.writeTimeout,
		IdleTimeout:  s.opts.idleTimeout,
	}
	s.httpServer.BaseContext = func(net.Listener) context.Context {
		return baseCtx
	}
	s.cancelCtx = cancel

	go func() {
		defer close(s.notify)
		s.notify <- s.httpServer.ListenAndServe()
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.opts.cleanupTimeout)
	defer cancel()

	defer s.cancelCtx()
	if err := s.httpServer.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("http server shutdown failed: %w", err)
	}

	return nil
}
