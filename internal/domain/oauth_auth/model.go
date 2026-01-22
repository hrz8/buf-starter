package oauth_auth

import (
	"time"

	"github.com/google/uuid"
)

// AuthorizationCode represents an OAuth authorization code.
type AuthorizationCode struct {
	ID                  int64
	Code                uuid.UUID
	ClientID            uuid.UUID
	UserID              int64
	RedirectURI         string
	Scope               string
	Nonce               *string
	CodeChallenge       *string
	CodeChallengeMethod *string
	ExpiresAt           time.Time
	ExchangeAt          *time.Time
	CreatedAt           time.Time
}

// RefreshToken represents an OAuth refresh token.
type RefreshToken struct {
	ID         int64
	Token      uuid.UUID
	ClientID   uuid.UUID
	UserID     int64
	Scope      string
	Nonce      *string
	ExpiresAt  time.Time
	ExchangeAt *time.Time
	CreatedAt  time.Time
}

// UserConsent represents a user's consent to an OAuth client.
type UserConsent struct {
	ID        int64
	UserID    int64
	ClientID  uuid.UUID
	Scope     string
	GrantedAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// UserConsentWithClient represents a user's consent with OAuth client details.
type UserConsentWithClient struct {
	ID         int64
	UserID     int64
	ClientID   uuid.UUID
	ClientName string
	Scope      string
	GrantedAt  time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
}

// CreateAuthCodeInput holds parameters for creating an authorization code.
type CreateAuthCodeInput struct {
	ClientID            uuid.UUID
	UserID              int64
	RedirectURI         string
	Scope               string
	Nonce               *string
	CodeChallenge       *string
	CodeChallengeMethod *string
	ExpiresAt           time.Time
}

// CreateRefreshTokenInput holds parameters for creating a refresh token.
type CreateRefreshTokenInput struct {
	ClientID  uuid.UUID
	UserID    int64
	Scope     string
	Nonce     *string
	ExpiresAt time.Time
}

// UserConsentInput holds parameters for creating or updating user consent.
type UserConsentInput struct {
	UserID   int64
	ClientID uuid.UUID
	Scope    string
}

// CodeExchangeResult holds the result of exchanging an authorization code.
type CodeExchangeResult struct {
	UserID int64
	Scope  string
	Nonce  *string
}

// TokenPair holds an access token and refresh token pair.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	Scope        string
}

// GenerateAuthCodeInput holds parameters for generating an authorization code.
type GenerateAuthCodeInput struct {
	ClientID            uuid.UUID
	UserID              int64
	RedirectURI         string
	Scope               string
	Nonce               *string
	CodeChallenge       *string
	CodeChallengeMethod *string
}

// OAuthClientInfo holds OAuth client information for authentication.
type OAuthClientInfo struct {
	ID           int64
	ClientID     uuid.UUID
	Name         string
	RedirectURIs []string
	PKCERequired bool
	IsDefault    bool
	SecretHash   *string // Nullable for public clients
	Confidential bool
}

// OTPToken represents a one-time password token for authentication.
type OTPToken struct {
	ID        int64
	Email     string
	OTPHash   string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// EmailVerificationToken represents a token for verifying user email addresses.
type EmailVerificationToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// UserInfo holds minimal user information for OTP/verification services.
type UserInfo struct {
	ID            int64
	PublicID      string
	Email         string
	FirstName     string
	LastName      string
	IsActive      bool
	EmailVerified bool
}
