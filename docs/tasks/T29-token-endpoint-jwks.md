# Task T29: Token Endpoint & JWKS

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Backend Handlers
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T28 (Authorization & Consent Endpoints)

## Objective

Implement the /oauth/token endpoint for exchanging authorization codes for tokens and refreshing access tokens, plus the /.well-known/jwks.json endpoint for publishing the public key used to verify JWT signatures.

## Acceptance Criteria

- [ ] POST /oauth/token accepts authorization_code grant type
- [ ] POST /oauth/token accepts refresh_token grant type
- [ ] Client authentication via Basic Auth header
- [ ] Authorization code validation (client_id, expiry, exchange status)
- [ ] PKCE code_verifier validation when code_challenge was used
- [ ] Access token generation (RS256 signed JWT)
- [ ] Refresh token generation (UUID v4)
- [ ] Token response follows OAuth 2.0 spec
- [ ] GET /.well-known/jwks.json returns public key in JWK format
- [ ] Error responses follow OAuth 2.0 error codes

## Technical Requirements

### Token Handler (`handlers/token.go`)

```go
func Token(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 1. Parse form data
        if err := r.ParseForm(); err != nil {
            writeTokenError(w, "invalid_request", "Invalid form data", http.StatusBadRequest)
            return
        }

        // 2. Authenticate client via Basic Auth
        clientID, clientSecret, ok := r.BasicAuth()
        if !ok {
            w.Header().Set("WWW-Authenticate", `Basic realm="OAuth"`)
            writeTokenError(w, "invalid_client", "Client authentication required", http.StatusUnauthorized)
            return
        }

        client, err := c.OAuthAuthService().AuthenticateClient(ctx, clientID, clientSecret)
        if err != nil {
            writeTokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
            return
        }

        // 3. Handle based on grant_type
        grantType := r.FormValue("grant_type")

        switch grantType {
        case "authorization_code":
            handleAuthorizationCodeGrant(w, r, ctx, c, client)
        case "refresh_token":
            handleRefreshTokenGrant(w, r, ctx, c, client)
        default:
            writeTokenError(w, "unsupported_grant_type", "Grant type not supported", http.StatusBadRequest)
        }
    }
}
```

### Authorization Code Grant

```go
func handleAuthorizationCodeGrant(w http.ResponseWriter, r *http.Request, ctx context.Context, c *Container, client *OAuthClient) {
    // 1. Get parameters
    code := r.FormValue("code")
    redirectURI := r.FormValue("redirect_uri")
    codeVerifier := r.FormValue("code_verifier")

    // 2. Validate required parameters
    if code == "" {
        writeTokenError(w, "invalid_request", "Missing code parameter", http.StatusBadRequest)
        return
    }
    if redirectURI == "" {
        writeTokenError(w, "invalid_request", "Missing redirect_uri parameter", http.StatusBadRequest)
        return
    }

    // 3. Validate and exchange authorization code
    result, err := c.OAuthAuthService().ValidateAndExchangeCode(ctx, code, client.ClientID, redirectURI, &codeVerifier)
    if err != nil {
        switch err {
        case oauth_auth.ErrInvalidAuthorizationCode:
            writeTokenError(w, "invalid_grant", "Invalid authorization code", http.StatusBadRequest)
        case oauth_auth.ErrCodeExpired:
            writeTokenError(w, "invalid_grant", "Authorization code has expired", http.StatusBadRequest)
        case oauth_auth.ErrCodeAlreadyUsed:
            writeTokenError(w, "invalid_grant", "Authorization code has already been used", http.StatusBadRequest)
        case oauth_auth.ErrClientMismatch:
            writeTokenError(w, "invalid_grant", "Authorization code was not issued to this client", http.StatusBadRequest)
        case oauth_auth.ErrRedirectURIMismatch:
            writeTokenError(w, "invalid_grant", "Redirect URI does not match", http.StatusBadRequest)
        case oauth_auth.ErrMissingCodeVerifier:
            writeTokenError(w, "invalid_request", "PKCE code_verifier required", http.StatusBadRequest)
        case oauth_auth.ErrInvalidCodeVerifier:
            writeTokenError(w, "invalid_grant", "Invalid PKCE code_verifier", http.StatusBadRequest)
        default:
            writeTokenError(w, "server_error", "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // 4. Get user info for token claims
    user, err := c.UserService().GetByID(ctx, result.UserID)
    if err != nil {
        writeTokenError(w, "server_error", "Failed to get user info", http.StatusInternalServerError)
        return
    }

    // 5. Generate token pair
    tokenPair, err := c.OAuthAuthService().GenerateTokenPair(ctx, GenerateTokenPairInput{
        UserID:   result.UserID,
        ClientID: client.ClientID,
        Scope:    result.Scope,
        User:     user,
    })
    if err != nil {
        writeTokenError(w, "server_error", "Failed to generate tokens", http.StatusInternalServerError)
        return
    }

    // 6. Write response
    writeTokenResponse(w, tokenPair)
}
```

