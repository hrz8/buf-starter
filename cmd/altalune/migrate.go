package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/hrz8/altalune/internal/config"
	"github.com/hrz8/altalune/internal/container"
	"github.com/spf13/cobra"
)

func NewMigrateCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  "Run database migrations for the altalune service",
	}

	cmd.AddCommand(
		newMigrateUpCommand(rootCmd),
		newMigrateDownCommand(rootCmd),
		newMigrateStatusCommand(rootCmd),
	)

	return cmd
}

func newMigrateUpCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run database migrations (upgrade schema)",
		RunE:  runMigration(rootCmd, "up"),
	}
}

func newMigrateDownCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Rollback database migrations",
		RunE:  runMigration(rootCmd, "down"),
	}
}

func newMigrateStatusCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		RunE:  runMigration(rootCmd, "status"),
	}
}

func runMigration(rootCmd *cobra.Command, action string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Get and load configuration
		configPath, _ := rootCmd.PersistentFlags().GetString("config")
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Bootstrapping
		c, err := container.CreateContainer(ctx, cfg)
		if err != nil {
			log.Fatalf("failed to create application container: %v\n", err)
		}
		if !c.IsHealthy(ctx) {
			return fmt.Errorf("container is not healthy, cannot run migration")
		}

		defer c.Shutdown()

		log.Printf("Running migration %s on database", action)

		migrationSvc := c.GetMigrationService()

		if migrationSvc == nil {
			return errors.New("unexpected error: service unregistered")
		}

		switch action {
		case "up":
			return migrationSvc.MigrateUp(ctx)
		case "down":
			return migrationSvc.MigrateDown(ctx)
		case "status":
			return migrationSvc.MigrationStatus(ctx)
		default:
			return fmt.Errorf("unknown migration action: %s", action)
		}
	}
}
