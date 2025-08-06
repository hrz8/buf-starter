package migration

import (
	"context"
	"fmt"

	"github.com/hrz8/altalune"
)

type Service struct {
	log        altalune.Logger
	gowapiRepo AltaluneRepositor
}

func NewService(
	log altalune.Logger,
	gowapiRepo AltaluneRepositor,
) *Service {
	return &Service{
		log:        log,
		gowapiRepo: gowapiRepo,
	}
}

// MigrateUp runs database migrations (schema upgrade)
func (s *Service) MigrateUp(ctx context.Context) error {
	s.log.Info("Starting database migration up...")

	s.log.Info("Running altalune database migration")
	if err := s.gowapiRepo.Up(ctx); err != nil {
		return fmt.Errorf("failed to run altalune migration: %w", err)
	}

	s.log.Info("Database migration up completed successfully")
	return nil
}

// MigrateDown runs database migrations (schema downgrade)
func (s *Service) MigrateDown(ctx context.Context) error {
	s.log.Info("Starting database migration down...")

	s.log.Info("Running altalune database migration")
	if err := s.gowapiRepo.Down(ctx); err != nil {
		return fmt.Errorf("failed to run altalune migration: %w", err)
	}

	s.log.Info("Database migration down completed successfully")
	return nil
}

// MigrationStatus is to print current migration statuses
func (s *Service) MigrationStatus(ctx context.Context) error {
	if err := s.gowapiRepo.PrintStatus(ctx); err != nil {
		return fmt.Errorf("failed to run altalune migration: %w", err)
	}
	return nil
}