### Refresh Token Grant

```go
func handleRefreshTokenGrant(w http.ResponseWriter, r *http.Request, ctx context.Context, c *Container, client *OAuthClient) {
    // 1. Get refresh token
    refreshToken := r.FormValue("refresh_token")
    if refreshToken == "" {
        writeTokenError(w, "invalid_request", "Missing refresh_token parameter", http.StatusBadRequest)
        return
    }

    // 2. Validate and refresh
    tokenPair, err := c.OAuthAuthService().ValidateAndRefreshToken(ctx, refreshToken, client.ClientID)
    if err != nil {
        switch err {
        case oauth_auth.ErrInvalidRefreshToken:
            writeTokenError(w, "invalid_grant", "Invalid refresh token", http.StatusBadRequest)
        case oauth_auth.ErrRefreshTokenExpired:
            writeTokenError(w, "invalid_grant", "Refresh token has expired", http.StatusBadRequest)
        case oauth_auth.ErrRefreshTokenUsed:
            writeTokenError(w, "invalid_grant", "Refresh token has already been used", http.StatusBadRequest)
        case oauth_auth.ErrClientMismatch:
            writeTokenError(w, "invalid_grant", "Refresh token was not issued to this client", http.StatusBadRequest)
        default:
            writeTokenError(w, "server_error", "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // 3. Write response
    writeTokenResponse(w, tokenPair)
}
```

### Response Writers

```go
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
    Scope        string `json:"scope,omitempty"`
}

type TokenErrorResponse struct {
    Error            string `json:"error"`
    ErrorDescription string `json:"error_description,omitempty"`
}

func writeTokenResponse(w http.ResponseWriter, tokenPair *TokenPair) {
    response := TokenResponse{
        AccessToken:  tokenPair.AccessToken,
        TokenType:    "Bearer",
        ExpiresIn:    tokenPair.ExpiresIn,
        RefreshToken: tokenPair.RefreshToken,
        Scope:        tokenPair.Scope,
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "no-store")
    w.Header().Set("Pragma", "no-cache")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func writeTokenError(w http.ResponseWriter, errorCode, description string, statusCode int) {
    response := TokenErrorResponse{
        Error:            errorCode,
        ErrorDescription: description,
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "no-store")
    w.Header().Set("Pragma", "no-cache")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}
```

### JWKS Handler (`handlers/jwks.go`)

```go
func JWKS(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Generate JWKS from public key
        jwks, err := c.JWTSigner().GenerateJWKS()
        if err != nil {
            http.Error(w, "Failed to generate JWKS", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
        json.NewEncoder(w).Encode(jwks)
    }
}
```

### Service Layer Token Generation

