# Task T26: serve-auth CLI Command & HTTP Server Setup

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Backend Infrastructure
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T24 (Session & JWT utilities), T25 (OAuth Auth Domain)

## Objective

Create the new `serve-auth` CLI command that starts the OAuth authorization server on a separate port from the main API server, with session management, proper middleware chain, and graceful shutdown.

## Acceptance Criteria

- [ ] `./bin/app serve-auth -c config.yaml` starts OAuth server
- [ ] Server runs on configurable port (default: 3101)
- [ ] Reads same config.yaml as `serve` command
- [ ] Initializes database connection (shared pool pattern)
- [ ] Loads RSA keys for JWT signing
- [ ] Initializes Gorilla session store
- [ ] Chi router with middleware chain
- [ ] Graceful shutdown on SIGTERM/SIGINT
- [ ] Startup logging includes port and host

## Technical Requirements

### CLI Command (`cmd/altalune/serve_auth.go`)

```go
func NewServeAuthCommand(rootCmd *cobra.Command) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "serve-auth",
        Short: "Start the OAuth authorization server",
        Long:  "Start the OAuth authorization server for user authentication and OAuth 2.0 flows",
        RunE:  serveAuth(rootCmd),
    }
    return cmd
}
```

Command should:
1. Load configuration from same config.yaml
2. Create database connection (same pattern as serve command)
3. Create auth server container with dependencies
4. Start HTTP server with graceful shutdown

### Auth Server Container (`internal/authserver/container.go`)

```go
type Container struct {
    config       altalune.Config
    logger       *slog.Logger
    db           postgres.Manager
    sessionStore *session.Store
    jwtSigner    *jwt.Signer

    // Domain services
    oauthAuthService   *oauth_auth.Service
    oauthClientService *oauth_client.Service
    oauthProviderService *oauth_provider.Service
    userService        *user.Service
}
```

Container should initialize:
- Session store (from T24)
- JWT signer (from T24)
- OAuth auth service (from T25)
- Reuse existing domain services (oauth_client, oauth_provider, user)

### HTTP Server (`internal/authserver/server.go`)

```go
type Server struct {
    container  *Container
    httpServer *http.Server
    router     chi.Router
}

func NewServer(container *Container) *Server {
    r := chi.NewRouter()

    // Apply middleware
    r.Use(middleware.RecoveryMiddleware(container.Logger()))
    r.Use(middleware.LoggingMiddleware(container.Logger()))
    r.Use(session.LoadSessionMiddleware(container.SessionStore()))

    // Register routes
    registerRoutes(r, container)

    return &Server{
        container: container,
        router:    r,
    }
}

func (s *Server) Start(ctx context.Context, host string, port int) error {
    addr := fmt.Sprintf("%s:%d", host, port)
    s.httpServer = &http.Server{
        Addr:    addr,
        Handler: s.router,
    }

    s.container.Logger().Info("Starting OAuth authorization server", "addr", addr)

    return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.httpServer.Shutdown(ctx)
}
```

### Routes (`internal/authserver/routes.go`)

Initial placeholder routes (actual handlers implemented in subsequent tasks):

```go
func registerRoutes(r chi.Router, c *Container) {
    // Health check (always available)
    r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    // Login pages (public)
    r.Get("/login", handlers.LoginPage(c))
    r.Get("/login/{provider}", handlers.LoginProvider(c))
    r.Get("/auth/callback", handlers.OAuthCallback(c))

    // Logout (requires session)
    r.Post("/logout", handlers.Logout(c))

    // OAuth endpoints
    r.Route("/oauth", func(r chi.Router) {
        r.Get("/authorize", handlers.Authorize(c))
        r.Post("/authorize", handlers.AuthorizeProcess(c))
        r.Post("/token", handlers.Token(c))
    })

    // JWKS (public)
    r.Get("/.well-known/jwks.json", handlers.JWKS(c))
}
```

### Middleware Chain

Order of middleware execution:
1. **Recovery** - Catches panics, returns 500
2. **Logging** - Logs request/response
3. **Session Loading** - Loads session from cookie into context

### Configuration

The auth server uses these config sections:

```yaml
auth:
  host: localhost
  port: 3101
  sessionSecret: "your-32-byte-session-secret"
  codeExpiry: 600          # Authorization code TTL (seconds)
  accessTokenExpiry: 3600  # Access token TTL (seconds)
  refreshTokenExpiry: 2592000  # Refresh token TTL (seconds)

security:
  jwtPrivateKeyPath: "keys/rsa-private.pem"
  jwtPublicKeyPath: "keys/rsa-public.pem"
  jwksKid: "altalune-oauth-2024"
```

## Implementation Details

### Command Implementation

