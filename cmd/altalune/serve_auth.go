package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/authserver"
	"github.com/hrz8/altalune/internal/config"
	"github.com/hrz8/altalune/internal/container"
	"github.com/hrz8/altalune/server/httpserver"
	"github.com/spf13/cobra"
)

func NewServeAuthCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "serve-auth",
		Short: "Start the OAuth authorization server",
		Long:  "Start the OAuth authorization server for user authentication and OAuth 2.0 flows",
		RunE:  serveAuth(rootCmd),
	}
}

func serveAuth(rootCmd *cobra.Command) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		configPath, _ := rootCmd.PersistentFlags().GetString("config")
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		if err := validateAuthConfig(cfg); err != nil {
			return err
		}

		c, err := container.CreateContainer(ctx, cfg)
		if err != nil {
			return fmt.Errorf("create container: %w", err)
		}
		if !c.IsHealthy(ctx) {
			return fmt.Errorf("container is not healthy")
		}

		if c.GetJWTSigner() == nil {
			return fmt.Errorf("JWT signer not initialized - check jwt key paths in config")
		}
		if c.GetSessionStore() == nil {
			return fmt.Errorf("session store not initialized - check sessionSecret in config")
		}

		srv := authserver.NewServer(c)
		handler := srv.Bootstrap()

		httpSrv := httpserver.NewHTTPServer(
			httpserver.WithHandler(handler),
			httpserver.WithPort(cfg.GetAuthPort()),
			httpserver.WithReadTimeout(cfg.GetServerReadTimeout()),
			httpserver.WithWriteTimeout(cfg.GetServerWriteTimeout()),
			httpserver.WithIdleTimeout(cfg.GetServerIdleTimeout()),
			httpserver.WithCleanupTimeout(cfg.GetServerCleanupTimeout()),
		)

		go func() {
			log.Printf("🚀 starting OAuth authorization server at port: %d\n", cfg.GetAuthPort())
			httpSrv.Start()
		}()

		defer cleanup(cfg,
			func() error {
				if err := httpSrv.Stop(); err != nil {
					log.Printf("failed shutdown auth server: %v\n", err)
					return err
				}
				return nil
			},
			func() error {
				if err := c.Shutdown(); err != nil {
					log.Printf("failed to shutdown container: %v\n", err)
					return err
				}
				return nil
			},
		)

		select {
		case <-ctx.Done():
			time.Sleep(100 * time.Millisecond)
			log.Println("performing graceful shutdown...")
		case err := <-httpSrv.Notify():
			if err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("auth server listen error: %w", err)
			}
		}

		return nil
	}
}

func validateAuthConfig(cfg altalune.Config) error {
	if cfg.GetJWTPrivateKeyPath() == "" {
		return fmt.Errorf("jwt private key path is required (security.jwtPrivateKeyPath)")
	}
	if cfg.GetJWTPublicKeyPath() == "" {
		return fmt.Errorf("jwt public key path is required (security.jwtPublicKeyPath)")
	}
	if cfg.GetSessionSecret() == "" {
		return fmt.Errorf("session secret is required (auth.sessionSecret)")
	}
	return nil
}
