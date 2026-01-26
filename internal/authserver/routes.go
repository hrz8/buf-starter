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
		s.cfg,
		s.c.GetJWTSigner(),
		s.c.GetSessionStore(),
		s.c.GetOAuthProviderRepo(),
		s.c.GetUserRepo(),
		s.c.GetRoleRepo(),
		s.c.GetIAMMapperRepo(),
		s.c.GetOTPService(),
		s.c.GetEmailVerificationService(),
		s.log,
	)

	mux.HandleFunc("GET /healthz", s.handleHealthz)

	// ============================================================================
	// Standalone IDP Routes (login without OAuth client)
	// ============================================================================
	mux.HandleFunc("GET /{$}", oauthAuthHandler.HandleRoot)
	mux.HandleFunc("GET /login/email", oauthAuthHandler.HandleEmailLoginPage)
	mux.HandleFunc("POST /login/email", oauthAuthHandler.HandleEmailLoginSubmit)
	mux.HandleFunc("GET /login/otp", oauthAuthHandler.HandleOTPPage)
	mux.HandleFunc("POST /login/otp/verify", oauthAuthHandler.HandleOTPVerify)
	mux.HandleFunc("GET /verify-email", oauthAuthHandler.HandleVerifyEmail)
	mux.HandleFunc("POST /resend-verification", oauthAuthHandler.HandleResendVerification)
	mux.HandleFunc("GET /pending-activation", oauthAuthHandler.HandlePendingActivation)

	// ============================================================================
	// OAuth Client Routes (this app acts as OAuth client to Providers e.g., Google/GitHub/etc)
	// Used for end-user authentication via external OAuth providers
	// ============================================================================
	mux.HandleFunc("GET /login", oauthAuthHandler.HandleLoginPage)
	mux.HandleFunc("GET /login/{provider}", oauthAuthHandler.HandleLoginProvider)
	mux.HandleFunc("GET /auth/callback", oauthAuthHandler.HandleOAuthCallback)
	mux.HandleFunc("GET /profile", oauthAuthHandler.HandleProfile)
	mux.HandleFunc("GET /edit-profile", oauthAuthHandler.HandleEditProfile)
	mux.HandleFunc("POST /edit-profile", oauthAuthHandler.HandleUpdateProfile)
	mux.HandleFunc("POST /profile/consents/revoke", oauthAuthHandler.HandleRevokeConsent)
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
