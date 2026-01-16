package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Issuer is the JWT issuer claim value.
const Issuer = "altalune-oauth"

// Signer handles JWT token generation and validation using RS256.
type Signer struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	kid        string
}

// NewSigner creates a new JWT signer from PEM key files.
func NewSigner(privateKeyPath, publicKeyPath, kid string) (*Signer, error) {
	privateKey, err := LoadPrivateKeyPEM(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load private key: %w", err)
	}

	publicKey, err := LoadPublicKeyPEM(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load public key: %w", err)
	}

	return &Signer{
		privateKey: privateKey,
		publicKey:  publicKey,
		kid:        kid,
	}, nil
}

// GetPublicKey returns the RSA public key.
func (s *Signer) GetPublicKey() *rsa.PublicKey {
	return s.publicKey
}

// GetKID returns the key ID used in JWT headers.
func (s *Signer) GetKID() string {
	return s.kid
}

// GenerateTokenParams holds parameters for access token generation.
type GenerateTokenParams struct {
	UserPublicID string        // User's public_id (nanoid) - used as JWT subject
	ClientID     string        // OAuth client ID (UUID string)
	Scope        string        // Space-separated OAuth scopes
	Email        string        // User email (if scope includes "email")
	Name         string        // User full name (if scope includes "profile")
	Perms        []string      // User permissions for stateless authorization
	Expiry       time.Duration // Token validity duration
}

// GenerateAccessToken creates a signed RS256 JWT access token.
func (s *Signer) GenerateAccessToken(params GenerateTokenParams) (string, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   params.UserPublicID,
			Audience:  jwt.ClaimStrings{params.ClientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(params.Expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		Scope: params.Scope,
		Email: params.Email,
		Name:  params.Name,
		Perms: params.Perms,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.kid

	return token.SignedString(s.privateKey)
}

// ValidateAccessToken parses and validates a JWT access token.
func (s *Signer) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GetTokenExpiry extracts the expiration time from a JWT without validation.
func (s *Signer) GetTokenExpiry(tokenString string) (time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &AccessTokenClaims{})
	if err != nil {
		return time.Time{}, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid token claims")
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, fmt.Errorf("token has no expiry")
	}

	return claims.ExpiresAt.Time, nil
}
