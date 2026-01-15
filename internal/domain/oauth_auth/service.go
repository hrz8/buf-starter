package oauth_auth

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/shared/jwt"
	"github.com/hrz8/altalune/internal/shared/password"
	"github.com/hrz8/altalune/internal/shared/pkce"
)

// Service handles OAuth authorization code and token operations.
type Service struct {
	repo      Repositor
	jwtSigner *jwt.Signer
	cfg       altalune.Config
	log       altalune.Logger
}

// NewService creates a new OAuth auth service.
func NewService(log altalune.Logger, repo Repositor, jwtSigner *jwt.Signer, cfg altalune.Config) *Service {
	return &Service{
		repo:      repo,
		jwtSigner: jwtSigner,
		cfg:       cfg,
		log:       log,
	}
}

// GenerateAuthorizationCode creates a new authorization code with the configured expiry.
func (s *Service) GenerateAuthorizationCode(ctx context.Context, input *GenerateAuthCodeInput) (*AuthorizationCode, error) {
	expiresAt := time.Now().Add(time.Duration(s.cfg.GetCodeExpiry()) * time.Second)

	createInput := &CreateAuthCodeInput{
		ClientID:            input.ClientID,
		UserID:              input.UserID,
		RedirectURI:         input.RedirectURI,
		Scope:               input.Scope,
		Nonce:               input.Nonce,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		ExpiresAt:           expiresAt,
	}

	code, err := s.repo.CreateAuthorizationCode(ctx, createInput)
	if err != nil {
		s.log.Error("failed to create authorization code",
			"error", err,
			"client_id", input.ClientID,
			"user_id", input.UserID,
		)
		return nil, err
	}

	return code, nil
}

// ValidateAndExchangeCode validates and exchanges an authorization code for tokens.
func (s *Service) ValidateAndExchangeCode(ctx context.Context, codeStr string, clientID uuid.UUID, redirectURI string, codeVerifier *string) (*CodeExchangeResult, error) {
	code, err := uuid.Parse(codeStr)
	if err != nil {
		return nil, ErrInvalidAuthorizationCode
	}

	authCode, err := s.repo.GetAuthorizationCodeByCode(ctx, code)
	if err != nil {
		return nil, ErrInvalidAuthorizationCode
	}

	if authCode.ClientID != clientID {
		return nil, ErrClientMismatch
	}

	if authCode.RedirectURI != redirectURI {
		return nil, ErrRedirectURIMismatch
	}

	if authCode.CodeChallenge != nil && *authCode.CodeChallenge != "" {
		if codeVerifier == nil || *codeVerifier == "" {
			return nil, ErrMissingCodeVerifier
		}
		method := "S256"
		if authCode.CodeChallengeMethod != nil {
			method = *authCode.CodeChallengeMethod
		}
		if !pkce.VerifyCodeChallenge(*codeVerifier, *authCode.CodeChallenge, method) {
			return nil, ErrInvalidCodeVerifier
		}
	}

	if err := s.repo.MarkCodeExchanged(ctx, code); err != nil {
		s.log.Error("failed to mark code exchanged",
			"error", err,
			"code", code,
		)
		return nil, err
	}

	return &CodeExchangeResult{
		UserID: authCode.UserID,
		Scope:  authCode.Scope,
		Nonce:  authCode.Nonce,
	}, nil
}

// GenerateTokenPair creates an access token and refresh token pair.
func (s *Service) GenerateTokenPair(ctx context.Context, userID int64, clientID uuid.UUID, scope, email, name string) (*TokenPair, error) {
	accessTokenExpiry := time.Duration(s.cfg.GetAccessTokenExpiry()) * time.Second
	refreshTokenExpiry := time.Duration(s.cfg.GetRefreshTokenExpiry()) * time.Second

	accessToken, err := s.jwtSigner.GenerateAccessToken(jwt.GenerateTokenParams{
		UserID:   userID,
		ClientID: clientID.String(),
		Scope:    scope,
		Email:    email,
		Name:     name,
		Expiry:   accessTokenExpiry,
	})
	if err != nil {
		s.log.Error("failed to generate access token",
			"error", err,
			"user_id", userID,
			"client_id", clientID,
		)
		return nil, err
	}

	refreshToken, err := s.repo.CreateRefreshToken(ctx, &CreateRefreshTokenInput{
		ClientID:  clientID,
		UserID:    userID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	})
	if err != nil {
		s.log.Error("failed to create refresh token",
			"error", err,
			"user_id", userID,
			"client_id", clientID,
		)
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token.String(),
		TokenType:    "Bearer",
		ExpiresIn:    s.cfg.GetAccessTokenExpiry(),
		Scope:        scope,
	}, nil
}

