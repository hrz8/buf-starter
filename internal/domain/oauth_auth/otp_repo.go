package oauth_auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hrz8/altalune/internal/postgres"
)

// OTPRepo implements OTPRepositor for database operations on OTP tokens.
type OTPRepo struct {
	db postgres.DB
}

// NewOTPRepo creates a new OTP repository.
func NewOTPRepo(db postgres.DB) *OTPRepo {
	return &OTPRepo{db: db}
}

// CreateOTP stores a new OTP token hash in the database.
func (r *OTPRepo) CreateOTP(ctx context.Context, email, otpHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO altalune_otp_tokens (email, otp_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, email, otpHash, expiresAt)
	if err != nil {
		return fmt.Errorf("create OTP: %w", err)
	}
	return nil
}

// GetValidOTP retrieves a valid (unused, not expired) OTP by email and hash.
func (r *OTPRepo) GetValidOTP(ctx context.Context, email, otpHash string) (*OTPToken, error) {
	query := `
		SELECT id, email, otp_hash, expires_at, used_at, created_at
		FROM altalune_otp_tokens
		WHERE email = $1 AND otp_hash = $2 AND used_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`
	var otp OTPToken
	err := r.db.QueryRowContext(ctx, query, email, otpHash).Scan(
		&otp.ID, &otp.Email, &otp.OTPHash, &otp.ExpiresAt, &otp.UsedAt, &otp.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidOTP
		}
		return nil, fmt.Errorf("get valid OTP: %w", err)
	}
	return &otp, nil
}

// MarkOTPUsed marks an OTP token as used.
func (r *OTPRepo) MarkOTPUsed(ctx context.Context, id int64) error {
	query := `UPDATE altalune_otp_tokens SET used_at = NOW() WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark OTP used: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrInvalidOTP
	}
	return nil
}

// CountRecentOTPs counts OTPs sent to an email within a time window (for rate limiting).
func (r *OTPRepo) CountRecentOTPs(ctx context.Context, email string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*) FROM altalune_otp_tokens
		WHERE email = $1 AND created_at > $2
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, email, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count recent OTPs: %w", err)
	}
	return count, nil
}
