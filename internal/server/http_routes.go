package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/gen/altalune/v1/altalunev1connect"
	"github.com/hrz8/altalune/gen/greeter/v1/greeterv1connect"
	api_key_domain "github.com/hrz8/altalune/internal/domain/api_key"
	employee_domain "github.com/hrz8/altalune/internal/domain/employee"
	greeter_domain "github.com/hrz8/altalune/internal/domain/greeter"
	iam_mapper_domain "github.com/hrz8/altalune/internal/domain/iam_mapper"
	oauth_client_domain "github.com/hrz8/altalune/internal/domain/oauth_client"
	oauth_provider_domain "github.com/hrz8/altalune/internal/domain/oauth_provider"
	permission_domain "github.com/hrz8/altalune/internal/domain/permission"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
	role_domain "github.com/hrz8/altalune/internal/domain/role"
	user_domain "github.com/hrz8/altalune/internal/domain/user"
)

func (s *Server) setupRoutes() *http.ServeMux {
	connectrpcMux := http.NewServeMux()

	// Examples
	greeterHandler := greeter_domain.NewHandler(s.c.GetGreeterService())
	employeeHandler := employee_domain.NewHandler(s.c.GetEmployeeService())
	greeterPath, greeterConnectHandler := greeterv1connect.NewGreeterServiceHandler(greeterHandler)
	employeePath, employeeConnectHandler := altalunev1connect.NewEmployeeServiceHandler(employeeHandler)
	connectrpcMux.Handle(greeterPath, greeterConnectHandler)
	connectrpcMux.Handle(employeePath, employeeConnectHandler)

	// Domains
	projectHandler := project_domain.NewHandler(s.c.GetProjectService())
	projectPath, projectConnectHandler := altalunev1connect.NewProjectServiceHandler(projectHandler)
	connectrpcMux.Handle(projectPath, projectConnectHandler)

	apiKeyHandler := api_key_domain.NewHandler(s.c.GetApiKeyService())
	apiKeyPath, apiKeyConnectHandler := altalunev1connect.NewApiKeyServiceHandler(apiKeyHandler)
	connectrpcMux.Handle(apiKeyPath, apiKeyConnectHandler)

	// IAM Domains
	userHandler := user_domain.NewHandler(s.c.GetUserService())
	userPath, userConnectHandler := altalunev1connect.NewUserServiceHandler(userHandler)
	connectrpcMux.Handle(userPath, userConnectHandler)

	roleHandler := role_domain.NewHandler(s.c.GetRoleService())
	rolePath, roleConnectHandler := altalunev1connect.NewRoleServiceHandler(roleHandler)
	connectrpcMux.Handle(rolePath, roleConnectHandler)

	permissionHandler := permission_domain.NewHandler(s.c.GetPermissionService())
	permissionPath, permissionConnectHandler := altalunev1connect.NewPermissionServiceHandler(permissionHandler)
	connectrpcMux.Handle(permissionPath, permissionConnectHandler)

	iamMapperHandler := iam_mapper_domain.NewHandler(s.c.GetIAMMapperService())
	iamMapperPath, iamMapperConnectHandler := altalunev1connect.NewIAMMapperServiceHandler(iamMapperHandler)
	connectrpcMux.Handle(iamMapperPath, iamMapperConnectHandler)

	oauthProviderHandler := oauth_provider_domain.NewHandler(s.c.GetOAuthProviderService())
	oauthProviderPath, oauthProviderConnectHandler := altalunev1connect.NewOAuthProviderServiceHandler(oauthProviderHandler)
	connectrpcMux.Handle(oauthProviderPath, oauthProviderConnectHandler)

	oauthClientHandler := oauth_client_domain.NewHandler(s.c.GetOAuthClientService())
	oauthClientPath, oauthClientConnectHandler := altalunev1connect.NewOAuthClientServiceHandler(oauthClientHandler)
	connectrpcMux.Handle(oauthClientPath, oauthClientConnectHandler)

	// main server mux
	mux := http.NewServeMux()

	// OAuth BFF endpoints (keeps tokens secure on backend via httpOnly cookies)
	mux.HandleFunc("/oauth/exchange", s.handleAuthExchange)
	mux.HandleFunc("/oauth/logout", s.handleAuthLogout)
	mux.HandleFunc("/oauth/refresh", s.handleAuthRefresh)
	mux.HandleFunc("/oauth/me", s.handleAuthMe)

	mux.Handle("/api/", http.StripPrefix("/api", connectrpcMux))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		health := map[string]any{
			"status": "ok",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// serve frontend
	websiteFS, _ := fs.Sub(altalune.FrontendEmbeddedFiles, "frontend/.output/public")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Exclude API and OAuth endpoints from SPA serving
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/oauth/") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

		if p == "" {
			serveFileOr404(w, r, websiteFS, "index.html")
			return
		}

		if exists(websiteFS, p) {
			if isDir(websiteFS, p) {
				serveFileOr404(w, r, websiteFS, path.Join(p, "index.html"))
				return
			}
			serveFileOr404(w, r, websiteFS, p)
			return
		}

		serve404Page(w, r, websiteFS)
	})

	return mux
}

