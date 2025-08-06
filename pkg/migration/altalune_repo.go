package migration

import (
	"context"
	"fmt"
	"sync"

	"github.com/hrz8/altalune/database"
	"github.com/hrz8/altalune/internal/postgres"
	"github.com/pressly/goose/v3"
)

type AltaluneMigrationRepo struct {
	db   postgres.DB
	once sync.Once
}

func NewAltaluneMigrationRepo(db postgres.DB) *AltaluneMigrationRepo {
	return &AltaluneMigrationRepo{db: db}
}

func (r *AltaluneMigrationRepo) configure() {
	r.once.Do(func() {
		goose.SetTableName(MigrationsTableName)
		goose.SetDialect(DatabaseDialect)
		goose.SetBaseFS(database.MigrationsFS)
	})
}

func (r *AltaluneMigrationRepo) Up(ctx context.Context) error {
	r.configure()
	db := r.db.GetDB()
	if db == nil {
		return fmt.Errorf("unknown database connection")
	}
	return goose.Up(db, MigrationsDir)
}

func (r *AltaluneMigrationRepo) Down(ctx context.Context) error {
	r.configure()
	db := r.db.GetDB()
	if db == nil {
		return fmt.Errorf("unknown database connection")
	}
	return goose.Down(db, MigrationsDir)
}

func (r *AltaluneMigrationRepo) PrintStatus(ctx context.Context) error {
	r.configure()
	db := r.db.GetDB()
	if db == nil {
		return fmt.Errorf("unknown database connection")
	}
	return goose.Status(db, MigrationsDir)
}
