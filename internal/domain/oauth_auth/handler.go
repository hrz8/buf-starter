package oauth_auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/authserver/views"
	oauth_provider_domain "github.com/hrz8/altalune/internal/domain/oauth_provider"
	user_domain "github.com/hrz8/altalune/internal/domain/user"
	"github.com/hrz8/altalune/internal/session"
	"github.com/hrz8/altalune/internal/shared/jwt"
	"github.com/hrz8/altalune/internal/shared/oauthprovider"
)

type Handler struct {
	svc               *Service
	jwtSigner         *jwt.Signer
	sessionStore      *session.Store
	oauthProviderRepo oauth_provider_domain.Repository
	userRepo          user_domain.Repository
	log               altalune.Logger
}

func NewHandler(
	svc *Service,
	jwtSigner *jwt.Signer,
	sessionStore *session.Store,
	oauthProviderRepo oauth_provider_domain.Repository,
	userRepo user_domain.Repository,
	log altalune.Logger,
) *Handler {
	return &Handler{
		svc:               svc,
		jwtSigner:         jwtSigner,
		sessionStore:      sessionStore,
		oauthProviderRepo: oauthProviderRepo,
		userRepo:          userRepo,
		log:               log,
	}
}

func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if h.sessionStore.IsAuthenticated(r) {
		http.Redirect(w, r, "/oauth/authorize", http.StatusFound)
		return
	}

	errorMsg := r.URL.Query().Get("error")

	data := views.LoginPageData{
		BaseData: views.BaseData{
			Title: "Sign In",
		},
		Providers:    views.GetProviders(),
		ErrorMessage: errorMsg,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "login.html", data); err != nil {
		h.log.Error("failed to render login page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) HandleLoginProvider(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")

	provider, err := h.oauthProviderRepo.GetByProviderType(r.Context(), oauth_provider_domain.ProviderType(providerName))
	if err != nil {
		h.log.Error("failed to get provider", "provider", providerName, "error", err)
		http.Redirect(w, r, "/login?error=invalid_provider", http.StatusFound)
		return
	}

	if !provider.Enabled {
		http.Redirect(w, r, "/login?error=provider_disabled", http.StatusFound)
		return
	}

	state := generateSecureRandomString(32)

	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData == nil {
		sessionData = &session.Data{}
	}

	sessionData.OAuthState = state
	if nextURL := r.URL.Query().Get("next"); nextURL != "" {
		sessionData.OriginalURL = nextURL
	}

	if err := h.sessionStore.SetData(r, w, sessionData); err != nil {
		h.log.Error("failed to save session", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	clientSecret, err := h.oauthProviderRepo.RevealClientSecret(r.Context(), provider.ID)
	if err != nil {
		h.log.Error("failed to reveal client secret", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var client oauthprovider.Client
	switch oauth_provider_domain.ProviderType(providerName) {
	case oauth_provider_domain.ProviderTypeGoogle:
		client = oauthprovider.NewGoogleClient(provider.ClientID, clientSecret, provider.RedirectURL)
	case oauth_provider_domain.ProviderTypeGithub:
		client = oauthprovider.NewGitHubClient(provider.ClientID, clientSecret, provider.RedirectURL)
	default:
		http.Redirect(w, r, "/login?error=unsupported_provider", http.StatusFound)
		return
	}

	authURL := client.GetAuthorizationURL(state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil {
		http.Redirect(w, r, "/login?error=session_error", http.StatusFound)
		return
	}

	state := r.URL.Query().Get("state")
	if state != sessionData.OAuthState {
		http.Redirect(w, r, "/login?error=invalid_state", http.StatusFound)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		errorMsg := r.URL.Query().Get("error")
		http.Redirect(w, r, "/login?error="+errorMsg, http.StatusFound)
		return
	}

	providerName := r.URL.Query().Get("provider")
	if providerName == "" {
		providerName = "google"
	}

	provider, err := h.oauthProviderRepo.GetByProviderType(r.Context(), oauth_provider_domain.ProviderType(providerName))
	if err != nil {
		h.log.Error("failed to get provider", "error", err)
		http.Redirect(w, r, "/login?error=provider_error", http.StatusFound)
		return
	}

	clientSecret, err := h.oauthProviderRepo.RevealClientSecret(r.Context(), provider.ID)
	if err != nil {
		h.log.Error("failed to reveal client secret", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var client oauthprovider.Client
	switch provider.ProviderType {
	case oauth_provider_domain.ProviderTypeGoogle:
		client = oauthprovider.NewGoogleClient(provider.ClientID, clientSecret, provider.RedirectURL)
	case oauth_provider_domain.ProviderTypeGithub:
		client = oauthprovider.NewGitHubClient(provider.ClientID, clientSecret, provider.RedirectURL)
	default:
		http.Redirect(w, r, "/login?error=unsupported_provider", http.StatusFound)
		return
	}

	userInfo, err := client.ExchangeCodeForUserInfo(r.Context(), code)
	if err != nil {
		h.log.Error("failed to exchange code", "error", err)
		http.Redirect(w, r, "/login?error=exchange_failed", http.StatusFound)
		return
	}

	existingIdentity, err := h.userRepo.GetUserIdentityByProvider(r.Context(), string(provider.ProviderType), userInfo.ID)
	if err != nil && err != user_domain.ErrUserNotFound {
		h.log.Error("failed to check existing identity", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var userID int64
	if existingIdentity != nil {
		userID = existingIdentity.UserID
		if err := h.userRepo.UpdateUserIdentityLastLogin(r.Context(), userID, string(provider.ProviderType)); err != nil {
			h.log.Error("failed to update last login", "error", err)
		}
	} else {
		user, err := h.userRepo.Create(r.Context(), &user_domain.CreateUserInput{
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
		})
		if err != nil {
			h.log.Error("failed to create user", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		userID = user.ID

		// Extract OAuth client info from the original authorization URL for historical tracking
		var oauthClientIDStr *string
		var originClientName *string
		if sessionData.OriginalURL != "" {
			if params, err := parseAuthorizationParamsFromURL(sessionData.OriginalURL); err == nil && params.ClientID != uuid.Nil {
				clientIDStr := params.ClientID.String()
				oauthClientIDStr = &clientIDStr

				// Look up client to get name for historical snapshot
				if client, err := h.svc.GetOAuthClient(r.Context(), clientIDStr); err == nil {
					originClientName = &client.Name
				}
			}
		}

		if err := h.userRepo.CreateUserIdentity(r.Context(), &user_domain.CreateUserIdentityInput{
			UserID:                userID,
			Provider:              string(provider.ProviderType),
			ProviderUserID:        userInfo.ID,
			Email:                 userInfo.Email,
			FirstName:             userInfo.FirstName,
			LastName:              userInfo.LastName,
			OAuthClientID:         oauthClientIDStr,
			OriginOAuthClientName: originClientName,
		}); err != nil {
			h.log.Error("failed to create user identity", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := h.userRepo.AddProjectMember(r.Context(), 1, userID, "user"); err != nil {
			h.log.Error("failed to add project member", "error", err)
		}
	}

	sessionData.UserID = userID
	sessionData.AuthenticatedAt = time.Now()
	if err := h.sessionStore.SetData(r, w, sessionData); err != nil {
		h.log.Error("failed to save session", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	redirectURL := sessionData.OriginalURL
	if redirectURL == "" {
		redirectURL = "/oauth/authorize"
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionStore.Clear(r, w); err != nil {
		h.log.Error("failed to clear session", "error", err)
	}

	data := views.BaseData{
		Title: "Logged Out",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "logout.html", data); err != nil {
		h.log.Error("failed to render logout page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData.UserID == 0 {
		sessionData = &session.Data{OriginalURL: r.URL.String()}
		h.sessionStore.SetData(r, w, sessionData)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	params, err := parseAuthorizationParams(r)
	if err != nil {
		h.renderAuthError(w, r, "", "", err)
		return
	}

	if params.ResponseType != "code" {
		h.renderAuthError(w, r, params.RedirectURI, params.State, ErrUnsupportedResponseType)
		return
	}

	client, err := h.svc.GetOAuthClient(r.Context(), params.ClientID.String())
	if err != nil {
		h.renderError(w, "invalid_client", "Unknown client_id")
		return
	}

	if !h.svc.ValidateRedirectURI(client, params.RedirectURI) {
		h.renderError(w, "invalid_redirect_uri", "Redirect URI does not match registered URIs")
		return
	}

	if client.PKCERequired {
		if params.CodeChallenge == nil || *params.CodeChallenge == "" {
			h.renderAuthError(w, r, params.RedirectURI, params.State, ErrMissingCodeChallenge)
			return
		}
		if params.CodeChallengeMethod != nil && *params.CodeChallengeMethod != "S256" && *params.CodeChallengeMethod != "plain" {
			h.renderAuthError(w, r, params.RedirectURI, params.State, ErrInvalidCodeChallengeMethod)
			return
		}
	}

	hasConsent, err := h.svc.CheckUserConsent(r.Context(), sessionData.UserID, params.ClientID, params.Scope)
	if err != nil {
		h.renderAuthError(w, r, params.RedirectURI, params.State, ErrServerError)
		return
	}

	if hasConsent {
		code, err := h.svc.GenerateAuthorizationCode(r.Context(), &GenerateAuthCodeInput{
			ClientID:            params.ClientID,
			UserID:              sessionData.UserID,
			RedirectURI:         params.RedirectURI,
			Scope:               params.Scope,
			Nonce:               params.Nonce,
			CodeChallenge:       params.CodeChallenge,
			CodeChallengeMethod: params.CodeChallengeMethod,
		})
		if err != nil {
			h.renderAuthError(w, r, params.RedirectURI, params.State, ErrServerError)
			return
		}

		redirectWithCode(w, r, params.RedirectURI, code.Code.String(), params.State)
		return
	}

	csrfToken := generateCSRFToken()
	sessionData.CSRFToken = csrfToken
	if err := h.sessionStore.SetData(r, w, sessionData); err != nil {
		h.log.Error("failed to save session", "error", err)
	}

	h.renderConsentPage(w, client, params, csrfToken)
}

func (h *Handler) HandleAuthorizeProcess(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	csrfToken := r.FormValue("csrf_token")
	if csrfToken == "" || csrfToken != sessionData.CSRFToken {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	params := AuthorizationParams{
		RedirectURI:         r.FormValue("redirect_uri"),
		Scope:               r.FormValue("scope"),
		State:               r.FormValue("state"),
		Nonce:               stringPtr(r.FormValue("nonce")),
		CodeChallenge:       stringPtr(r.FormValue("code_challenge")),
		CodeChallengeMethod: stringPtr(r.FormValue("code_challenge_method")),
	}

	clientIDStr := r.FormValue("client_id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}
	params.ClientID = clientID

	decision := r.FormValue("decision")

	if decision == "deny" {
		redirectWithError(w, r, params.RedirectURI, "access_denied", "User denied the request", params.State)
		return
	}

	code, err := h.svc.GenerateAuthorizationCode(r.Context(), &GenerateAuthCodeInput{
		ClientID:            params.ClientID,
		UserID:              sessionData.UserID,
		RedirectURI:         params.RedirectURI,
		Scope:               params.Scope,
		Nonce:               params.Nonce,
		CodeChallenge:       params.CodeChallenge,
		CodeChallengeMethod: params.CodeChallengeMethod,
	})
	if err != nil {
		h.renderAuthError(w, r, params.RedirectURI, params.State, ErrServerError)
		return
	}

	if err := h.svc.SaveUserConsent(r.Context(), sessionData.UserID, params.ClientID, params.Scope); err != nil {
		h.log.Warn("failed to save user consent", "error", err)
	}

	redirectWithCode(w, r, params.RedirectURI, code.Code.String(), params.State)
}

func (h *Handler) HandleToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeTokenError(w, "invalid_request", "Invalid form data", http.StatusBadRequest)
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="OAuth"`)
		writeTokenError(w, "invalid_client", "Client authentication required", http.StatusUnauthorized)
		return
	}

	client, err := h.svc.AuthenticateClient(r.Context(), clientID, clientSecret)
	if err != nil {
		writeTokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
		return
	}

	grantType := r.FormValue("grant_type")

	switch grantType {
	case "authorization_code":
		h.handleAuthorizationCodeGrant(w, r, client)
	case "refresh_token":
		h.handleRefreshTokenGrant(w, r, client)
	default:
		writeTokenError(w, "unsupported_grant_type", "Grant type not supported", http.StatusBadRequest)
	}
}

func (h *Handler) handleAuthorizationCodeGrant(w http.ResponseWriter, r *http.Request, client *OAuthClientInfo) {
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	codeVerifier := r.FormValue("code_verifier")

	if code == "" {
		writeTokenError(w, "invalid_request", "Missing code parameter", http.StatusBadRequest)
		return
	}
	if redirectURI == "" {
		writeTokenError(w, "invalid_request", "Missing redirect_uri parameter", http.StatusBadRequest)
		return
	}

	var codeVerifierPtr *string
	if codeVerifier != "" {
		codeVerifierPtr = &codeVerifier
	}

	result, err := h.svc.ValidateAndExchangeCode(r.Context(), code, client.ClientID, redirectURI, codeVerifierPtr)
	if err != nil {
		switch err {
		case ErrInvalidAuthorizationCode:
			writeTokenError(w, "invalid_grant", "Invalid authorization code", http.StatusBadRequest)
		case ErrCodeExpired:
			writeTokenError(w, "invalid_grant", "Authorization code has expired", http.StatusBadRequest)
		case ErrCodeAlreadyUsed:
			writeTokenError(w, "invalid_grant", "Authorization code has already been used", http.StatusBadRequest)
		case ErrClientMismatch:
			writeTokenError(w, "invalid_grant", "Authorization code was not issued to this client", http.StatusBadRequest)
		case ErrRedirectURIMismatch:
			writeTokenError(w, "invalid_grant", "Redirect URI does not match", http.StatusBadRequest)
		case ErrMissingCodeVerifier:
			writeTokenError(w, "invalid_request", "PKCE code_verifier required", http.StatusBadRequest)
		case ErrInvalidCodeVerifier:
			writeTokenError(w, "invalid_grant", "Invalid PKCE code_verifier", http.StatusBadRequest)
		default:
			h.log.Error("token exchange error", "error", err)
			writeTokenError(w, "server_error", "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), fmt.Sprintf("%d", result.UserID))
	if err != nil {
		h.log.Error("failed to get user", "error", err, "user_id", result.UserID)
		writeTokenError(w, "server_error", "Failed to get user info", http.StatusInternalServerError)
		return
	}

	email := ""
	name := ""
	if strings.Contains(result.Scope, "email") {
		email = user.Email
	}
	if strings.Contains(result.Scope, "profile") {
		name = user.FirstName + " " + user.LastName
	}

	tokenPair, err := h.svc.GenerateTokenPair(r.Context(), result.UserID, client.ClientID, result.Scope, email, name)
	if err != nil {
		h.log.Error("failed to generate tokens", "error", err)
		writeTokenError(w, "server_error", "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	writeTokenResponse(w, tokenPair)
}

func (h *Handler) handleRefreshTokenGrant(w http.ResponseWriter, r *http.Request, client *OAuthClientInfo) {
	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		writeTokenError(w, "invalid_request", "Missing refresh_token parameter", http.StatusBadRequest)
		return
	}

	tokenPair, err := h.svc.ValidateAndRefreshToken(r.Context(), refreshToken, client.ClientID, "", "")
	if err != nil {
		switch err {
		case ErrInvalidRefreshToken:
			writeTokenError(w, "invalid_grant", "Invalid refresh token", http.StatusBadRequest)
		case ErrRefreshTokenExpired:
			writeTokenError(w, "invalid_grant", "Refresh token has expired", http.StatusBadRequest)
		case ErrRefreshTokenUsed:
			writeTokenError(w, "invalid_grant", "Refresh token has already been used", http.StatusBadRequest)
		case ErrClientMismatch:
			writeTokenError(w, "invalid_grant", "Refresh token was not issued to this client", http.StatusBadRequest)
		default:
			h.log.Error("refresh token error", "error", err)
			writeTokenError(w, "server_error", "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeTokenResponse(w, tokenPair)
}

func (h *Handler) HandleJWKS(w http.ResponseWriter, r *http.Request) {
	if h.jwtSigner == nil {
		http.Error(w, "JWKS not available", http.StatusInternalServerError)
		return
	}

	jwks := jwt.GenerateJWKS(h.jwtSigner.GetPublicKey(), h.jwtSigner.GetKID())
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(jwks)
}

func (h *Handler) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("WWW-Authenticate", `Bearer realm="OAuth"`)
		writeJSONError(w, "invalid_token", "Missing authorization header", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		w.Header().Set("WWW-Authenticate", `Bearer realm="OAuth"`)
		writeJSONError(w, "invalid_token", "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	accessToken := parts[1]
	claims, err := h.jwtSigner.ValidateAccessToken(accessToken)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Bearer realm="OAuth", error="invalid_token"`)
		writeJSONError(w, "invalid_token", "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), claims.Subject)
	if err != nil {
		h.log.Error("failed to get user", "error", err, "user_id", claims.Subject)
		writeJSONError(w, "server_error", "Failed to retrieve user info", http.StatusInternalServerError)
		return
	}

	scope := claims.Scope

	userInfo := map[string]interface{}{
		"sub": user.ID,
	}

	if strings.Contains(scope, "profile") {
		if user.FirstName != "" {
			userInfo["given_name"] = user.FirstName
		}
		if user.LastName != "" {
			userInfo["family_name"] = user.LastName
		}
		if user.FirstName != "" || user.LastName != "" {
			userInfo["name"] = strings.TrimSpace(user.FirstName + " " + user.LastName)
		}
	}

	if strings.Contains(scope, "email") {
		userInfo["email"] = user.Email
		userInfo["email_verified"] = true
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	json.NewEncoder(w).Encode(userInfo)
}

func (h *Handler) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSONError(w, "invalid_request", "Invalid form data", http.StatusBadRequest)
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="OAuth"`)
		writeJSONError(w, "invalid_client", "Client authentication required", http.StatusUnauthorized)
		return
	}

	_, err := h.svc.AuthenticateClient(r.Context(), clientID, clientSecret)
	if err != nil {
		writeJSONError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
		return
	}

	token := r.FormValue("token")
	if token == "" {
		writeJSONError(w, "invalid_request", "Missing token parameter", http.StatusBadRequest)
		return
	}

	tokenTypeHint := r.FormValue("token_type_hint")

	if err := h.svc.RevokeToken(r.Context(), token, tokenTypeHint); err != nil {
		h.log.Error("failed to revoke token", "error", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleIntrospect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSONError(w, "invalid_request", "Invalid form data", http.StatusBadRequest)
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="OAuth"`)
		writeJSONError(w, "invalid_client", "Client authentication required", http.StatusUnauthorized)
		return
	}

	client, err := h.svc.AuthenticateClient(r.Context(), clientID, clientSecret)
	if err != nil {
		writeJSONError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
		return
	}

	token := r.FormValue("token")
	if token == "" {
		writeJSONError(w, "invalid_request", "Missing token parameter", http.StatusBadRequest)
		return
	}

	introspection, err := h.svc.IntrospectToken(r.Context(), token, client.ClientID)
	if err != nil {
		h.log.Error("introspection error", "error", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"active": false})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(introspection)
}

func (h *Handler) HandleOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	config := map[string]interface{}{
		"issuer":                 baseURL,
		"authorization_endpoint": baseURL + "/oauth/authorize",
		"token_endpoint":         baseURL + "/oauth/token",
		"userinfo_endpoint":      baseURL + "/oauth/userinfo",
		"jwks_uri":               baseURL + "/.well-known/jwks.json",
		"revocation_endpoint":    baseURL + "/oauth/revoke",
		"introspection_endpoint": baseURL + "/oauth/introspect",
		"scopes_supported": []string{
			"openid",
			"profile",
			"email",
			"offline_access",
		},
		"response_types_supported": []string{
			"code",
		},
		"grant_types_supported": []string{
			"authorization_code",
			"refresh_token",
		},
		"subject_types_supported": []string{
			"public",
		},
		"id_token_signing_alg_values_supported": []string{
			"RS256",
		},
		"token_endpoint_auth_methods_supported": []string{
			"client_secret_basic",
		},
		"code_challenge_methods_supported": []string{
			"S256",
			"plain",
		},
		"claims_supported": []string{
			"sub",
			"name",
			"given_name",
			"family_name",
			"email",
			"email_verified",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(config)
}

func generateSecureRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

type AuthorizationParams struct {
	ResponseType        string
	ClientID            uuid.UUID
	RedirectURI         string
	Scope               string
	State               string
	Nonce               *string
	CodeChallenge       *string
	CodeChallengeMethod *string
}

func parseAuthorizationParams(r *http.Request) (*AuthorizationParams, error) {
	params := &AuthorizationParams{
		ResponseType: r.URL.Query().Get("response_type"),
		RedirectURI:  r.URL.Query().Get("redirect_uri"),
		Scope:        r.URL.Query().Get("scope"),
		State:        r.URL.Query().Get("state"),
	}

	nonce := r.URL.Query().Get("nonce")
	if nonce != "" {
		params.Nonce = &nonce
	}

	codeChallenge := r.URL.Query().Get("code_challenge")
	if codeChallenge != "" {
		params.CodeChallenge = &codeChallenge
	}

	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")
	if codeChallengeMethod == "" && codeChallenge != "" {
		codeChallengeMethod = "S256"
	}
	if codeChallengeMethod != "" {
		params.CodeChallengeMethod = &codeChallengeMethod
	}

	clientIDStr := r.URL.Query().Get("client_id")
	if clientIDStr == "" {
		return nil, ErrMissingClientID
	}
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return nil, ErrInvalidClientID
	}
	params.ClientID = clientID

	if params.ResponseType == "" {
		return nil, ErrMissingResponseType
	}
	if params.RedirectURI == "" {
		return nil, ErrMissingRedirectURI
	}

	return params, nil
}

// parseAuthorizationParamsFromURL parses authorization params from a URL string (for session originalURL)
func parseAuthorizationParamsFromURL(urlStr string) (*AuthorizationParams, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	params := &AuthorizationParams{
		ResponseType: u.Query().Get("response_type"),
		RedirectURI:  u.Query().Get("redirect_uri"),
		Scope:        u.Query().Get("scope"),
		State:        u.Query().Get("state"),
	}

	nonce := u.Query().Get("nonce")
	if nonce != "" {
		params.Nonce = &nonce
	}

	codeChallenge := u.Query().Get("code_challenge")
	if codeChallenge != "" {
		params.CodeChallenge = &codeChallenge
	}

	codeChallengeMethod := u.Query().Get("code_challenge_method")
	if codeChallengeMethod == "" && codeChallenge != "" {
		codeChallengeMethod = "S256"
	}
	if codeChallengeMethod != "" {
		params.CodeChallengeMethod = &codeChallengeMethod
	}

	clientIDStr := u.Query().Get("client_id")
	if clientIDStr != "" {
		clientID, err := uuid.Parse(clientIDStr)
		if err == nil {
			params.ClientID = clientID
		}
	}

	return params, nil
}

func redirectWithCode(w http.ResponseWriter, r *http.Request, redirectURI, code, state string) {
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("code", code)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func redirectWithError(w http.ResponseWriter, r *http.Request, redirectURI, errorCode, errorDesc, state string) {
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("error", errorCode)
	q.Set("error_description", errorDesc)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func generateCSRFToken() string {
	return generateSecureRandomString(32)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (h *Handler) renderAuthError(w http.ResponseWriter, r *http.Request, redirectURI, state string, err error) {
	if redirectURI != "" {
		errorCode := "server_error"
		errorDesc := err.Error()

		switch err {
		case ErrUnsupportedResponseType:
			errorCode = "unsupported_response_type"
		case ErrMissingCodeChallenge:
			errorCode = "invalid_request"
			errorDesc = "code_challenge is required for this client"
		case ErrInvalidCodeChallengeMethod:
			errorCode = "invalid_request"
			errorDesc = "code_challenge_method must be S256 or plain"
		}

		redirectWithError(w, r, redirectURI, errorCode, errorDesc, state)
		return
	}

	h.renderError(w, "invalid_request", err.Error())
}

func (h *Handler) renderError(w http.ResponseWriter, errorCode, errorDesc string) {
	data := views.ErrorPageData{
		BaseData: views.BaseData{
			Title: "Error",
		},
		Error:            errorCode,
		ErrorDescription: errorDesc,
		ShowBackToLogin:  true,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := views.Render(w, "error.html", data); err != nil {
		h.log.Error("failed to render error page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) renderConsentPage(w http.ResponseWriter, client *OAuthClientInfo, params *AuthorizationParams, csrfToken string) {
	scopes := parseScopes(params.Scope)

	data := views.ConsentPageData{
		BaseData: views.BaseData{
			Title: "Authorize",
		},
		ClientName:          client.Name,
		Scopes:              scopes,
		CSRFToken:           csrfToken,
		ClientID:            params.ClientID.String(),
		RedirectURI:         params.RedirectURI,
		Scope:               params.Scope,
		State:               params.State,
		Nonce:               params.Nonce,
		CodeChallenge:       params.CodeChallenge,
		CodeChallengeMethod: params.CodeChallengeMethod,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "consent.html", data); err != nil {
		h.log.Error("failed to render consent page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

var scopeDescriptions = map[string]string{
	"openid":         "Verify your identity",
	"profile":        "Access your profile information (name)",
	"email":          "Access your email address",
	"offline_access": "Access your data while you're offline",
}

func parseScopes(scopeString string) []views.ScopeInfo {
	if scopeString == "" {
		return []views.ScopeInfo{{Name: "openid", Description: "Verify your identity"}}
	}

	scopes := []views.ScopeInfo{}
	for _, scope := range strings.Fields(scopeString) {
		description := scopeDescriptions[scope]
		if description == "" {
			description = scope
		}
		scopes = append(scopes, views.ScopeInfo{
			Name:        scope,
			Description: description,
		})
	}

	return scopes
}

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
		TokenType:    tokenPair.TokenType,
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

func writeJSONError(w http.ResponseWriter, errorCode, description string, statusCode int) {
	response := map[string]string{
		"error":             errorCode,
		"error_description": description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
