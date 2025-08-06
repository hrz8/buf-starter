package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGXConnection struct {
	_ *sql.DB
	_ ConnectionOptions
}

var _ DB = (*PGXConnection)(nil)

// MustConnect creates a new database connection manager
func MustConnectWithPGX(cfg ConnectionOptions) *PGXConnection {
	return nil
}

// TestConnection tests database connectivity
func (cm *PGXConnection) TestConnection(ctx context.Context) error {
	panic("unimplemented")
}

// GetDB implements DB.
func (cm *PGXConnection) GetDB() *sql.DB {
	panic("unimplemented")
}

// QueryContext implements DB.
func (cm *PGXConnection) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	panic("unimplemented")
}

// QueryRowContext implements DB.
func (cm *PGXConnection) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	panic("unimplemented")
}

// ExecContext implements DB.
func (cm *PGXConnection) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	panic("unimplemented")
}

// PingContext implements DB.
func (cm *PGXConnection) PingContext(ctx context.Context) error {
	panic("unimplemented")
}

// Close implements DB.
func (cm *PGXConnection) Close() error {
	panic("unimplemented")
}