// ValidateAndRefreshToken validates and rotates a refresh token.
func (s *Service) ValidateAndRefreshToken(ctx context.Context, refreshTokenStr string, clientID uuid.UUID, email, name string) (*TokenPair, error) {
	tokenUUID, err := uuid.Parse(refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	refreshToken, err := s.repo.GetRefreshTokenByToken(ctx, tokenUUID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if refreshToken.ClientID != clientID {
		return nil, ErrClientMismatch
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	if refreshToken.ExchangeAt != nil {
		return nil, ErrRefreshTokenUsed
	}

	if err := s.repo.MarkRefreshTokenExchanged(ctx, tokenUUID); err != nil {
		s.log.Error("failed to mark refresh token exchanged",
			"error", err,
			"token", tokenUUID,
		)
		return nil, err
	}

	return s.GenerateTokenPair(ctx, refreshToken.UserID, clientID, refreshToken.Scope, email, name)
}

// CheckUserConsent checks if a user has granted consent for the requested scopes.
func (s *Service) CheckUserConsent(ctx context.Context, userID int64, clientID uuid.UUID, requestedScope string) (bool, error) {
	consent, err := s.repo.GetUserConsent(ctx, userID, clientID)
	if err != nil {
		if err == ErrUserConsentNotFound {
			return false, nil
		}
		return false, err
	}

	if consent.RevokedAt != nil {
		return false, nil
	}

	requestedScopes := strings.Fields(requestedScope)
	grantedScopes := strings.Fields(consent.Scope)

	for _, requested := range requestedScopes {
		if !slices.Contains(grantedScopes, requested) {
			return false, nil
		}
	}

	return true, nil
}

// SaveUserConsent saves or updates a user's consent for the given scopes.
func (s *Service) SaveUserConsent(ctx context.Context, userID int64, clientID uuid.UUID, scope string) error {
	_, err := s.repo.CreateOrUpdateUserConsent(ctx, &UserConsentInput{
		UserID:   userID,
		ClientID: clientID,
		Scope:    scope,
	})
	if err != nil {
		s.log.Error("failed to save user consent",
			"error", err,
			"user_id", userID,
			"client_id", clientID,
		)
		return err
	}
	return nil
}

// GetUserConsents retrieves all consents for a user.
func (s *Service) GetUserConsents(ctx context.Context, userID int64) ([]*UserConsentWithClient, error) {
	return s.repo.GetUserConsents(ctx, userID)
}

// RevokeUserConsent revokes a user's consent for a specific client.
func (s *Service) RevokeUserConsent(ctx context.Context, userID int64, clientID uuid.UUID) error {
	return s.repo.RevokeUserConsent(ctx, userID, clientID)
}

// AuthenticateClient verifies client credentials.
// Public clients: only validates client_id exists (no secret required)
// Confidential clients: validates client_id and client_secret
func (s *Service) AuthenticateClient(ctx context.Context, clientIDStr, clientSecret string) (*OAuthClientInfo, error) {
	clientUUID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return nil, ErrInvalidClientID
	}

	client, err := s.repo.GetOAuthClientByClientID(ctx, clientUUID)
	if err != nil {
		if err == ErrOAuthClientNotFound {
			return nil, ErrInvalidClientID
		}
		return nil, err
	}

	// Public clients: no secret required
	if !client.Confidential {
		s.log.Info("public client authenticated", "client_id", clientIDStr, "client_name", client.Name)
		return client, nil
	}

	// Confidential clients: secret required
	if clientSecret == "" {
		return nil, ErrClientSecretRequired
	}

	if client.SecretHash == nil {
		s.log.Error("confidential client missing secret hash", "client_id", clientIDStr)
		return nil, ErrInvalidClientSecret
	}

	valid, err := password.VerifyPassword(clientSecret, *client.SecretHash)
	if err != nil || !valid {
		return nil, ErrInvalidClientSecret
	}

	return client, nil
}

// GetOAuthClient retrieves an OAuth client by client_id.
func (s *Service) GetOAuthClient(ctx context.Context, clientIDStr string) (*OAuthClientInfo, error) {
	clientUUID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return nil, ErrInvalidClientID
	}

	client, err := s.repo.GetOAuthClientByClientID(ctx, clientUUID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ValidateRedirectURI checks if a redirect URI is registered for the client.
func (s *Service) ValidateRedirectURI(client *OAuthClientInfo, redirectURI string) bool {
	return slices.Contains(client.RedirectURIs, redirectURI)
}

// RevokeToken revokes a refresh token or access token.
func (s *Service) RevokeToken(ctx context.Context, token, tokenTypeHint string) error {
	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return nil
	}

	if tokenTypeHint == "access_token" {
		return nil
	}

	refreshToken, err := s.repo.GetRefreshTokenByToken(ctx, tokenUUID)
	if err != nil {
		return nil
	}

	if refreshToken.ExchangeAt != nil {
		return nil
	}

	if err := s.repo.MarkRefreshTokenExchanged(ctx, tokenUUID); err != nil {
		s.log.Error("failed to revoke refresh token", "error", err, "token", tokenUUID)
		return err
	}

	return nil
}

// IntrospectToken inspects a token and returns its metadata.
func (s *Service) IntrospectToken(ctx context.Context, token string, clientID uuid.UUID) (map[string]interface{}, error) {
	claims, err := s.jwtSigner.ValidateAccessToken(token)
	if err == nil {
		aud := ""
		if len(claims.Audience) > 0 {
			aud = claims.Audience[0]
		}

		if aud != clientID.String() {
			return map[string]interface{}{"active": false}, nil
		}

		return map[string]interface{}{
			"active":     true,
			"scope":      claims.Scope,
			"client_id":  aud,
			"username":   claims.Subject,
			"token_type": "Bearer",
			"exp":        claims.ExpiresAt.Unix(),
			"iat":        claims.IssuedAt.Unix(),
			"sub":        claims.Subject,
			"iss":        claims.Issuer,
		}, nil
	}

	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return map[string]interface{}{"active": false}, nil
	}

	refreshToken, err := s.repo.GetRefreshTokenByToken(ctx, tokenUUID)
	if err != nil {
		return map[string]interface{}{"active": false}, nil
	}

	if refreshToken.ClientID != clientID {
		return map[string]interface{}{"active": false}, nil
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return map[string]interface{}{"active": false}, nil
	}

	if refreshToken.ExchangeAt != nil {
		return map[string]interface{}{"active": false}, nil
	}

	return map[string]interface{}{
		"active":     true,
		"scope":      refreshToken.Scope,
		"client_id":  refreshToken.ClientID.String(),
		"token_type": "refresh_token",
		"exp":        refreshToken.ExpiresAt.Unix(),
	}, nil
}
