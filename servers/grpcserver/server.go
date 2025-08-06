package grpcserver

import (
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener

	notify chan error

	opts *options
}

func NewGRPCServer(provider Provider, opts ...Option) *Server {
	var grpcServer *grpc.Server
	if provider != nil {
		grpcServer = provider.GRPCServer()
	}

	if grpcServer == nil {
		grpcServer = grpc.NewServer()
	}

	s := &Server{
		grpcServer: grpcServer,
		notify:     make(chan error, 1),
		opts:       defaultOptions(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.port))
	if err != nil {
		s.notify <- fmt.Errorf("failed to create listener: %w", err)
		close(s.notify)
		return
	}
	s.listener = listener

	go func() {
		defer close(s.notify)
		s.notify <- s.grpcServer.Serve(listener)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Stop() error {
	// Create a done channel to signal when graceful stop is complete
	done := make(chan struct{})

	go func() {
		defer close(done)
		s.grpcServer.GracefulStop()
	}()

	// Wait for graceful stop or timeout
	select {
	case <-done:
		// Graceful stop completed
	case <-time.After(s.opts.cleanupTimeout):
		// Timeout exceeded, force stop
		s.grpcServer.Stop()
		return errors.New("gRPC server graceful stop timeout exceeded")
	}

	return nil
}
