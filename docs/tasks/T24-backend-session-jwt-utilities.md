# Task T24: Backend Infrastructure - Session & JWT Utilities

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Backend Foundation
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** None (foundational task)

## Objective

Create foundational utilities for session management, JWT token generation/validation, and PKCE verification that will be used by the OAuth authorization server.

## Acceptance Criteria

- [ ] Session store using Gorilla sessions with HMAC-signed cookies
- [ ] Session middleware for authentication and CSRF protection
- [ ] JWT utilities for RS256 signing and verification
- [ ] JWKS JSON generation from RSA public key
- [ ] PKCE code_challenge validation (S256 and plain methods)
- [ ] All utilities have proper error handling

## Technical Requirements

### Session Management (`internal/session/`)

#### Store (`store.go`)
- Use `github.com/gorilla/sessions` for cookie-based sessions
- HMAC-signed cookies with session secret from config
- Session options:
  - HttpOnly: true
  - Secure: true (in production)
  - SameSite: Lax
  - MaxAge: configurable (default 24 hours)
  - Path: "/"

#### Session Data Model (`model.go`)
```go
type SessionData struct {
    UserID          int64     // Authenticated user ID (0 if not authenticated)
    AuthenticatedAt time.Time // When user authenticated
    OAuthState      string    // State parameter for OAuth CSRF
    OriginalURL     string    // URL to redirect after login
}
```

#### Middleware (`middleware.go`)
- `LoadSession`: Load session from cookie into context
- `AuthenticatedOnly`: Redirect to /login if not authenticated
- `ValidateCSRF`: Validate CSRF token on POST requests

### JWT Utilities (`internal/shared/jwt/`)

#### Key Management (`jwt.go`)
- Load RSA private key from PEM file (PKCS8 format)
- Load RSA public key from PEM file (PKIX format)
- Support 2048 and 4096 bit keys

#### Claims Structure (`claims.go`)
```go
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    Scope string `json:"scope,omitempty"`
    Email string `json:"email,omitempty"`
    Name  string `json:"name,omitempty"`
}
```

Standard claims to include:
- `iss`: "altalune-oauth"
- `sub`: user ID (as string)
- `aud`: client_id
- `exp`: expiration timestamp
- `iat`: issued at timestamp

#### Token Operations
- `GenerateAccessToken(claims) (string, error)`: Generate RS256 signed JWT
- `ValidateAccessToken(tokenString) (*AccessTokenClaims, error)`: Parse and validate JWT
- `GetTokenExpiry(tokenString) (time.Time, error)`: Extract expiry without full validation

#### JWKS Generation (`jwks.go`)
```go
type JWK struct {
    Kty string `json:"kty"` // "RSA"
    Use string `json:"use"` // "sig"
    Kid string `json:"kid"` // Key ID from config
    Alg string `json:"alg"` // "RS256"
    N   string `json:"n"`   // Base64url-encoded modulus
    E   string `json:"e"`   // Base64url-encoded exponent ("AQAB")
}

type JWKS struct {
    Keys []JWK `json:"keys"`
}
```
- `GenerateJWKS(publicKey, kid) (*JWKS, error)`: Generate JWKS from public key

### PKCE Utilities (`internal/shared/pkce/`)

#### PKCE Validation (`pkce.go`)
```go
// VerifyCodeChallenge verifies code_verifier matches code_challenge
func VerifyCodeChallenge(verifier, challenge, method string) bool

// GenerateCodeChallenge generates code_challenge from code_verifier (for testing)
func GenerateCodeChallenge(verifier, method string) (string, error)
```

Methods to support:
- `S256`: SHA256 hash of verifier, base64url-encoded
- `plain`: verifier equals challenge (not recommended but must support)

## Implementation Details

### Session Store Implementation

