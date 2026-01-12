package authserver

import (
	"encoding/json"
	"net/http"

	oauth_auth_domain "github.com/hrz8/altalune/internal/domain/oauth_auth"
)

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	oauthAuthHandler := oauth_auth_domain.NewHandler(
		s.c.GetOAuthAuthService(),
		s.c.GetJWTSigner(),
		s.c.GetSessionStore(),
		s.c.GetOAuthProviderRepo(),
		s.c.GetUserRepo(),
		s.log,
	)

	mux.HandleFunc("GET /healthz", s.handleHealthz)

	// ============================================================================
	// OAuth Client Routes (this app acts as OAuth client to Providers e.g., Google/GitHub/etc)
	// Used for end-user authentication via external OAuth providers
	// ============================================================================
	mux.HandleFunc("GET /login", oauthAuthHandler.HandleLoginPage)
	mux.HandleFunc("GET /login/{provider}", oauthAuthHandler.HandleLoginProvider)
	mux.HandleFunc("GET /auth/callback", oauthAuthHandler.HandleOAuthCallback)
	mux.HandleFunc("POST /logout", oauthAuthHandler.HandleLogout)

	// ============================================================================
	// OAuth Provider/Authorization Server Routes (this app acts as OAuth provider)
	// Used by third-party OAuth clients to authenticate users and get tokens
	// ============================================================================

	// Authorization endpoints - interactive with end-user
	mux.HandleFunc("GET /oauth/authorize", oauthAuthHandler.HandleAuthorize)
	mux.HandleFunc("POST /oauth/authorize", oauthAuthHandler.HandleAuthorizeProcess)

	// Token endpoint - machine-to-machine
	mux.HandleFunc("POST /oauth/token", oauthAuthHandler.HandleToken)

	// UserInfo endpoint - returns user claims based on access token
	mux.HandleFunc("GET /oauth/userinfo", oauthAuthHandler.HandleUserInfo)

	// Token management endpoints
	mux.HandleFunc("POST /oauth/revoke", oauthAuthHandler.HandleRevoke)
	mux.HandleFunc("POST /oauth/introspect", oauthAuthHandler.HandleIntrospect)

	// JWKS endpoint - public key for token verification
	mux.HandleFunc("GET /.well-known/jwks.json", oauthAuthHandler.HandleJWKS)

	// OpenID Connect Discovery endpoint
	mux.HandleFunc("GET /.well-known/openid-configuration", oauthAuthHandler.HandleOpenIDConfiguration)

	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
