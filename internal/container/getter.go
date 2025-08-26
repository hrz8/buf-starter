package container

import (
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
	migration_domain "github.com/hrz8/altalune/internal/domain/migration"
	"github.com/hrz8/altalune/internal/postgres"
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

// GetDB returns the database connection
func (c *Container) GetDBManager() postgres.Manager {
	if mgr, ok := c.db.(postgres.Manager); ok {
		return mgr
	}
	return nil
}

// GetMigrationService returns the migration service
func (c *Container) GetMigrationService() *migration_domain.Service {
	return c.migrationService
}

// GetGreeterService returns the greeter service (Only Example)
func (c *Container) GetGreeterService() greeterv1.GreeterServiceServer {
	return c.greeterService
}

// GetEmployeeService returns the employee service (Only Example)
func (c *Container) GetEmployeeService() altalunev1.EmployeeServiceServer {
	return c.employeeService
}

// GetProjectService returns the greeter service
func (c *Container) GetProjectService() altalunev1.ProjectServiceServer {
	return c.projectService
}
