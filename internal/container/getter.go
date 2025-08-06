package container

import (
	"github.com/hrz8/altalune"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
	"github.com/hrz8/altalune/internal/postgres"
	greeter_domain "github.com/hrz8/altalune/pkg/greeter"
	migration_domain "github.com/hrz8/altalune/pkg/migration"
)

// Public getter methods for accessing private components

// GetLogger returns the logger instance
func (c *Container) GetLogger() altalune.Logger {
	return c.logger
}

// GetConfig returns the configuration instance
func (c *Container) GetConfig() altalune.Config {
	return c.config
}

// GetDB returns the database connection
func (c *Container) GetDB() postgres.DB {
	return c.db
}

// GetGreeterRepo returns the greeter repository
func (c *Container) GetGreeterRepo() greeter_domain.Repositor {
	return c.greeterRepo
}

// GetMigrationRepo returns the migration repository
func (c *Container) GetMigrationRepo() migration_domain.AltaluneRepositor {
	return c.migrationRepo
}

// GetGreeterService returns the greeter service
func (c *Container) GetGreeterService() greeterv1.GreeterServiceServer {
	return c.greeterService
}

// GetMigrationService returns the migration service
func (c *Container) GetMigrationService() *migration_domain.Service {
	return c.migrationService
}
