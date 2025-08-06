package migration

import "context"

type AltaluneRepositor interface {
	Down(ctx context.Context) error
	PrintStatus(ctx context.Context) error
	Up(ctx context.Context) error
}
