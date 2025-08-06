package postgres

import (
	"context"
	"database/sql"
)

// DB defines a minimal interface compatible with both *sql.DB and *pgxpool.Pool
type DB interface {
	GetDB() *sql.DB
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PingContext(ctx context.Context) error
}

// Manager defines an interface for managing database connections
type Manager interface {
	DB
	TestConnection(ctx context.Context) error
	Close() error
}
