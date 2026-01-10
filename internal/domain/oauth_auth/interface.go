package oauth_auth

import (
	"context"

	"github.com/google/uuid"
)

// Repositor defines the interface for OAuth auth repository operations.
type Repositor interface {
	CreateAuthorizationCode(ctx context.Context, input *CreateAuthCodeInput) (*AuthorizationCode, error)
	GetAuthorizationCodeByCode(ctx context.Context, code uuid.UUID) (*AuthorizationCode, error)
	MarkCodeExchanged(ctx context.Context, code uuid.UUID) error

	CreateRefreshToken(ctx context.Context, input *CreateRefreshTokenInput) (*RefreshToken, error)
	GetRefreshTokenByToken(ctx context.Context, token uuid.UUID) (*RefreshToken, error)
	MarkRefreshTokenExchanged(ctx context.Context, token uuid.UUID) error

	GetUserConsent(ctx context.Context, userID int64, clientID uuid.UUID) (*UserConsent, error)
	CreateOrUpdateUserConsent(ctx context.Context, input *UserConsentInput) (*UserConsent, error)

	GetOAuthClientByClientID(ctx context.Context, clientID uuid.UUID) (*OAuthClientInfo, error)
}
