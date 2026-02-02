package auth

import (
	"context"
)

type contextKey string

const authContextKey contextKey = "auth_context"

// AuthContext holds the authenticated user's information from JWT.
type AuthContext struct {
	UserID          string            // JWT subject (user public_id)
	Email           string
	Name            string
	Permissions     []string
	Memberships     map[string]string // project_public_id -> role
	EmailVerified   bool
	IsAuthenticated bool
}

// FromContext extracts AuthContext from request context.
func FromContext(ctx context.Context) *AuthContext {
	auth, ok := ctx.Value(authContextKey).(*AuthContext)
	if !ok {
		return &AuthContext{IsAuthenticated: false}
	}
	return auth
}

// WithAuthContext adds AuthContext to request context.
func WithAuthContext(ctx context.Context, auth *AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey, auth)
}

// NewAuthContextFromClaims creates AuthContext from JWT claims.
func NewAuthContextFromClaims(claims *AccessTokenClaims) *AuthContext {
	return &AuthContext{
		UserID:          claims.Subject,
		Email:           claims.Email,
		Name:            claims.Name,
		Permissions:     claims.Perms,
		Memberships:     claims.Memberships,
		EmailVerified:   claims.EmailVerified,
		IsAuthenticated: true,
	}
}
