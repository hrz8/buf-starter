package project

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hrz8/altalune/internal/postgres"
)

type Repo struct {
	db postgres.DB
}

func NewRepo(db postgres.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetIDByPublicID(ctx context.Context, publicID string) (int64, error) {
	query := `
		SELECT id 
		FROM altalune_projects 
		WHERE public_id = $1
	`

	var projectID int64
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrProjectNotFound
		}
		return 0, fmt.Errorf("get project ID by public ID: %w", err)
	}

	return projectID, nil
}
