package oauth_auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hrz8/altalune/internal/postgres"
)

// UserRepo implements UserLookupRepositor and UserEmailVerificationRepositor.
type UserRepo struct {
	db postgres.DB
}

// NewUserRepo creates a new user repository for OTP/verification services.
func NewUserRepo(db postgres.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetUserByEmail retrieves user info by email address.
func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*UserInfo, error) {
	query := `
		SELECT id, public_id, email, first_name, last_name, is_active, email_verified
		FROM altalune_users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`
	var user UserInfo
	var firstName, lastName sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&firstName,
		&lastName,
		&user.IsActive,
		&user.EmailVerified,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if firstName.Valid {
		user.FirstName = firstName.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}

	return &user, nil
}

// GetUserByPublicID retrieves user info by public ID (UUID string).
func (r *UserRepo) GetUserByPublicID(ctx context.Context, publicID string) (*UserInfo, error) {
	query := `
		SELECT id, public_id, email, first_name, last_name, is_active, email_verified
		FROM altalune_users
		WHERE public_id = $1
	`
	var user UserInfo
	var firstName, lastName sql.NullString

	err := r.db.QueryRowContext(ctx, query, publicID).Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&firstName,
		&lastName,
		&user.IsActive,
		&user.EmailVerified,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by public ID: %w", err)
	}

	if firstName.Valid {
		user.FirstName = firstName.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}

	return &user, nil
}

// GetUserByID retrieves user info by internal database ID.
func (r *UserRepo) GetUserByID(ctx context.Context, userID int64) (*UserInfo, error) {
	query := `
		SELECT id, public_id, email, first_name, last_name, is_active, email_verified
		FROM altalune_users
		WHERE id = $1
	`
	var user UserInfo
	var firstName, lastName sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&firstName,
		&lastName,
		&user.IsActive,
		&user.EmailVerified,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	if firstName.Valid {
		user.FirstName = firstName.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}

	return &user, nil
}

// SetEmailVerified updates the email_verified status and activated_at timestamp for a user.
func (r *UserRepo) SetEmailVerified(ctx context.Context, userID int64, verified bool) error {
	var query string
	if verified {
		// When verifying, also set activated_at if not already set
		query = `
			UPDATE altalune_users
			SET email_verified = $2,
			    activated_at = COALESCE(activated_at, NOW()),
			    updated_at = NOW()
			WHERE id = $1
		`
	} else {
		query = `
			UPDATE altalune_users
			SET email_verified = $2, updated_at = NOW()
			WHERE id = $1
		`
	}

	result, err := r.db.ExecContext(ctx, query, userID, verified)
	if err != nil {
		return fmt.Errorf("set email verified: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
