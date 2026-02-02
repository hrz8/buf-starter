package auth

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
)

// authInterceptor implements connect.Interceptor for JWT validation.
type authInterceptor struct {
	validator *JWTValidator
}

// NewAuthInterceptor creates a Connect-RPC interceptor for JWT validation.
// It extracts Bearer tokens from Authorization header or access_token cookie,
// validates the JWT, and injects AuthContext into the request context.
func NewAuthInterceptor(validator *JWTValidator) connect.Interceptor {
	return &authInterceptor{validator: validator}
}

// WrapUnary implements connect.Interceptor for unary RPC calls.
func (i *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		tokenString := extractToken(req.Header())

		// If no token, continue with unauthenticated context
		// Individual handlers will decide if auth is required
		if tokenString == "" {
			authCtx := &AuthContext{IsAuthenticated: false}
			ctx = WithAuthContext(ctx, authCtx)
			return next(ctx, req)
		}

		// Validate JWT
		claims, err := i.validator.Validate(ctx, tokenString)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		// Create AuthContext from claims
		authCtx := NewAuthContextFromClaims(claims)
		ctx = WithAuthContext(ctx, authCtx)

		return next(ctx, req)
	}
}

// WrapStreamingClient implements connect.Interceptor for client streaming.
// This is a pass-through for server-side interceptors.
func (i *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler implements connect.Interceptor for server streaming.
func (i *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		tokenString := extractToken(conn.RequestHeader())

		if tokenString == "" {
			authCtx := &AuthContext{IsAuthenticated: false}
			ctx = WithAuthContext(ctx, authCtx)
			return next(ctx, conn)
		}

		claims, err := i.validator.Validate(ctx, tokenString)
		if err != nil {
			return connect.NewError(connect.CodeUnauthenticated, err)
		}

		authCtx := NewAuthContextFromClaims(claims)
		ctx = WithAuthContext(ctx, authCtx)

		return next(ctx, conn)
	}
}

// extractToken extracts the JWT token from Authorization header or cookie.
// Priority: Authorization header > access_token cookie
func extractToken(headers interface {
	Get(key string) string
}) string {
	// Try Authorization header first
	authHeader := headers.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		// Invalid format, return empty to fail authentication
		return ""
	}

	// Try Cookie header (for httpOnly cookie-based auth)
	cookieHeader := headers.Get("Cookie")
	if cookieHeader != "" {
		// Parse cookies to find access_token
		cookies := strings.Split(cookieHeader, ";")
		for _, cookie := range cookies {
			cookie = strings.TrimSpace(cookie)
			if strings.HasPrefix(cookie, "access_token=") {
				return strings.TrimPrefix(cookie, "access_token=")
			}
		}
	}

	return ""
}

// RequireAuth is a helper that returns an error if not authenticated.
// Use this at the start of handlers that require authentication.
func RequireAuth(ctx context.Context) error {
	auth := FromContext(ctx)
	if !auth.IsAuthenticated {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("authentication required"))
	}
	return nil
}
