package oauth_auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hrz8/altalune/internal/postgres"
)

// EmailVerificationRepo implements EmailVerificationRepositor for database operations.
type EmailVerificationRepo struct {
	db postgres.DB
}

// NewEmailVerificationRepo creates a new email verification repository.
func NewEmailVerificationRepo(db postgres.DB) *EmailVerificationRepo {
	return &EmailVerificationRepo{db: db}
}

// CreateVerificationToken stores a new email verification token hash in the database.
func (r *EmailVerificationRepo) CreateVerificationToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO altalune_email_verification_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("create verification token: %w", err)
	}
	return nil
}

// GetValidToken retrieves a valid (unused, not expired) verification token by hash.
func (r *EmailVerificationRepo) GetValidToken(ctx context.Context, tokenHash string) (*EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM altalune_email_verification_tokens
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()
	`
	var token EmailVerificationToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidVerificationToken
		}
		return nil, fmt.Errorf("get valid token: %w", err)
	}
	return &token, nil
}

// MarkTokenUsed marks a verification token as used.
func (r *EmailVerificationRepo) MarkTokenUsed(ctx context.Context, id int64) error {
	query := `UPDATE altalune_email_verification_tokens SET used_at = NOW() WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark token used: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrInvalidVerificationToken
	}
	return nil
}

// InvalidateUserTokens marks all unused tokens for a user as used (invalidates them).
func (r *EmailVerificationRepo) InvalidateUserTokens(ctx context.Context, userID int64) error {
	query := `UPDATE altalune_email_verification_tokens SET used_at = NOW() WHERE user_id = $1 AND used_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("invalidate user tokens: %w", err)
	}
	return nil
}
