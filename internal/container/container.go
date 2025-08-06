package container

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/logger"
	greeter_domain "github.com/hrz8/altalune/pkg/greeter"
	migration_domain "github.com/hrz8/altalune/pkg/migration"
)

// Container manages dependency injection with private fields
type Container struct {
	// Configuration and logger
	config altalune.Config
	logger altalune.Logger

	// Database connection and manager
	db        postgres.DB
	dbManager postgres.Manager

	// Repositories
	greeterRepo   greeter_domain.Repositor
	migrationRepo migration_domain.AltaluneRepositor

	// Services
	greeterService   greeterv1.GreeterServiceServer
	migrationService *migration_domain.Service
}

// CreateContainer creates a new dependency injection container with proper error handling
func CreateContainer(ctx context.Context, cfg altalune.Config) (*Container, error) {
	container := &Container{
		config: cfg,
		logger: logger.New(cfg.GetServerLogLevel()),
	}

	// Initialize components in dependency order
	if err := container.initDatabase(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	if err := container.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}
	if err := container.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}
	return container, nil
}

// Private initialization methods
func (c *Container) initDatabase(ctx context.Context) error {
	conn := postgres.MustConnect(postgres.ConnectionOptions{
		URL:            c.config.GetDatabaseURL(),
		MaxConnections: c.config.GetDatabaseMaxConnections(),
		MaxIdleTime:    c.config.GetDatabaseMaxIdleTime(),
		ConnectTimeout: c.config.GetDatabaseConnectTimeout(),
	})
	if err := conn.TestConnection(ctx); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}
	c.db = conn
	c.dbManager = conn
	return nil
}

func (c *Container) initRepositories() error {
	c.greeterRepo = greeter_domain.NewRepo()
	c.migrationRepo = migration_domain.NewAltaluneMigrationRepo(c.db)
	return nil
}

func (c *Container) initServices() error {
	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}
	c.greeterService = greeter_domain.NewService(validator, c.logger, c.greeterRepo)
	c.migrationService = migration_domain.NewService(c.logger, c.migrationRepo)
	return nil
}
