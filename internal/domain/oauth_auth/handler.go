package oauth_auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/authserver/views"
	iam_mapper_domain "github.com/hrz8/altalune/internal/domain/iam_mapper"
	oauth_provider_domain "github.com/hrz8/altalune/internal/domain/oauth_provider"
	role_domain "github.com/hrz8/altalune/internal/domain/role"
	user_domain "github.com/hrz8/altalune/internal/domain/user"
	"github.com/hrz8/altalune/internal/session"
	"github.com/hrz8/altalune/internal/shared/jwt"
	"github.com/hrz8/altalune/internal/shared/oauthprovider"
)

type Handler struct {
	svc                 *Service
	cfg                 altalune.Config
	jwtSigner           *jwt.Signer
	sessionStore        *session.Store
	oauthProviderRepo   oauth_provider_domain.Repository
	userRepo            user_domain.Repository
	roleRepo            role_domain.Repository
	iamMapperRepo       iam_mapper_domain.Repository
	otpService          *OTPService
	verificationService *EmailVerificationService
	log                 altalune.Logger
}

func NewHandler(
	svc *Service,
	cfg altalune.Config,
	jwtSigner *jwt.Signer,
	sessionStore *session.Store,
	oauthProviderRepo oauth_provider_domain.Repository,
	userRepo user_domain.Repository,
	roleRepo role_domain.Repository,
	iamMapperRepo iam_mapper_domain.Repository,
	otpService *OTPService,
	verificationService *EmailVerificationService,
	log altalune.Logger,
) *Handler {
	return &Handler{
		svc:                 svc,
		cfg:                 cfg,
		jwtSigner:           jwtSigner,
		sessionStore:        sessionStore,
		oauthProviderRepo:   oauthProviderRepo,
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		iamMapperRepo:       iamMapperRepo,
		otpService:          otpService,
		verificationService: verificationService,
		log:                 log,
	}
}

// baseData creates a BaseData struct with branding information.
func (h *Handler) baseData(title string) views.BaseData {
	return views.BaseData{
		Title: title,
		Branding: views.BrandingData{
			Name: h.cfg.GetAuthServerBrandingName(),
		},
	}
}

