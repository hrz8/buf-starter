package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenClaims mirrors the JWT claims structure from internal/shared/jwt.
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Scope         string            `json:"scope,omitempty"`
	Email         string            `json:"email,omitempty"`
	Name          string            `json:"name,omitempty"`
	Perms         []string          `json:"perms"`
	Memberships   map[string]string `json:"memberships,omitempty"`
	EmailVerified bool              `json:"email_verified"`
}

// JWTValidator validates JWT tokens using JWKS.
type JWTValidator struct {
	fetcher   *JWKSFetcher
	issuer    string
	audiences []string // Optional audience validation
}

// NewJWTValidator creates a new JWT validator.
func NewJWTValidator(jwksURL, issuer string, audiences []string, cacheTTL int, refreshLimit int) *JWTValidator {
	ttl := time.Duration(cacheTTL) * time.Second
	if cacheTTL == 0 {
		ttl = time.Hour // Default 1 hour
	}
	if refreshLimit == 0 {
		refreshLimit = 3 // Default 3 per minute
	}

	return &JWTValidator{
		fetcher:   NewJWKSFetcher(jwksURL, ttl, refreshLimit),
		issuer:    issuer,
		audiences: audiences,
	}
}

// Validate validates a JWT token and returns the claims.
func (v *JWTValidator) Validate(ctx context.Context, tokenString string) (*AccessTokenClaims, error) {
	// Parse token without validation to get kid
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, &AccessTokenClaims{})
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing kid in token header")
	}

	// Get public key from JWKS
	publicKey, err := v.fetcher.GetPublicKey(ctx, kid)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Parse and validate token with signature verification
	claims := &AccessTokenClaims{}
	token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		// Try refreshing JWKS on validation failure (key rotation)
		if strings.Contains(err.Error(), "signature") {
			if refreshErr := v.fetcher.ForceRefresh(ctx); refreshErr == nil {
				publicKey, _ = v.fetcher.GetPublicKey(ctx, kid)
				token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
					return publicKey, nil
				})
			}
		}

		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				return nil, fmt.Errorf("token has expired")
			}
			return nil, fmt.Errorf("invalid token signature")
		}
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Validate issuer
	if v.issuer != "" && claims.Issuer != v.issuer {
		return nil, fmt.Errorf("invalid token issuer")
	}

	// Validate audience (optional)
	if len(v.audiences) > 0 {
		valid := false
		for _, aud := range v.audiences {
			for _, tokenAud := range claims.Audience {
				if aud == tokenAud {
					valid = true
					break
				}
			}
			if valid {
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("invalid token audience")
		}
	}

	return claims, nil
}

// RefreshJWKS forces a refresh of the JWKS cache.
func (v *JWTValidator) RefreshJWKS(ctx context.Context) error {
	return v.fetcher.ForceRefresh(ctx)
}
