package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/hrz8/altalune/internal/config"
	"github.com/hrz8/altalune/internal/container"
	"github.com/hrz8/altalune/internal/domain/oauth_seeder"
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
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Run database migrations (upgrade schema)",
		RunE:  runMigration(rootCmd, "up"),
	}

	// Add --skip-seed flag for migrate up
	cmd.Flags().Bool("skip-seed", false, "Skip database seeding after migrations")

	return cmd
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
			// Run migrations first
			if err := migrationSvc.MigrateUp(ctx); err != nil {
				return err
			}

			// Check if seeding should be skipped
			skipSeed, _ := cmd.Flags().GetBool("skip-seed")
			if !skipSeed {
				log.Println("Running database seeder...")

				// Get database connection from container
				dbManager := c.GetDBManager()
				if dbManager == nil {
					return errors.New("database manager not available")
				}

				// Initialize seeder with config
				// Type assert to *config.AppConfig
				appCfg, ok := cfg.(*config.AppConfig)
				if !ok {
					return errors.New("failed to type assert config to *config.AppConfig")
				}

				seeder, err := oauth_seeder.NewSeeder(dbManager.GetDB(), appCfg)
				if err != nil {
					return fmt.Errorf("failed to initialize seeder: %w", err)
				}

				// Run seeder
				if err := seeder.Seed(ctx); err != nil {
					return fmt.Errorf("seeding failed: %w", err)
				}

				log.Println("Database seeding completed successfully")
			} else {
				log.Println("Skipping database seeding (--skip-seed flag set)")
			}

			return nil

		case "down":
			return migrationSvc.MigrateDown(ctx)
		case "status":
			return migrationSvc.MigrationStatus(ctx)
		default:
			return fmt.Errorf("unknown migration action: %s", action)
		}
	}
}