func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if h.sessionStore.IsAuthenticated(r) {
		// If user is already authenticated, redirect based on user status
		sessionData, _ := h.sessionStore.GetData(r)
		if sessionData != nil && sessionData.UserID != 0 {
			user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
			if err == nil {
				// Redirect inactive users to pending activation
				if !user.IsActive {
					http.Redirect(w, r, "/pending-activation", http.StatusFound)
					return
				}
			}
		}
		// Check for original URL first (OAuth flow)
		if sessionData != nil && sessionData.OriginalURL != "" {
			http.Redirect(w, r, sessionData.OriginalURL, http.StatusFound)
			return
		}
		// Default to profile for standalone login
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	errorMsg := r.URL.Query().Get("error")

	// Check if this login is for an OAuth client (UX transparency)
	var clientName string
	if clientID := r.URL.Query().Get("client_id"); clientID != "" {
		if client, err := h.svc.GetOAuthClient(r.Context(), clientID); err == nil {
			clientName = client.Name
		}
	}

	data := views.LoginPageData{
		BaseData:     h.baseData("Sign In"),
		Providers:    views.GetProviders(),
		ErrorMessage: errorMsg,
		ClientName:   clientName,
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
	sessionData.OAuthProvider = providerName
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

	// Get provider from session (set during HandleLoginProvider)
	providerName := sessionData.OAuthProvider
	if providerName == "" {
		providerName = r.URL.Query().Get("provider")
		if providerName == "" {
			h.log.Error("no provider found in session or query")
			http.Redirect(w, r, "/login?error=missing_provider", http.StatusFound)
			return
		}
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

	// Step 1: Check if identity already exists for this provider + provider_user_id
	existingIdentity, err := h.userRepo.GetUserIdentityByProvider(r.Context(), string(provider.ProviderType), userInfo.ID)
	if err != nil && err != user_domain.ErrUserNotFound {
		h.log.Error("failed to check existing identity", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var userID int64

	if existingIdentity != nil {
		// Identity exists - just update last login and use the existing user
		userID = existingIdentity.UserID
		if err := h.userRepo.UpdateUserIdentityLastLogin(r.Context(), userID, string(provider.ProviderType)); err != nil {
			h.log.Error("failed to update last login", "error", err)
		}
		h.log.Info("user logged in via existing identity",
			"userID", userID,
			"provider", provider.ProviderType,
			"email", userInfo.Email,
		)
	} else {
		// No identity for this provider - check if user exists by email (identity linking)
		existingUserID, err := h.userRepo.GetInternalIDByEmail(r.Context(), userInfo.Email)
		if err != nil && err != user_domain.ErrUserNotFound {
			h.log.Error("failed to check existing user by email", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Extract OAuth client info for historical tracking
		var oauthClientIDStr *string
		var originClientName *string
		var clientID string
		if sessionData.OriginalURL != "" {
			if params, err := parseAuthorizationParamsFromURL(sessionData.OriginalURL); err == nil && params.ClientID != uuid.Nil {
				clientIDStr := params.ClientID.String()
				oauthClientIDStr = &clientIDStr
				clientID = clientIDStr

				// Look up client to get name for historical snapshot
				if client, err := h.svc.GetOAuthClient(r.Context(), clientIDStr); err == nil {
					originClientName = &client.Name
				}
			}
		}

		if existingUserID > 0 {
			// Step 2: User exists with same email - link new identity to existing user
			userID = existingUserID

			// Create new identity linked to existing user
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
				h.log.Error("failed to create linked user identity", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			h.log.Info("linked new OAuth provider to existing user",
				"userID", userID,
				"provider", provider.ProviderType,
				"email", userInfo.Email,
			)
		} else {
			// Step 3: No user exists - create new user and identity
			regCtx := DetermineRegistrationContext(clientID, h.cfg.GetDefaultOAuthClientID())
			projectRole := GetProjectRoleForContext(regCtx)
			autoActivate := h.cfg.IsAutoActivate()

			user, err := h.userRepo.Create(r.Context(), &user_domain.CreateUserInput{
				Email:     userInfo.Email,
				FirstName: userInfo.FirstName,
				LastName:  userInfo.LastName,
				AvatarURL: userInfo.AvatarURL,
				IsActive:  &autoActivate,
			})
			if err != nil {
				h.log.Error("failed to create user", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			userID = user.ID

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

			// Assign user to default project with context-appropriate role
			if err := h.userRepo.AddProjectMember(r.Context(), 1, userID, projectRole); err != nil {
				h.log.Error("failed to add project member", "error", err, "role", projectRole)
			}

			// Assign global 'user' role to new user
			if h.roleRepo != nil && h.iamMapperRepo != nil {
				userRoleID, err := h.roleRepo.GetInternalIDByName(r.Context(), "user")
				if err != nil {
					h.log.Warn("failed to get 'user' role for assignment", "error", err)
				} else {
					if err := h.iamMapperRepo.AssignUserRoles(r.Context(), userID, []int64{userRoleID}); err != nil {
						h.log.Warn("failed to assign global 'user' role", "error", err, "userID", userID)
					}
				}
			}

			// Send verification email if user is auto-activated
			if autoActivate && h.verificationService != nil {
				if err := h.verificationService.GenerateAndSendVerificationEmail(r.Context(), userID); err != nil {
					h.log.Warn("failed to send verification email", "error", err, "userID", userID)
				}
			}

			h.log.Info("created new user via OAuth",
				"userID", userID,
				"email", userInfo.Email,
				"regContext", regCtx,
				"projectRole", projectRole,
				"autoActivated", autoActivate,
			)
		}
	}

	sessionData.UserID = userID
	sessionData.AuthenticatedAt = time.Now()
	if err := h.sessionStore.SetData(r, w, sessionData); err != nil {
		h.log.Error("failed to save session", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check user activation status for standalone login redirect
	var redirectURL string
	if sessionData.OriginalURL != "" {
		// OAuth client flow - redirect to original URL
		redirectURL = sessionData.OriginalURL
	} else {
		// Standalone IDP login - check user status
		user, err := h.userRepo.GetByInternalID(r.Context(), userID)
		if err != nil {
			h.log.Error("failed to get user for redirect", "error", err)
			redirectURL = "/profile"
		} else if !user.IsActive {
			// Inactive users go to pending activation
			redirectURL = "/pending-activation"
		} else {
			redirectURL = "/profile"
		}
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionStore.Clear(r, w); err != nil {
		h.log.Error("failed to clear session", "error", err)
	}

	data := h.baseData("Logged Out")

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

		loginURL := "/login"
		if clientID := r.URL.Query().Get("client_id"); clientID != "" {
			loginURL = "/login?client_id=" + clientID
		}
		http.Redirect(w, r, loginURL, http.StatusFound)
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

	// Check if user is active before allowing authorization
	// This prevents inactive users from completing OAuth flow to client applications
	user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user for authorization check", "error", err)
		h.renderAuthError(w, r, params.RedirectURI, params.State, ErrServerError)
		return
	}
	if !user.IsActive {
		// Return OAuth error to client - user account is not activated
		redirectWithError(w, r, params.RedirectURI, "access_denied", "account_not_activated", params.State)
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

	// Handle prompt=consent: always show consent page regardless of existing consent
	if params.Prompt == "consent" {
		csrfToken := generateCSRFToken()
		sessionData.CSRFToken = csrfToken
		h.sessionStore.SetData(r, w, sessionData)
		h.renderConsentPage(w, client, params, csrfToken)
		return
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

	// Check if user is still active (could have been deactivated while on consent page)
	user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user for consent processing", "error", err)
		h.renderAuthError(w, r, params.RedirectURI, params.State, ErrServerError)
		return
	}
	if !user.IsActive {
		redirectWithError(w, r, params.RedirectURI, "access_denied", "account_not_activated", params.State)
		return
	}

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

	// Try Basic Auth first (confidential clients)
	clientID, clientSecret, hasBasicAuth := r.BasicAuth()

	// If no Basic Auth, try form body (public clients)
	if !hasBasicAuth {
		clientID = r.FormValue("client_id")
		clientSecret = ""
	}

	if clientID == "" {
		writeTokenError(w, "invalid_client", "client_id is required", http.StatusBadRequest)
		return
	}

	client, err := h.svc.AuthenticateClient(r.Context(), clientID, clientSecret)
	if err != nil {
		switch err {
		case ErrClientSecretRequired:
			w.Header().Set("WWW-Authenticate", `Basic realm="OAuth"`)
			writeTokenError(w, "invalid_client", "Client authentication required", http.StatusUnauthorized)
		case ErrInvalidClientID:
			writeTokenError(w, "invalid_client", "Unknown client", http.StatusUnauthorized)
		case ErrInvalidClientSecret:
			writeTokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
		default:
			h.log.Error("client authentication error", "error", err)
			writeTokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
		}
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

	// Public clients MUST provide code_verifier (PKCE)
	if !client.Confidential && codeVerifier == "" {
		writeTokenError(w, "invalid_request", "PKCE code_verifier required for public clients", http.StatusBadRequest)
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

	user, err := h.userRepo.GetByInternalID(r.Context(), result.UserID)
	if err != nil {
		h.log.Error("failed to get user", "error", err, "user_id", result.UserID)
		writeTokenError(w, "server_error", "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Prevent inactive users from exchanging authorization codes
	if !user.IsActive {
		writeTokenError(w, "invalid_grant", "User account is not active", http.StatusBadRequest)
		return
	}

	scopeClaims, err := h.svc.BuildUserInfoClaims(r.Context(), result.Scope, &ScopeUser{
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AvatarURL:     user.AvatarURL,
		EmailVerified: user.EmailVerified,
	})
	if err != nil {
		h.log.Error("failed to build scope claims", "error", err)
		writeTokenError(w, "server_error", "Failed to process scopes", http.StatusInternalServerError)
		return
	}

	email, _ := scopeClaims["email"].(string)
	name, _ := scopeClaims["name"].(string)

	tokenPair, err := h.svc.GenerateTokenPair(r.Context(), &GenerateTokenPairParams{
		UserID:        result.UserID,
		UserPublicID:  user.ID,
		ClientID:      client.ClientID,
		Scope:         result.Scope,
		Email:         email,
		Name:          name,
		EmailVerified: user.EmailVerified,
	})
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

	result, err := h.svc.ValidateRefreshToken(r.Context(), refreshToken, client.ClientID)
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

	user, err := h.userRepo.GetByInternalID(r.Context(), result.UserID)
	if err != nil {
		h.log.Error("failed to get user for refresh token", "error", err, "user_id", result.UserID)
		writeTokenError(w, "server_error", "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Prevent inactive users from refreshing tokens
	if !user.IsActive {
		writeTokenError(w, "invalid_grant", "User account is not active", http.StatusBadRequest)
		return
	}

	scopeUser := &ScopeUser{
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AvatarURL:     user.AvatarURL,
		EmailVerified: user.EmailVerified,
	}
	scopeClaims, err := h.svc.BuildUserInfoClaims(r.Context(), result.Scope, scopeUser)
	if err != nil {
		h.log.Error("failed to build scope claims", "error", err)
		writeTokenError(w, "server_error", "Failed to process scopes", http.StatusInternalServerError)
		return
	}

	email, _ := scopeClaims["email"].(string)
	name, _ := scopeClaims["name"].(string)

	tokenPair, err := h.svc.GenerateTokenPair(r.Context(), &GenerateTokenPairParams{
		UserID:        result.UserID,
		UserPublicID:  user.ID,
		ClientID:      client.ClientID,
		Scope:         result.Scope,
		Email:         email,
		Name:          name,
		EmailVerified: user.EmailVerified,
	})
	if err != nil {
		h.log.Error("failed to generate tokens", "error", err)
		writeTokenError(w, "server_error", "Failed to generate tokens", http.StatusInternalServerError)
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

	// Subject is now user's public_id (nanoid), not internal ID
	userPublicID := claims.Subject

	user, err := h.userRepo.GetByID(r.Context(), userPublicID)
	if err != nil {
		h.log.Error("failed to get user", "error", err, "public_id", userPublicID)
		writeJSONError(w, "server_error", "Failed to retrieve user info", http.StatusInternalServerError)
		return
	}

	scope := claims.Scope

	// Build userinfo response using scope handler registry
	// Sub is always included per OIDC spec
	userInfo := map[string]interface{}{
		"sub": user.ID, // Return public_id as sub
	}

	// Use scope handler registry to build claims based on requested scopes
	scopeUser := &ScopeUser{
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AvatarURL:     user.AvatarURL,
		EmailVerified: user.EmailVerified,
	}
	scopeClaims, err := h.svc.BuildUserInfoClaims(r.Context(), scope, scopeUser)
	if err != nil {
		h.log.Error("failed to build userinfo claims", "error", err)
		writeJSONError(w, "server_error", "Failed to build user info", http.StatusInternalServerError)
		return
	}

	// Merge scope-based claims into userInfo
	for k, v := range scopeClaims {
		userInfo[k] = v
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
	// Use configured issuer for consistency with JWT tokens
	// This must match the "iss" claim in access tokens for JWKS discovery
	issuer := h.cfg.GetJWTIssuer()

	config := map[string]interface{}{
		"issuer":                 issuer,
		"authorization_endpoint": issuer + "/oauth/authorize",
		"token_endpoint":         issuer + "/oauth/token",
		"userinfo_endpoint":      issuer + "/oauth/userinfo",
		"jwks_uri":               issuer + "/.well-known/jwks.json",
		"revocation_endpoint":    issuer + "/oauth/revoke",
		"introspection_endpoint": issuer + "/oauth/introspect",
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
		// Note: This server returns claims in access tokens (JWT) rather than
		// separate ID tokens. Access tokens are RS256-signed and can be validated
		// using the JWKS endpoint for stateless authorization.
		"token_endpoint_auth_methods_supported": []string{
			"client_secret_basic",
			"none", // For public clients using PKCE
		},
		"introspection_endpoint_auth_methods_supported": []string{
			"client_secret_basic",
		},
		"revocation_endpoint_auth_methods_supported": []string{
			"client_secret_basic",
		},
		"code_challenge_methods_supported": []string{
			"S256",
			"plain",
		},
		// Claims available via userinfo endpoint when corresponding scopes are requested
		"claims_supported": []string{
			"sub",            // Always included (user public_id)
			"name",           // profile scope
			"given_name",     // profile scope
			"family_name",    // profile scope
			"picture",        // profile scope
			"email",          // email scope
			"email_verified", // email scope
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
	Prompt              string
}

func parseAuthorizationParams(r *http.Request) (*AuthorizationParams, error) {
	params := &AuthorizationParams{
		ResponseType: r.URL.Query().Get("response_type"),
		RedirectURI:  r.URL.Query().Get("redirect_uri"),
		Scope:        r.URL.Query().Get("scope"),
		State:        r.URL.Query().Get("state"),
		Prompt:       r.URL.Query().Get("prompt"),
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
		BaseData:         h.baseData("Error"),
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
		BaseData:            h.baseData("Authorize"),
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

// HandleProfile displays the user's profile page with authorized applications.
func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Redirect inactive users to pending activation page
	if !user.IsActive {
		http.Redirect(w, r, "/pending-activation", http.StatusFound)
		return
	}

	userIdentities, err := h.userRepo.GetUserIdentities(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user identities", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	consents, err := h.svc.GetUserConsents(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user consents", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check for verification email status from query params
	verificationStatus := r.URL.Query().Get("verification")
	verificationEmailSent := verificationStatus == "sent"
	verificationEmailError := verificationStatus == "error"

	data := views.ProfileData{
		BaseData:                   h.baseData("Your Profile"),
		User:                       user,
		Identities:                 userIdentities,
		Consents:                   consents,
		ShowEmailVerificationAlert: !user.EmailVerified,
		UserEmail:                  user.Email,
		VerificationEmailSent:      verificationEmailSent,
		VerificationEmailError:     verificationEmailError,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "profile.html", data); err != nil {
		h.log.Error("failed to render profile page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleRevokeConsent revokes a user's consent for a specific OAuth client.
func (h *Handler) HandleRevokeConsent(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	clientIDStr := r.FormValue("client_id")
	if clientIDStr == "" {
		http.Error(w, "Missing client_id", http.StatusBadRequest)
		return
	}

	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	if err := h.svc.RevokeUserConsent(r.Context(), sessionData.UserID, clientID); err != nil {
		if err == ErrUserConsentNotFound {
			http.Error(w, "Consent not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to revoke consent", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusFound)
}

// HandleRoot redirects based on authentication state.
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if h.sessionStore.IsAuthenticated(r) {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// HandleEmailLoginPage shows the email input form for OTP login.
func (h *Handler) HandleEmailLoginPage(w http.ResponseWriter, r *http.Request) {
	// Redirect logged-in users based on their status
	if h.sessionStore.IsAuthenticated(r) {
		sessionData, _ := h.sessionStore.GetData(r)
		if sessionData != nil && sessionData.UserID != 0 {
			user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
			if err == nil {
				if !user.IsActive {
					http.Redirect(w, r, "/pending-activation", http.StatusFound)
					return
				}
				http.Redirect(w, r, "/profile", http.StatusFound)
				return
			}
		}
	}

	errorMsg := r.URL.Query().Get("error")

	data := views.EmailLoginPageData{
		BaseData: h.baseData("Login with Email"),
		Error:    errorMsg,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "email_input.html", data); err != nil {
		h.log.Error("failed to render email login page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleEmailLoginSubmit processes email submission and sends OTP.
func (h *Handler) HandleEmailLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login/email?error=invalid_request", http.StatusFound)
		return
	}

	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	if email == "" {
		http.Redirect(w, r, "/login/email?error=email_required", http.StatusFound)
		return
	}

	// Check if OTP service is available
	if h.otpService == nil {
		h.log.Error("OTP service not configured")
		http.Redirect(w, r, "/login/email?error=server_error", http.StatusFound)
		return
	}

	// Generate and send OTP
	err := h.otpService.GenerateAndSendOTP(r.Context(), email)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailNotRegistered):
			http.Redirect(w, r, "/login/email?error=email_not_registered", http.StatusFound)
		case errors.Is(err, ErrOTPRateLimited):
			http.Redirect(w, r, "/login/email?error=rate_limited", http.StatusFound)
		default:
			h.log.Error("failed to send OTP", "error", err)
			http.Redirect(w, r, "/login/email?error=server_error", http.StatusFound)
		}
		return
	}

	// Store email in session for OTP verification
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData == nil {
		sessionData = &session.Data{}
	}
	sessionData.PendingOTPEmail = email
	if err := h.sessionStore.SetData(r, w, sessionData); err != nil {
		h.log.Error("failed to save session", "error", err)
		http.Redirect(w, r, "/login/email?error=server_error", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/login/otp", http.StatusFound)
}

// HandleOTPPage shows the OTP input form.
func (h *Handler) HandleOTPPage(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData == nil || sessionData.PendingOTPEmail == "" {
		http.Redirect(w, r, "/login/email", http.StatusFound)
		return
	}

	// Mask email for display (j***n@example.com)
	maskedEmail := maskEmail(sessionData.PendingOTPEmail)
	errorMsg := r.URL.Query().Get("error")

	data := views.OTPPageData{
		BaseData:   h.baseData("Enter Code"),
		Email:      maskedEmail,
		Error:      errorMsg,
		ExpiryMins: 5,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "otp_input.html", data); err != nil {
		h.log.Error("failed to render OTP page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleOTPVerify validates the OTP and creates a session.
func (h *Handler) HandleOTPVerify(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login/otp?error=invalid_request", http.StatusFound)
		return
	}

	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData == nil || sessionData.PendingOTPEmail == "" {
		http.Redirect(w, r, "/login/email", http.StatusFound)
		return
	}

	email := sessionData.PendingOTPEmail
	otp := r.FormValue("otp")

	if otp == "" {
		http.Redirect(w, r, "/login/otp?error=invalid_otp", http.StatusFound)
		return
	}

	// Check if OTP service is available
	if h.otpService == nil {
		h.log.Error("OTP service not configured")
		http.Redirect(w, r, "/login/otp?error=invalid_request", http.StatusFound)
		return
	}

	// Validate OTP
	user, err := h.otpService.ValidateOTP(r.Context(), email, otp)
	if err != nil {
		h.log.Debug("invalid OTP attempt", "email", email, "error", err)
		http.Redirect(w, r, "/login/otp?error=invalid_otp", http.StatusFound)
		return
	}

	// Create session first (so pending-activation page can access user info)
	sessionData.UserID = user.ID
	sessionData.AuthenticatedAt = time.Now()
	sessionData.PendingOTPEmail = "" // Clear pending email
	if err := h.sessionStore.SetData(r, w, sessionData); err != nil {
		h.log.Error("failed to save session", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if user is active - redirect inactive users to pending activation
	if !user.IsActive {
		http.Redirect(w, r, "/pending-activation", http.StatusFound)
		return
	}

	// Redirect to original URL or profile
	redirectURL := sessionData.OriginalURL
	if redirectURL == "" {
		redirectURL = "/profile"
	}
	sessionData.OriginalURL = ""
	h.sessionStore.SetData(r, w, sessionData)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// HandleVerifyEmail handles email verification link clicks.
func (h *Handler) HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	data := views.VerifyEmailResultData{
		BaseData: h.baseData("Email Verification"),
	}

	if token == "" {
		data.Success = false
		data.Error = "missing_token"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.Render(w, "verify_email_result.html", data); err != nil {
			h.log.Error("failed to render verify email result page", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Check if verification service is available
	if h.verificationService == nil {
		h.log.Error("verification service not configured")
		data.Success = false
		data.Error = "invalid_token"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.Render(w, "verify_email_result.html", data)
		return
	}

	err := h.verificationService.VerifyEmail(r.Context(), token)
	if err != nil {
		h.log.Debug("email verification failed", "error", err)
		data.Success = false
		if errors.Is(err, ErrInvalidVerificationToken) {
			data.Error = "expired_or_used"
		} else {
			data.Error = "invalid_token"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.Render(w, "verify_email_result.html", data)
		return
	}

	data.Success = true
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "verify_email_result.html", data); err != nil {
		h.log.Error("failed to render verify email result page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleResendVerification resends verification email for authenticated user.
func (h *Handler) HandleResendVerification(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData == nil || sessionData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Check if verification service is available
	if h.verificationService == nil {
		h.log.Error("verification service not configured")
		http.Redirect(w, r, "/profile?verification=error", http.StatusFound)
		return
	}

	err = h.verificationService.ResendVerificationEmail(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to resend verification email", "error", err, "userID", sessionData.UserID)
		http.Redirect(w, r, "/profile?verification=error", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/profile?verification=sent", http.StatusFound)
}

// HandlePendingActivation shows the pending activation page.
func (h *Handler) HandlePendingActivation(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Redirect active users to their original destination or profile
	if user.IsActive {
		redirectURL := "/profile"
		if sessionData.OriginalURL != "" {
			redirectURL = sessionData.OriginalURL
			// Clear the original URL after using it
			sessionData.OriginalURL = ""
			h.sessionStore.SetData(r, w, sessionData)
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	data := views.PendingActivationData{
		BaseData:  h.baseData("Account Pending"),
		UserEmail: user.Email,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "pending_activation.html", data); err != nil {
		h.log.Error("failed to render pending activation page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// maskEmail masks an email for display (e.g., j***n@example.com).
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	local := parts[0]
	if len(local) <= 2 {
		return local[0:1] + "***@" + parts[1]
	}
	return local[0:1] + "***" + local[len(local)-1:] + "@" + parts[1]
}

// HandleEditProfile shows the edit profile form.
func (h *Handler) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Redirect inactive users to pending activation page
	if !user.IsActive {
		http.Redirect(w, r, "/pending-activation", http.StatusFound)
		return
	}

	data := views.EditProfileData{
		BaseData: h.baseData("Edit Profile"),
		User:     user,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "edit_profile.html", data); err != nil {
		h.log.Error("failed to render edit profile page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleUpdateProfile processes the edit profile form submission.
func (h *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData, err := h.sessionStore.GetData(r)
	if err != nil || sessionData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := h.userRepo.GetByInternalID(r.Context(), sessionData.UserID)
	if err != nil {
		h.log.Error("failed to get user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Redirect inactive users to pending activation page
	if !user.IsActive {
		http.Redirect(w, r, "/pending-activation", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))

	// Validate input lengths
	if len(firstName) > 100 || len(lastName) > 100 {
		data := views.EditProfileData{
			BaseData:     h.baseData("Edit Profile"),
			User:         user,
			ErrorMessage: "Name fields must be 100 characters or less",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.Render(w, "edit_profile.html", data)
		return
	}

	// Update the profile
	updatedUser, err := h.userRepo.UpdateProfileByInternalID(r.Context(), sessionData.UserID, firstName, lastName)
	if err != nil {
		h.log.Error("failed to update user profile", "error", err)
		data := views.EditProfileData{
			BaseData:     h.baseData("Edit Profile"),
			User:         user,
			ErrorMessage: "Failed to update profile. Please try again.",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.Render(w, "edit_profile.html", data)
		return
	}

	// Show success message
	data := views.EditProfileData{
		BaseData: h.baseData("Edit Profile"),
		User:     updatedUser,
		Success:  true,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "edit_profile.html", data); err != nil {
		h.log.Error("failed to render edit profile page", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
