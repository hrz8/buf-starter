package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"examples/oauth-client-ssr/views"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// Config holds the OAuth client configuration
type Config struct {
	AuthServerURL string
	ClientID      string
	ClientSecret  string
	RedirectURI   string
	Scopes        []string
	Port          int
}

// TokenResponse represents the token endpoint response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// OAuthError represents OAuth error response
type OAuthError struct {
	Error       string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

// SessionStore manages user sessions
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// Session holds user session data
type Session struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	UserInfo     map[string]any
}

// PendingAuth tracks pending OAuth authorization requests
type PendingAuth struct {
	State        string
	CodeVerifier string
	CreatedAt    time.Time
}

var (
	config       *Config
	sessionStore = &SessionStore{sessions: make(map[string]*Session)}
	pendingAuths = &sync.Map{}
)

func main() {
	config = parseFlags()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/public", handlePublic)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/callback", handleCallback)
	mux.HandleFunc("/logout", handleLogout)

	// Protected routes
	mux.HandleFunc("/private/dashboard", requireAuth(handlePrivateDashboard))
	mux.HandleFunc("/private/profile", requireAuth(handlePrivateProfile))
	mux.HandleFunc("/private/settings", requireAuth(handlePrivateSettings))

	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("OAuth Example Client starting at http://localhost%s", addr)
	log.Printf("Authorization Server: %s", config.AuthServerURL)
	log.Printf("Client ID: %s", config.ClientID)
	log.Println()
	log.Println("Open your browser and navigate to the home page to begin")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func parseFlags() *Config {
	authServer := flag.String("auth-server", "http://localhost:3300", "OAuth authorization server URL")
	clientID := flag.String("client-id", "", "OAuth client ID (required)")
	clientSecret := flag.String("client-secret", "", "OAuth client secret (required)")
	port := flag.Int("port", 8080, "Local web server port")
	scopes := flag.String("scopes", "openid profile email", "Space-separated scopes")

	flag.Parse()

	if *clientID == "" || *clientSecret == "" {
		fmt.Println("Usage: oauth-client -client-id <id> -client-secret <secret> [options]")
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(1)
	}

	return &Config{
		AuthServerURL: *authServer,
		ClientID:      *clientID,
		ClientSecret:  *clientSecret,
		RedirectURI:   fmt.Sprintf("http://localhost:%d/callback", *port),
		Scopes:        strings.Split(*scopes, " "),
		Port:          *port,
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	session := sessionStore.Get(sessionID)

	data := map[string]any{
		"Authenticated": session != nil && !session.IsExpired(),
	}

	if session != nil && !session.IsExpired() {
		data["AccessToken"] = truncateToken(session.AccessToken)
		if session.UserInfo != nil {
			userInfoJSON, _ := json.MarshalIndent(session.UserInfo, "", "  ")
			data["UserInfo"] = string(userInfoJSON)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	sessionID := getSessionID(r)
	session := sessionStore.Get(sessionID)

	if session != nil && !session.IsExpired() {
		// Already logged in, redirect to home
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if this is the start of OAuth flow
	if r.URL.Query().Get("start") != "1" {
		// Show login page with button
		renderLoginPage(w)
		return
	}

	// Generate PKCE code verifier and challenge
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Generate state for CSRF protection
	state := generateRandomString(32)

	// Store pending auth
	pendingAuths.Store(state, &PendingAuth{
		State:        state,
		CodeVerifier: codeVerifier,
		CreatedAt:    time.Now(),
	})

	// Build authorization URL
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {config.ClientID},
		"redirect_uri":          {config.RedirectURI},
		"scope":                 {strings.Join(config.Scopes, " ")},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	// Forward prompt parameter if provided
	if prompt := r.URL.Query().Get("prompt"); prompt != "" {
		params.Set("prompt", prompt)
	}

	authURL := fmt.Sprintf("%s/oauth/authorize?%s", config.AuthServerURL, params.Encode())

	// Show loading page before redirect
	renderLoadingPage(w, "Redirecting to Altalune...", authURL)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// Check for error response
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		errDesc := r.URL.Query().Get("error_description")
		renderError(w, fmt.Sprintf("OAuth Error: %s - %s", errParam, errDesc))
		return
	}

	// Get and validate state
	state := r.URL.Query().Get("state")
	if state == "" {
		renderError(w, "Missing state parameter")
		return
	}

	pending, ok := pendingAuths.Load(state)
	if !ok {
		renderError(w, "Invalid or expired state parameter")
		return
	}
	pendingAuths.Delete(state)

	pendingAuth := pending.(*PendingAuth)

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		renderError(w, "No authorization code in response")
		return
	}

	// Exchange code for tokens
	tokens, err := exchangeCodeForTokens(code, pendingAuth.CodeVerifier)
	if err != nil {
		renderError(w, fmt.Sprintf("Token exchange failed: %v", err))
		return
	}

	// Validate access token
	if err := validateAccessToken(tokens.AccessToken); err != nil {
		log.Printf("Warning: JWT validation failed: %v", err)
	}

	// Get user info from JWT claims
	userInfo := extractUserInfo(tokens.AccessToken)

	// Create session
	sessionID := generateRandomString(32)
	session := &Session{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second),
		UserInfo:     userInfo,
	}
	sessionStore.Set(sessionID, session)

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "example_oauthclient_session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   tokens.ExpiresIn,
	})

	// Check for return_to cookie
	returnTo := "/"
	if cookie, err := r.Cookie("return_to"); err == nil && cookie.Value != "" {
		returnTo = cookie.Value
		// Clear the return_to cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "return_to",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}

	// Show success page with redirect
	renderSuccessPage(w, "Authentication Successful!", returnTo)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID != "" {
		sessionStore.Delete(sessionID)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "example_oauthclient_session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

// requireAuth is middleware that protects routes requiring authentication
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := getSessionID(r)
		session := sessionStore.Get(sessionID)

		if session == nil || session.IsExpired() {
			// Store the original URL to redirect back after login
			http.SetCookie(w, &http.Cookie{
				Name:     "return_to",
				Value:    r.URL.Path,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   300, // 5 minutes
			})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}

func handlePublic(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	session := sessionStore.Get(sessionID)

	data := map[string]any{
		"Authenticated": session != nil && !session.IsExpired(),
	}

	if session != nil && !session.IsExpired() {
		if email, ok := session.UserInfo["email"].(string); ok {
			data["UserEmail"] = email
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "public.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePrivateDashboard(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	session := sessionStore.Get(sessionID)

	data := map[string]any{
		"UserName": getUserName(session),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePrivateProfile(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	session := sessionStore.Get(sessionID)

	claimsJSON, _ := json.MarshalIndent(session.UserInfo, "", "  ")
	data := map[string]any{
		"UserName":  getUserName(session),
		"UserEmail": getUserEmail(session),
		"Claims":    string(claimsJSON),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "profile.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePrivateSettings(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	session := sessionStore.Get(sessionID)

	data := map[string]any{
		"AccessToken":  truncateToken(session.AccessToken),
		"RefreshToken": truncateToken(session.RefreshToken),
		"ExpiresAt":    session.ExpiresAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "settings.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func exchangeCodeForTokens(code, codeVerifier string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {config.RedirectURI},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequest("POST", config.AuthServerURL+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.ClientID, config.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var oauthErr OAuthError
		if err := json.Unmarshal(body, &oauthErr); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("%s: %s", oauthErr.Error, oauthErr.Description)
	}

	var tokens TokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	return &tokens, nil
}

func validateAccessToken(accessToken string) error {
	// Fetch JWKS
	jwksURL := config.AuthServerURL + "/.well-known/jwks.json"
	keySet, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %v", err)
	}

	// Parse and validate token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		// Verify algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get key ID from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in token header")
		}

		// Find key in JWKS
		key, found := keySet.LookupKeyID(kid)
		if !found {
			return nil, fmt.Errorf("key %s not found in JWKS", kid)
		}

		// Convert to crypto key
		var rawKey any
		if err := key.Raw(&rawKey); err != nil {
			return nil, fmt.Errorf("failed to get raw key: %v", err)
		}

		return rawKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("token is invalid")
	}

	return nil
}

func extractUserInfo(accessToken string) map[string]any {
	token, _, _ := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if token == nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	return claims
}

func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}

func getSessionID(r *http.Request) string {
	cookie, err := r.Cookie("example_oauthclient_session_id")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func truncateToken(token string) string {
	if len(token) <= 50 {
		return token
	}
	return token[:50] + "..."
}

func getUserName(session *Session) string {
	if session == nil || session.UserInfo == nil {
		return "Unknown"
	}

	// Try to get name from various JWT claim fields
	if name, ok := session.UserInfo["name"].(string); ok && name != "" {
		return name
	}
	if givenName, ok := session.UserInfo["given_name"].(string); ok {
		if familyName, ok := session.UserInfo["family_name"].(string); ok {
			return givenName + " " + familyName
		}
		return givenName
	}
	if email, ok := session.UserInfo["email"].(string); ok {
		return email
	}
	return "Unknown User"
}

func getUserEmail(session *Session) string {
	if session == nil || session.UserInfo == nil {
		return "unknown@example.com"
	}

	if email, ok := session.UserInfo["email"].(string); ok {
		return email
	}
	return "unknown@example.com"
}

func renderLoginPage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderLoadingPage(w http.ResponseWriter, message, redirectURL string) {
	data := map[string]any{
		"Message":     message,
		"RedirectURL": redirectURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "loading.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderSuccessPage(w http.ResponseWriter, message, redirectURL string) {
	data := map[string]any{
		"Message":     message,
		"RedirectURL": redirectURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.Render(w, "success.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderError(w http.ResponseWriter, message string) {
	data := map[string]any{
		"Message": message,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := views.Render(w, "error.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SessionStore methods
func (s *SessionStore) Get(id string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[id]
}

func (s *SessionStore) Set(id string, session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = session
}

func (s *SessionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}

// Session methods
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
