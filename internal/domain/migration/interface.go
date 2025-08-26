package migration

import "context"

type Migrator interface {
	Down(ctx context.Context) error
	PrintStatus(ctx context.Context) error
	Up(ctx context.Context) error
}