```go
// internal/domain/oauth_auth/service.go

type GenerateTokenPairInput struct {
    UserID   int64
    ClientID uuid.UUID
    Scope    string
    User     *user.User
}

func (s *Service) GenerateTokenPair(ctx context.Context, input GenerateTokenPairInput) (*TokenPair, error) {
    // 1. Generate access token (JWT)
    accessTokenExpiry := time.Duration(s.config.GetAccessTokenExpiry()) * time.Second

    claims := jwt.AccessTokenClaims{
        Issuer:   "altalune-oauth",
        Subject:  strconv.FormatInt(input.UserID, 10),
        Audience: input.ClientID.String(),
        Scope:    input.Scope,
    }

    // Add optional claims based on scopes
    if strings.Contains(input.Scope, "email") {
        claims.Email = input.User.Email
    }
    if strings.Contains(input.Scope, "profile") {
        claims.Name = input.User.FullName()
    }

    accessToken, err := s.jwtSigner.GenerateAccessToken(claims, accessTokenExpiry)
    if err != nil {
        return nil, fmt.Errorf("generate access token: %w", err)
    }

    // 2. Generate refresh token (UUID)
    refreshToken := uuid.New()
    refreshTokenExpiry := time.Duration(s.config.GetRefreshTokenExpiry()) * time.Second

    // 3. Store refresh token
    _, err = s.repo.CreateRefreshToken(ctx, CreateRefreshTokenInput{
        Token:     refreshToken,
        ClientID:  input.ClientID,
        UserID:    input.UserID,
        Scope:     input.Scope,
        ExpiresAt: time.Now().Add(refreshTokenExpiry),
    })
    if err != nil {
        return nil, fmt.Errorf("store refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken.String(),
        TokenType:    "Bearer",
        ExpiresIn:    int(accessTokenExpiry.Seconds()),
        Scope:        input.Scope,
    }, nil
}

func (s *Service) ValidateAndRefreshToken(ctx context.Context, tokenStr string, clientID uuid.UUID) (*TokenPair, error) {
    // 1. Parse refresh token
    token, err := uuid.Parse(tokenStr)
    if err != nil {
        return nil, ErrInvalidRefreshToken
    }

    // 2. Get refresh token from database
    refreshToken, err := s.repo.GetRefreshTokenByToken(ctx, token)
    if err != nil {
        return nil, ErrInvalidRefreshToken
    }

    // 3. Validate client matches
    if refreshToken.ClientID != clientID {
        return nil, ErrClientMismatch
    }

    // 4. Check expiration
    if time.Now().After(refreshToken.ExpiresAt) {
        return nil, ErrRefreshTokenExpired
    }

    // 5. Check if already used
    if refreshToken.ExchangeAt != nil {
        return nil, ErrRefreshTokenUsed
    }

    // 6. Mark as used (soft-delete)
    if err := s.repo.MarkRefreshTokenExchanged(ctx, token); err != nil {
        return nil, fmt.Errorf("mark token exchanged: %w", err)
    }

    // 7. Get user for new token claims
    user, err := s.userRepo.GetByID(ctx, refreshToken.UserID)
    if err != nil {
        return nil, fmt.Errorf("get user: %w", err)
    }

    // 8. Generate new token pair
    return s.GenerateTokenPair(ctx, GenerateTokenPairInput{
        UserID:   refreshToken.UserID,
        ClientID: clientID,
        Scope:    refreshToken.Scope,
        User:     user,
    })
}
```

## Files to Create

- `internal/authserver/handlers/token.go` - Token endpoint handler
- `internal/authserver/handlers/jwks.go` - JWKS endpoint handler

## Files to Modify

- `internal/authserver/routes.go` - Update route handlers
- `internal/domain/oauth_auth/service.go` - Token generation methods

## Testing Requirements

- Test authorization_code grant with valid code
- Test with expired code
- Test with already-used code
- Test with wrong client_id
- Test with wrong redirect_uri
- Test with PKCE (valid and invalid verifier)
- Test refresh_token grant
- Test with expired refresh token
- Test with already-used refresh token
- Test JWKS endpoint returns valid JSON

## Commands to Run

```bash
# Build application
make build

# Start auth server
./bin/app serve-auth -c config.yaml

# Test token endpoint (after getting authorization code)
curl -X POST http://localhost:3101/oauth/token \
  -u "client_id:client_secret" \
  -d "grant_type=authorization_code" \
  -d "code=your-auth-code" \
  -d "redirect_uri=http://localhost:8089/callback"

# Test JWKS endpoint
curl http://localhost:3101/.well-known/jwks.json

# Verify JWT with JWKS
# Use jwt.io or similar to verify token signature
```

## Validation Checklist

- [ ] Authorization code exchange returns tokens
- [ ] Refresh token exchange returns new tokens
- [ ] Old refresh token is invalidated after use
- [ ] Access token is valid RS256 JWT
- [ ] JWT contains correct claims (iss, sub, aud, exp, scope)
- [ ] JWKS endpoint returns valid JSON Web Key Set
- [ ] JWKS kid matches JWT header kid
- [ ] Error responses follow OAuth 2.0 format
- [ ] No tokens in logs

## Definition of Done

- [ ] Token endpoint handles authorization_code grant
- [ ] Token endpoint handles refresh_token grant
- [ ] Client authentication via Basic Auth
- [ ] All authorization code validations implemented
- [ ] PKCE verification working
- [ ] Access token is RS256 signed JWT
- [ ] Refresh token is stored and single-use
- [ ] JWKS endpoint returns public key
- [ ] Error responses use standard OAuth error codes
- [ ] Response headers include Cache-Control: no-store

## Dependencies

- T28: Authorization endpoint (generates codes to exchange)
- T24: JWT signer for access token generation
- T25: OAuth auth domain for token storage

## Risk Factors

- **Low Risk**: Standard OAuth 2.0 token endpoint
- **Medium Risk**: JWT signing - must use correct key and algorithm
- **High Risk**: Refresh token single-use - must mark as used atomically

## Notes

- Token endpoint MUST use POST method only
- Client credentials MUST be in Basic Auth header (not form body)
- Refresh tokens are single-use - always generate new one on refresh
- Access tokens should be short-lived (1 hour default)
- Refresh tokens can be longer-lived (30 days default)
- JWKS should be cacheable (Cache-Control header)
- Never log tokens (security sensitive)
