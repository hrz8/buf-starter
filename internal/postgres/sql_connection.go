package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLConnection struct {
	db     *sql.DB
	config ConnectionOptions
}

var _ DB = (*SQLConnection)(nil)

// MustConnect creates a new database connection manager
func MustConnect(cfg ConnectionOptions) *SQLConnection {
	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		panic(fmt.Errorf("failed opening database connection: %w", err))
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections / 2)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	return &SQLConnection{
		db:     db,
		config: cfg,
	}
}

// TestConnection tests database connectivity
func (cm *SQLConnection) TestConnection(ctx context.Context) error {
	connCtx, cancel := context.WithTimeout(ctx, cm.config.ConnectTimeout)
	defer cancel()

	if err := cm.db.PingContext(connCtx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	var result int
	err := cm.db.QueryRowContext(connCtx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("failed to execute test query: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected test query result: expected 1, got %d", result)
	}

	return nil
}

// GetDB returns the database connection
func (cm *SQLConnection) GetDB() *sql.DB {
	return cm.db
}

// QueryContext implements DB
func (c *SQLConnection) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

// QueryRowContext implements DB
func (c *SQLConnection) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return c.db.QueryRowContext(ctx, query, args...)
}

// ExecContext implements DB.
func (c *SQLConnection) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

// PingContext implements DB.
func (c *SQLConnection) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Close implements DB.
func (c *SQLConnection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