```go
// cmd/altalune/serve_auth.go
func serveAuth(rootCmd *cobra.Command) func(cmd *cobra.Command, args []string) error {
    return func(cmd *cobra.Command, args []string) error {
        ctx := cmd.Context()

        // Load config
        configPath, _ := rootCmd.PersistentFlags().GetString("config")
        cfg, err := config.Load(configPath)
        if err != nil {
            return fmt.Errorf("load config: %w", err)
        }

        // Validate required config
        if cfg.GetJWTPrivateKeyPath() == "" {
            return fmt.Errorf("jwt private key path is required")
        }
        if cfg.GetAuthSessionSecret() == "" {
            return fmt.Errorf("auth session secret is required")
        }

        // Create database connection
        db, err := postgres.MustConnect(postgres.ConnectionOptions{
            URL:            cfg.GetDatabaseURL(),
            MaxConnections: cfg.GetDatabaseMaxConnections(),
            MaxIdleTime:    cfg.GetDatabaseMaxIdleTime(),
            ConnectTimeout: cfg.GetDatabaseConnectTimeout(),
        })
        if err != nil {
            return fmt.Errorf("connect database: %w", err)
        }
        defer db.Close()

        // Create auth server container
        container, err := authserver.NewContainer(ctx, cfg, db)
        if err != nil {
            return fmt.Errorf("create container: %w", err)
        }

        // Create and start server
        server := authserver.NewServer(container)

        // Graceful shutdown
        errChan := make(chan error, 1)
        go func() {
            errChan <- server.Start(ctx, cfg.GetAuthHost(), cfg.GetAuthPort())
        }()

        // Wait for interrupt or error
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

        select {
        case err := <-errChan:
            if err != nil && err != http.ErrServerClosed {
                return err
            }
        case <-sigChan:
            container.Logger().Info("Shutting down OAuth server...")
            shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            defer cancel()
            return server.Shutdown(shutdownCtx)
        }

        return nil
    }
}
```

### Container Initialization

```go
// internal/authserver/container.go
func NewContainer(ctx context.Context, cfg altalune.Config, db postgres.Manager) (*Container, error) {
    logger := slog.Default()

    // Initialize session store
    sessionStore := session.NewStore(
        cfg.GetAuthSessionSecret(),
        cfg.IsProduction(),
        cfg.GetAuthSessionMaxAge(),
    )

    // Load JWT signer
    jwtSigner, err := jwt.NewSigner(
        cfg.GetJWTPrivateKeyPath(),
        cfg.GetJWTPublicKeyPath(),
        cfg.GetJWKSKid(),
    )
    if err != nil {
        return nil, fmt.Errorf("initialize jwt signer: %w", err)
    }

    // Initialize domain services (reuse patterns from main container)
    // ...

    return &Container{
        config:       cfg,
        logger:       logger,
        db:           db,
        sessionStore: sessionStore,
        jwtSigner:    jwtSigner,
        // ... services
    }, nil
}
```

## Files to Create

- `cmd/altalune/serve_auth.go` - CLI command implementation
- `internal/authserver/server.go` - HTTP server lifecycle
- `internal/authserver/routes.go` - Route registration
- `internal/authserver/container.go` - Dependency injection container

## Files to Modify

- `cmd/altalune/main.go` - Register serve-auth command
- `config.go` (if needed) - Add auth config getters

## Testing Requirements

- Manual testing: `./bin/app serve-auth -c config.yaml`
- Verify server starts on correct port
- Verify /healthz endpoint responds
- Verify graceful shutdown works

## Commands to Run

```bash
# Build application
make build

# Start auth server
./bin/app serve-auth -c config.yaml

# In another terminal, test health endpoint
curl http://localhost:3101/healthz
# Expected: ok

# Test graceful shutdown
# Press Ctrl+C in server terminal
# Server should log "Shutting down OAuth server..."
```

## Validation Checklist

- [ ] Command registered and shows in `./bin/app --help`
- [ ] Server starts on configured port
- [ ] Health endpoint responds with 200 OK
- [ ] Session store initialized (check logs)
- [ ] JWT signer initialized (check logs)
- [ ] Graceful shutdown completes within timeout
- [ ] Database connection is established

## Definition of Done

- [ ] serve-auth command implemented and registered
- [ ] Server starts on port 3101 (or configured port)
- [ ] Database connection established (shared pool)
- [ ] Session store initialized with config secret
- [ ] JWT signer loads RSA keys successfully
- [ ] Chi router configured with middleware chain
- [ ] Routes registered (placeholder handlers OK)
- [ ] Health endpoint working
- [ ] Graceful shutdown handles SIGTERM/SIGINT
- [ ] Startup logging shows host:port

## Dependencies

- T24: Session store and JWT signer utilities
- T25: OAuth auth domain service
- Existing: oauth_client, oauth_provider, user domain services
- `github.com/go-chi/chi/v5` - HTTP router

## Risk Factors

- **Low Risk**: Similar to existing serve command pattern
- **Medium Risk**: RSA key loading - must handle file not found gracefully

## Notes

- The auth server reuses the same database connection pool as the main server
- Both servers can run simultaneously (different ports: 3100 and 3101)
- Session cookies are scoped to the auth server domain
- In production, auth server should be behind HTTPS
- Consider adding metrics endpoint in future (/metrics for Prometheus)
