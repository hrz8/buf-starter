package oauth_auth

import (
	"context"
	"maps"
	"strings"
)

// ScopeUser represents user data needed by scope handlers.
// This abstraction prevents coupling to specific user domain types.
type ScopeUser struct {
	Email         string
	FirstName     string
	LastName      string
	EmailVerified bool
}

// ScopeHandler processes a specific OAuth scope and returns claims to add to JWT.
type ScopeHandler interface {
	// Scope returns the scope name this handler processes (e.g., "email", "profile")
	Scope() string

	// Handle processes the scope and returns claims to merge into the JWT.
	// Returns nil claims if scope doesn't apply additional data.
	Handle(ctx context.Context, user *ScopeUser) (map[string]interface{}, error)
}

// ScopeHandlerRegistry manages scope handlers and processes requested scopes.
type ScopeHandlerRegistry struct {
	handlers map[string]ScopeHandler
}

// NewScopeHandlerRegistry creates a registry with default OIDC handlers.
func NewScopeHandlerRegistry() *ScopeHandlerRegistry {
	registry := &ScopeHandlerRegistry{
		handlers: make(map[string]ScopeHandler),
	}
	// Register default OIDC scope handlers
	registry.Register(&EmailScopeHandler{})
	registry.Register(&ProfileScopeHandler{})
	return registry
}

// Register adds a scope handler to the registry.
func (r *ScopeHandlerRegistry) Register(handler ScopeHandler) {
	r.handlers[handler.Scope()] = handler
}

// ProcessScopes executes handlers for requested scopes and returns merged claims.
func (r *ScopeHandlerRegistry) ProcessScopes(ctx context.Context, scopes string, user *ScopeUser) (map[string]interface{}, error) {
	claims := make(map[string]any)
	for scope := range strings.FieldsSeq(scopes) {
		if handler, ok := r.handlers[scope]; ok {
			scopeClaims, err := handler.Handle(ctx, user)
			if err != nil {
				return nil, err
			}
			maps.Copy(claims, scopeClaims)
		}
	}
	return claims, nil
}

// EmailScopeHandler handles the "email" OIDC scope.
type EmailScopeHandler struct{}

func (h *EmailScopeHandler) Scope() string { return "email" }

func (h *EmailScopeHandler) Handle(_ context.Context, user *ScopeUser) (map[string]interface{}, error) {
	if user.Email == "" {
		return nil, nil
	}
	return map[string]any{
		"email":          user.Email,
		"email_verified": user.EmailVerified,
	}, nil
}

// ProfileScopeHandler handles the "profile" OIDC scope.
type ProfileScopeHandler struct{}

func (h *ProfileScopeHandler) Scope() string { return "profile" }

func (h *ProfileScopeHandler) Handle(_ context.Context, user *ScopeUser) (map[string]interface{}, error) {
	claims := make(map[string]any)
	if user.FirstName != "" || user.LastName != "" {
		claims["name"] = strings.TrimSpace(user.FirstName + " " + user.LastName)
	}
	if user.FirstName != "" {
		claims["given_name"] = user.FirstName
	}
	if user.LastName != "" {
		claims["family_name"] = user.LastName
	}
	return claims, nil
}
