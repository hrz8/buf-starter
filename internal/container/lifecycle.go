package container

import (
	"context"
	"fmt"
)

// Shutdown gracefully shuts down all components
func (c *Container) Shutdown() error {
	if c.GetDBManager() != nil {
		if err := c.GetDBManager().Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}
	return nil
}

// Health returns the health status of all components
func (c *Container) Health(ctx context.Context) map[string]error {
	health := make(map[string]error)

	// Check database health
	if c.db == nil {
		health["database"] = fmt.Errorf("database not initialized")
	} else if err := c.GetDBManager().TestConnection(ctx); err != nil {
		health["database"] = err
	} else {
		health["database"] = nil
	}

	// Add other health checks as needed
	health["logger"] = nil
	health["config"] = nil

	return health
}

// IsHealthy returns true if all components are healthy
func (c *Container) IsHealthy(ctx context.Context) bool {
	health := c.Health(ctx)
	for _, err := range health {
		if err != nil {
			return false
		}
	}
	return true
}
