package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/config"
	"github.com/hrz8/altalune/internal/container"
	"github.com/hrz8/altalune/internal/server"
	"github.com/hrz8/altalune/server/grpcserver"
	"github.com/hrz8/altalune/server/httpserver"
	"github.com/spf13/cobra"
)

func NewServeCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the altalune server",
		Long:  "Start the altalune server with HTTP REST and gRPC APIs",
		RunE:  serve(rootCmd),
	}

	return cmd
}

func serve(rootCmd *cobra.Command) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Get and load configuration
		configPath, _ := rootCmd.PersistentFlags().GetString("config")
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("error loading configuration file: %w", err)
		}

		// Bootstrapping
		c, err := container.CreateContainer(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create application container: %w", err)
		}
		if !c.IsHealthy(ctx) {
			return fmt.Errorf("container is not healthy, cannot run migration")
		}
		srv := server.NewServer(c)
		httpHandler, grpcServer := srv.Bootstrap()

		// Create HTTP server
		httpSrv := httpserver.NewHTTPServer(
			httpserver.WithHandler(httpHandler),
			httpserver.WithPort(cfg.GetServerPort()),
			httpserver.WithReadTimeout(cfg.GetServerReadTimeout()),
			httpserver.WithWriteTimeout(cfg.GetServerWriteTimeout()),
			httpserver.WithIdleTimeout(cfg.GetServerIdleTimeout()),
			httpserver.WithCleanupTimeout(cfg.GetServerCleanupTimeout()),
		)

		// Create gRPC server
		grpcSrv := grpcserver.NewGRPCServer(
			grpcserver.WithHandler(grpcServer),
			grpcserver.WithPort(cfg.GetServerPort()+1),
			grpcserver.WithCleanupTimeout(cfg.GetServerCleanupTimeout()),
		)

		// Start servers
		go func() {
			log.Printf("🚀 starting HTTP server at port: %d\n", cfg.GetServerPort())
			httpSrv.Start()
		}()

		go func() {
			log.Printf("🚀 starting gRPC server at port: %d\n", cfg.GetServerPort()+1)
			grpcSrv.Start()
		}()

		defer cleanup(cfg,
			func() error {
				if err := httpSrv.Stop(); err != nil {
					log.Printf("failed shutdown HTTP server: %v\n", err)
					return err
				}
				return nil
			},
			func() error {
				if err := grpcSrv.Stop(); err != nil {
					log.Printf("failed shutdown gRPC server: %v\n", err)
					return err
				}
				return nil
			},
			func() error {
				if err := c.Shutdown(); err != nil {
					log.Printf("failed to shutdown application container: %v\n", err)
					return err
				}
				return nil
			},
		)

		select {
		case <-ctx.Done():
			time.Sleep(100 * time.Millisecond)
			log.Println("🍀 performing graceful shutdown...")
		case err := <-httpSrv.Notify():
			if err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("HTTP server listen error: %w", err)
			}
		case err := <-grpcSrv.Notify():
			if err != nil {
				return fmt.Errorf("gRPC server listen error: %w", err)
			}
		}

		return nil
	}
}

func cleanup(cfg altalune.Config, cleanupFuncs ...func() error) {
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), cfg.GetServerCleanupTimeout())
	defer cleanupCancel()

	var wg sync.WaitGroup

	for _, fn := range cleanupFuncs {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				log.Printf("cleanup function error: %v\n", err)
			}
		}(fn)
	}

	cleanupDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(cleanupDone)
	}()

	select {
	case <-cleanupCtx.Done():
		log.Println("⚠️ cleanup done partially, because it takes longer than it should")
	case <-cleanupDone:
		log.Println("✨ cleanup done")
	}
}
