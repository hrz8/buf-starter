package oauth_auth

import "errors"

var (
	ErrAuthorizationCodeNotFound = errors.New("authorization code not found or expired")
	ErrRefreshTokenNotFound      = errors.New("refresh token not found or expired")
	ErrUserConsentNotFound       = errors.New("user consent not found")
	ErrOAuthClientNotFound       = errors.New("oauth client not found")
	ErrInvalidAuthorizationCode  = errors.New("invalid authorization code")
	ErrClientMismatch            = errors.New("client_id does not match authorization code")
	ErrRedirectURIMismatch       = errors.New("redirect_uri does not match authorization code")
	ErrMissingCodeVerifier       = errors.New("code_verifier required for PKCE")
	ErrInvalidCodeVerifier       = errors.New("code_verifier does not match code_challenge")
	ErrInvalidClientID           = errors.New("invalid client_id format")
	ErrInvalidClientSecret       = errors.New("invalid client_secret")
	ErrInvalidRefreshToken       = errors.New("invalid refresh token")
	ErrClientSecretRequired      = errors.New("client_secret required")
	ErrPKCERequired              = errors.New("PKCE is required for this client")
)
