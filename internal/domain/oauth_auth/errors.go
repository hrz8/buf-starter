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

	ErrMissingClientID            = errors.New("client_id is required")
	ErrMissingResponseType        = errors.New("response_type is required")
	ErrMissingRedirectURI         = errors.New("redirect_uri is required")
	ErrUnsupportedResponseType    = errors.New("unsupported response_type")
	ErrMissingCodeChallenge       = errors.New("code_challenge is required")
	ErrInvalidCodeChallengeMethod = errors.New("invalid code_challenge_method")
	ErrServerError                = errors.New("internal server error")

	ErrRefreshTokenExpired = errors.New("refresh token has expired")
	ErrRefreshTokenUsed    = errors.New("refresh token has already been used")
	ErrCodeExpired         = errors.New("authorization code has expired")
	ErrCodeAlreadyUsed     = errors.New("authorization code has already been used")

	// OTP errors
	ErrEmailNotRegistered = errors.New("email not registered")
	ErrOTPRateLimited     = errors.New("too many OTP requests, please try again later")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrOTPAlreadyUsed     = errors.New("OTP has already been used")

	// Email verification errors
	ErrInvalidVerificationToken = errors.New("invalid or expired verification token")
	ErrTokenAlreadyUsed         = errors.New("verification token has already been used")
	ErrUserNotFound             = errors.New("user not found")
)
