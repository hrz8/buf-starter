package oauth_auth

import "errors"

var (
	ErrInvalidAuthorizationCode     = errors.New("invalid authorization code")
	ErrAuthorizationCodeExpired     = errors.New("authorization code has expired")
	ErrAuthorizationCodeAlreadyUsed = errors.New("authorization code already used")
	ErrInvalidRefreshToken          = errors.New("invalid refresh token")
	ErrRefreshTokenExpired          = errors.New("refresh token has expired")
	ErrRefreshTokenAlreadyUsed      = errors.New("refresh token already used")
	ErrInvalidPKCEVerifier          = errors.New("invalid PKCE code verifier")
	ErrPKCERequired                 = errors.New("PKCE is required for this client")
	ErrInvalidScope                 = errors.New("invalid scope requested")
	ErrUserConsentRequired          = errors.New("user consent required")
	ErrInvalidClientCredentials     = errors.New("invalid client credentials")
)
