package project

import (
	"context"
)

type Repositor interface {
	GetIDByPublicID(ctx context.Context, publicID string) (int64, error)
}
