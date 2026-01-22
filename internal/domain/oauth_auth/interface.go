package oauth_auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/altalune/internal/domain/permission"
)

// UserPermissionProvider defines the interface for fetching user permissions.
type UserPermissionProvider interface {
	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
}

// IAMMapperRepositor defines the interface for fetching user permissions.
type IAMMapperRepositor interface {
	GetUserPermissions(ctx context.Context, userID int64) ([]*permission.Permission, error)
}

// Repositor defines the interface for OAuth auth repository operations.
type Repositor interface {
	CreateAuthorizationCode(ctx context.Context, input *CreateAuthCodeInput) (*AuthorizationCode, error)
	GetAuthorizationCodeByCode(ctx context.Context, code uuid.UUID) (*AuthorizationCode, error)
	MarkCodeExchanged(ctx context.Context, code uuid.UUID) error

	CreateRefreshToken(ctx context.Context, input *CreateRefreshTokenInput) (*RefreshToken, error)
	GetRefreshTokenByToken(ctx context.Context, token uuid.UUID) (*RefreshToken, error)
	MarkRefreshTokenExchanged(ctx context.Context, token uuid.UUID) error

	GetUserConsent(ctx context.Context, userID int64, clientID uuid.UUID) (*UserConsent, error)
	GetUserConsents(ctx context.Context, userID int64) ([]*UserConsentWithClient, error)
	CreateOrUpdateUserConsent(ctx context.Context, input *UserConsentInput) (*UserConsent, error)
	RevokeUserConsent(ctx context.Context, userID int64, clientID uuid.UUID) error

	GetOAuthClientByClientID(ctx context.Context, clientID uuid.UUID) (*OAuthClientInfo, error)
}

// OTPRepositor defines the interface for OTP repository operations.
type OTPRepositor interface {
	CreateOTP(ctx context.Context, email, otpHash string, expiresAt time.Time) error
	GetValidOTP(ctx context.Context, email, otpHash string) (*OTPToken, error)
	MarkOTPUsed(ctx context.Context, id int64) error
	CountRecentOTPs(ctx context.Context, email string, since time.Time) (int, error)
}

// EmailVerificationRepositor defines the interface for email verification repository operations.
type EmailVerificationRepositor interface {
	CreateVerificationToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	GetValidToken(ctx context.Context, tokenHash string) (*EmailVerificationToken, error)
	MarkTokenUsed(ctx context.Context, id int64) error
	InvalidateUserTokens(ctx context.Context, userID int64) error
}

// UserLookupRepositor defines the interface for looking up users by email, public ID, or internal ID (for OTP service and introspection).
type UserLookupRepositor interface {
	GetUserByEmail(ctx context.Context, email string) (*UserInfo, error)
	GetUserByPublicID(ctx context.Context, publicID string) (*UserInfo, error)
	GetUserByID(ctx context.Context, userID int64) (*UserInfo, error)
}

// UserEmailVerificationRepositor defines the interface for user email verification operations.
type UserEmailVerificationRepositor interface {
	GetUserByID(ctx context.Context, userID int64) (*UserInfo, error)
	SetEmailVerified(ctx context.Context, userID int64, verified bool) error
}