// AuthExchangeRequest is the request body for token exchange
type AuthExchangeRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
	RedirectURI  string `json:"redirect_uri"`
}

// TokenResponse is the response from the auth server token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthUserInfo represents user information from OIDC /userinfo endpoint
type AuthUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Name          string `json:"name,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
}

// AuthExchangeSuccessResponse is the success response from token exchange
type AuthExchangeSuccessResponse struct {
	User      AuthUserInfo `json:"user"`
	ExpiresIn int          `json:"expires_in"`
}

// AuthErrorResponse is the error response for auth endpoints
type AuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// handleAuthExchange proxies OAuth token exchange requests to the auth server.
// Sets httpOnly cookies for tokens and returns user info in response body.
//
// POST /oauth/exchange
// Request:  { "code": "...", "code_verifier": "...", "redirect_uri": "..." }
// Response: { "user": { "sub": "...", "email": "...", "name": "..." }, "expires_in": 3600 }
// Cookies:  access_token (httpOnly), refresh_token (httpOnly)
func (s *Server) handleAuthExchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeAuthError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST method is allowed")
		return
	}

	var req AuthExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeAuthError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.Code == "" {
		s.writeAuthError(w, http.StatusBadRequest, "invalid_request", "Missing required field: code")
		return
	}
	if req.CodeVerifier == "" {
		s.writeAuthError(w, http.StatusBadRequest, "invalid_request", "Missing required field: code_verifier")
		return
	}
	if req.RedirectURI == "" {
		s.writeAuthError(w, http.StatusBadRequest, "invalid_request", "Missing required field: redirect_uri")
		return
	}

	// Exchange code for tokens with auth server
	tokenResp, err := s.exchangeCodeForTokens(req.Code, req.CodeVerifier, req.RedirectURI)
	if err != nil {
		s.log.Error("Token exchange failed", "error", err)
		s.writeAuthError(w, http.StatusBadGateway, "server_error", "Token exchange failed")
		return
	}

	// Set httpOnly cookies for tokens
	s.setAuthCookies(w, tokenResp)

	// Fetch user info from /userinfo endpoint (proper OIDC pattern)
	userInfo, err := s.fetchUserInfo(tokenResp.AccessToken)
	if err != nil {
		s.log.Warn("Failed to fetch userinfo, falling back to JWT", "error", err)
		// Fallback to JWT extraction if /userinfo fails
		jwtUserInfo := extractUserInfoFromJWT(tokenResp.AccessToken)
		userInfo = &jwtUserInfo
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthExchangeSuccessResponse{
		User:      *userInfo,
		ExpiresIn: tokenResp.ExpiresIn,
	})
}

// handleAuthLogout clears auth cookies.
//
// POST /oauth/logout
func (s *Server) handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeAuthError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST method is allowed")
		return
	}

	s.clearAuthCookies(w)
	w.WriteHeader(http.StatusOK)
}

// handleAuthRefresh refreshes tokens using the refresh_token cookie.
//
// POST /oauth/refresh
// Response: { "user": { "sub": "...", "email": "...", "name": "..." }, "expires_in": 3600 }
func (s *Server) handleAuthRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeAuthError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST method is allowed")
		return
	}

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		s.writeAuthError(w, http.StatusUnauthorized, "unauthorized", "No refresh token")
		return
	}

	tokenResp, err := s.refreshTokens(refreshCookie.Value)
	if err != nil {
		s.log.Error("Token refresh failed", "error", err)
		s.clearAuthCookies(w)
		s.writeAuthError(w, http.StatusUnauthorized, "unauthorized", "Token refresh failed")
		return
	}

	s.setAuthCookies(w, tokenResp)

	// Fetch user info from /userinfo endpoint (proper OIDC pattern)
	userInfo, err := s.fetchUserInfo(tokenResp.AccessToken)
	if err != nil {
		s.log.Warn("Failed to fetch userinfo, falling back to JWT", "error", err)
		jwtUserInfo := extractUserInfoFromJWT(tokenResp.AccessToken)
		userInfo = &jwtUserInfo
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthExchangeSuccessResponse{
		User:      *userInfo,
		ExpiresIn: tokenResp.ExpiresIn,
	})
}

// handleAuthMe returns current user info from the access_token cookie.
//
// GET /oauth/me
// Response: { "user": { "sub": "...", "email": "...", "name": "...", "given_name": "...", ... }, "expires_in": 3600 }
func (s *Server) handleAuthMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeAuthError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET method is allowed")
		return
	}

	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		s.writeAuthError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	// Fetch user info from /userinfo endpoint (proper OIDC pattern)
	userInfo, err := s.fetchUserInfo(accessCookie.Value)
	if err != nil {
		s.log.Warn("Failed to fetch userinfo, falling back to JWT", "error", err)
		jwtUserInfo := extractUserInfoFromJWT(accessCookie.Value)
		if jwtUserInfo.Sub == "" {
			s.writeAuthError(w, http.StatusUnauthorized, "unauthorized", "Invalid token")
			return
		}
		userInfo = &jwtUserInfo
	}

	// Calculate remaining expiry from JWT
	expiresIn := getTokenExpirySeconds(accessCookie.Value)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthExchangeSuccessResponse{
		User:      *userInfo,
		ExpiresIn: expiresIn,
	})
}

// exchangeCodeForTokens exchanges an authorization code for tokens with the auth server
func (s *Server) exchangeCodeForTokens(code, codeVerifier, redirectURI string) (*TokenResponse, error) {
	authServerURL := s.cfg.GetDashboardOAuthServerURL()
	clientID := s.cfg.GetDefaultOAuthClientID()
	clientSecret := s.cfg.GetDefaultOAuthClientSecret()

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"code_verifier": {codeVerifier},
	}

	tokenURL := authServerURL + "/oauth/token"
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &AuthServerError{StatusCode: resp.StatusCode}
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// refreshTokens exchanges a refresh token for new tokens with the auth server
func (s *Server) refreshTokens(refreshToken string) (*TokenResponse, error) {
	authServerURL := s.cfg.GetDashboardOAuthServerURL()
	clientID := s.cfg.GetDefaultOAuthClientID()
	clientSecret := s.cfg.GetDefaultOAuthClientSecret()

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	tokenURL := authServerURL + "/oauth/token"
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &AuthServerError{StatusCode: resp.StatusCode}
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// fetchUserInfo calls the OAuth server's /userinfo endpoint to get user profile
func (s *Server) fetchUserInfo(accessToken string) (*AuthUserInfo, error) {
	authServerURL := s.cfg.GetDashboardOAuthServerURL()
	userInfoURL := authServerURL + "/oauth/userinfo"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &AuthServerError{StatusCode: resp.StatusCode}
	}

	var userInfo AuthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// setAuthCookies sets httpOnly cookies for access and refresh tokens
func (s *Server) setAuthCookies(w http.ResponseWriter, tokenResp *TokenResponse) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenResp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   tokenResp.ExpiresIn,
	})

	if tokenResp.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokenResp.RefreshToken,
			Path:     "/oauth",
			HttpOnly: true,
			Secure:   false, // TODO: true in production with HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400 * 7, // 7 days
		})
	}
}

// clearAuthCookies clears auth cookies by setting MaxAge to -1
func (s *Server) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/oauth",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// extractUserInfoFromJWT extracts user info from a JWT without validation
// (we trust the auth server response)
func extractUserInfoFromJWT(accessToken string) AuthUserInfo {
	token, _, _ := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if token == nil {
		return AuthUserInfo{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthUserInfo{}
	}

	var userInfo AuthUserInfo
	if sub, ok := claims["sub"].(string); ok {
		userInfo.Sub = sub
	}
	if email, ok := claims["email"].(string); ok {
		userInfo.Email = email
	}
	if name, ok := claims["name"].(string); ok {
		userInfo.Name = name
	}

	return userInfo
}

// getTokenExpirySeconds returns remaining seconds until token expires
func getTokenExpirySeconds(accessToken string) int {
	token, _, _ := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if token == nil {
		return 0
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return 0
	}

	remaining := int(exp) - int(time.Now().Unix())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// AuthServerError represents an error from the auth server
type AuthServerError struct {
	StatusCode int
}

func (e *AuthServerError) Error() string {
	return "auth server returned status " + http.StatusText(e.StatusCode)
}

// writeAuthError writes an OAuth-style error response
func (s *Server) writeAuthError(w http.ResponseWriter, statusCode int, errorCode, errorDescription string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(AuthErrorResponse{
		Error:            errorCode,
		ErrorDescription: errorDescription,
	})
}