```go
// internal/session/store.go
package session

import (
    "github.com/gorilla/sessions"
)

const SessionName = "altalune_auth"

type Store struct {
    store   *sessions.CookieStore
    options *sessions.Options
}

func NewStore(secret string, secure bool, maxAge int) *Store {
    store := sessions.NewCookieStore([]byte(secret))
    store.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   maxAge,
        HttpOnly: true,
        Secure:   secure,
        SameSite: http.SameSiteLaxMode,
    }
    return &Store{store: store, options: store.Options}
}

func (s *Store) Get(r *http.Request) (*sessions.Session, error) {
    return s.store.Get(r, SessionName)
}

func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
    return s.store.Save(r, w, session)
}
```

### JWT Signer Implementation

```go
// internal/shared/jwt/jwt.go
package jwt

import (
    "crypto/rsa"
    "github.com/golang-jwt/jwt/v5"
)

type Signer struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
    kid        string
    issuer     string
}

func NewSigner(privateKeyPath, publicKeyPath, kid string) (*Signer, error) {
    // Load keys from PEM files
    privateKey, err := loadPrivateKey(privateKeyPath)
    publicKey, err := loadPublicKey(publicKeyPath)
    // ...
}

func (s *Signer) GenerateAccessToken(userID int64, clientID, scope string, expiry time.Duration) (string, error) {
    claims := AccessTokenClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    s.issuer,
            Subject:   strconv.FormatInt(userID, 10),
            Audience:  jwt.ClaimStrings{clientID},
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
        Scope: scope,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    token.Header["kid"] = s.kid

    return token.SignedString(s.privateKey)
}
```

### PKCE Verification

```go
// internal/shared/pkce/pkce.go
package pkce

import (
    "crypto/sha256"
    "encoding/base64"
)

func VerifyCodeChallenge(verifier, challenge, method string) bool {
    switch method {
    case "S256":
        hash := sha256.Sum256([]byte(verifier))
        computed := base64.RawURLEncoding.EncodeToString(hash[:])
        return computed == challenge
    case "plain":
        return verifier == challenge
    default:
        return false
    }
}
```

## Files to Create

- `internal/session/store.go` - Session store implementation
- `internal/session/middleware.go` - Session and CSRF middleware
- `internal/session/model.go` - Session data structures
- `internal/shared/jwt/jwt.go` - JWT signer and key loading
- `internal/shared/jwt/claims.go` - JWT claims structures
- `internal/shared/jwt/jwks.go` - JWKS generation
- `internal/shared/pkce/pkce.go` - PKCE verification utilities

## Files to Modify

- None (all new files)

## Testing Requirements

- Unit tests for JWT signing and verification
- Unit tests for PKCE challenge verification
- Unit tests for session data serialization/deserialization
- Test with various key sizes (2048, 4096 bits)

## Commands to Run

```bash
# Build and verify no compile errors
make build

# Run tests (when implemented)
go test ./internal/session/...
go test ./internal/shared/jwt/...
go test ./internal/shared/pkce/...
```

## Validation Checklist

- [ ] Session cookies are HttpOnly and Secure (in production)
- [ ] JWT signing uses RS256 algorithm
- [ ] JWT includes kid in header
- [ ] JWKS endpoint returns valid JSON
- [ ] PKCE S256 verification matches RFC 7636
- [ ] Error handling is consistent

## Definition of Done

- [ ] Session store implementation complete with cookie options
- [ ] Session middleware handles authentication and CSRF
- [ ] JWT signer generates valid RS256 tokens
- [ ] JWT validator properly verifies signatures and expiry
- [ ] JWKS generation produces valid JSON Web Key Set
- [ ] PKCE verification supports S256 and plain methods
- [ ] All utilities have proper error types
- [ ] Code follows established patterns and guidelines

## Dependencies

- `github.com/gorilla/sessions` - Session management
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- Standard library: crypto/rsa, crypto/sha256, encoding/base64

## Risk Factors

- **Low Risk**: Well-established libraries (Gorilla sessions, golang-jwt)
- **Medium Risk**: RSA key loading - must handle various PEM formats correctly

## Notes

- Session secret should be at least 32 bytes for security
- JWT private key should have restricted file permissions (600)
- JWKS kid should match the kid in JWT headers for key rotation support
- Consider adding key rotation support in future (multiple keys in JWKS)
