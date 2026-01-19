# Task T33: Example OAuth Client for Testing

**Story Reference:** US7-oauth-authorization-server.md, US9-oauth-testing-example-client.md
**Type:** Testing/Example
**Priority:** Low
**Estimated Effort:** 2-3 hours
**Prerequisites:** T29 (Token Endpoint & JWKS)

## Objective

Create a minimal Go-based OAuth client application for testing the OAuth authorization server. This client demonstrates the complete authorization code flow with PKCE and validates that all OAuth endpoints work correctly.

## Acceptance Criteria

- [ ] Example client lives in `examples/oauth-client/`
- [ ] Implements OAuth 2.0 authorization code flow with PKCE
- [ ] Configurable via command-line flags
- [ ] Starts local HTTP server for callback
- [ ] Opens browser for authorization (or displays URL)
- [ ] Exchanges code for tokens
- [ ] Displays received tokens
- [ ] Validates JWT access token signature against JWKS
- [ ] Demonstrates token refresh flow
- [ ] Includes clear documentation

## Technical Requirements

### Project Structure

```
examples/oauth-client/
├── main.go           # Main client application
├── go.mod            # Module definition
├── go.sum            # Dependencies
└── README.md         # Usage documentation
```

### Main Client Implementation

**File:** `examples/oauth-client/main.go`

```go
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
	"os/exec"
	"runtime"
	"strings"
	"time"

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
	CallbackPort  int
}

// TokenResponse represents the token endpoint response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// OAuthError represents OAuth error response
type OAuthError struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func main() {
	// Parse command-line flags
	config := parseFlags()

	// Generate PKCE code verifier and challenge
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Generate state for CSRF protection
	state := generateRandomString(32)

	// Build authorization URL
	authURL := buildAuthorizationURL(config, codeChallenge, state)

	// Create channels for callback handling
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Start local callback server
	server := startCallbackServer(config.CallbackPort, state, codeChan, errChan)

	// Print instructions
	fmt.Println("=== OAuth Client Test ===")
	fmt.Println()
	fmt.Println("Opening browser for authorization...")
	fmt.Println()
	fmt.Println("Authorization URL:")
	fmt.Println(authURL)
	fmt.Println()

	// Try to open browser
	openBrowser(authURL)

	fmt.Println("Waiting for callback...")
	fmt.Println()

	// Wait for callback or timeout
	var authCode string
	select {
	case authCode = <-codeChan:
		fmt.Println("Authorization code received!")
	case err := <-errChan:
		log.Fatalf("Error: %v", err)
	case <-time.After(5 * time.Minute):
		log.Fatal("Timeout waiting for authorization callback")
	}

	// Shutdown callback server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// Exchange code for tokens
	fmt.Println()
	fmt.Println("Exchanging authorization code for tokens...")
	fmt.Println()

	tokens, err := exchangeCodeForTokens(config, authCode, codeVerifier)
	if err != nil {
		log.Fatalf("Token exchange failed: %v", err)
	}

	// Display tokens
	fmt.Println("=== Tokens Received ===")
	fmt.Println()
	fmt.Printf("Access Token: %s...\n", truncateString(tokens.AccessToken, 50))
	fmt.Printf("Token Type: %s\n", tokens.TokenType)
	fmt.Printf("Expires In: %d seconds\n", tokens.ExpiresIn)
	fmt.Printf("Refresh Token: %s\n", tokens.RefreshToken)
	fmt.Printf("Scope: %s\n", tokens.Scope)
	fmt.Println()

	// Validate JWT signature
	fmt.Println("=== JWT Validation ===")
	fmt.Println()
	if err := validateAccessToken(config.AuthServerURL, tokens.AccessToken); err != nil {
		fmt.Printf("JWT validation failed: %v\n", err)
	} else {
		fmt.Println("JWT signature verified successfully!")
		printJWTClaims(tokens.AccessToken)
	}

	// Test token refresh
	fmt.Println()
	fmt.Println("=== Token Refresh Test ===")
	fmt.Println()
	fmt.Println("Refreshing tokens...")

	newTokens, err := refreshTokens(config, tokens.RefreshToken)
	if err != nil {
		fmt.Printf("Token refresh failed: %v\n", err)
	} else {
		fmt.Println("Token refresh successful!")
		fmt.Printf("New Access Token: %s...\n", truncateString(newTokens.AccessToken, 50))
		fmt.Printf("New Refresh Token: %s\n", newTokens.RefreshToken)
	}

	fmt.Println()
	fmt.Println("=== Test Complete ===")
}

func parseFlags() *Config {
	authServer := flag.String("auth-server", "http://localhost:3101", "OAuth authorization server URL")
	clientID := flag.String("client-id", "", "OAuth client ID (required)")
	clientSecret := flag.String("client-secret", "", "OAuth client secret (required)")
	port := flag.Int("port", 8085, "Local callback server port")
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
		CallbackPort:  *port,
	}
}

func generateCodeVerifier() string {
	// Generate 32 random bytes
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	// S256: BASE64URL(SHA256(code_verifier))
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}

func buildAuthorizationURL(config *Config, codeChallenge, state string) string {
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {config.ClientID},
		"redirect_uri":          {config.RedirectURI},
		"scope":                 {strings.Join(config.Scopes, " ")},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	return fmt.Sprintf("%s/oauth/authorize?%s", config.AuthServerURL, params.Encode())
}

func startCallbackServer(port int, expectedState string, codeChan chan string, errChan chan error) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Check for error response
		if errParam := r.URL.Query().Get("error"); errParam != "" {
			errDesc := r.URL.Query().Get("error_description")
			errChan <- fmt.Errorf("%s: %s", errParam, errDesc)
			w.Write([]byte("Authorization failed. Check terminal for details."))
			return
		}

		// Validate state
		state := r.URL.Query().Get("state")
		if state != expectedState {
			errChan <- fmt.Errorf("invalid state parameter")
			w.Write([]byte("Error: Invalid state parameter"))
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code in response")
			w.Write([]byte("Error: No authorization code received"))
			return
		}

		// Send code to main goroutine
		codeChan <- code

		// Return success page
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head><title>Authorization Successful</title></head>
			<body style="font-family: sans-serif; text-align: center; padding: 50px;">
				<h1>Authorization Successful!</h1>
				<p>You can close this window and return to the terminal.</p>
			</body>
			</html>
		`))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- fmt.Errorf("callback server error: %v", err)
		}
	}()

	return server
}

func exchangeCodeForTokens(config *Config, code, codeVerifier string) (*TokenResponse, error) {
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

func refreshTokens(config *Config, refreshToken string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
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

func validateAccessToken(authServerURL, accessToken string) error {
	// Fetch JWKS
	jwksURL := authServerURL + "/.well-known/jwks.json"
	keySet, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %v", err)
	}

	// Parse and validate token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
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
		var rawKey interface{}
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

func printJWTClaims(accessToken string) {
	// Parse without validation to print claims
	token, _, _ := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if token == nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	fmt.Println()
	fmt.Println("JWT Claims:")
	for key, value := range claims {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return
	}

	cmd.Start()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
```

### Module Definition

**File:** `examples/oauth-client/go.mod`

```go
module github.com/hrz8/buf-starter/examples/oauth-client

go 1.23

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/lestrrat-go/jwx/v2 v2.1.3
)
```

### README Documentation

**File:** `examples/oauth-client/README.md`

```markdown
# OAuth Client Example

A minimal OAuth 2.0 client for testing the Altalune OAuth authorization server.

## Features

- OAuth 2.0 Authorization Code flow
- PKCE (S256) support
- JWT signature validation via JWKS
- Token refresh demonstration

## Prerequisites

1. OAuth authorization server running (`./bin/app serve-auth -c config.yaml`)
2. An OAuth client registered in the system with:
   - `redirect_uri` set to `http://localhost:8085/callback`
   - `pkce_required` set to `true`
3. Client ID and Client Secret

## Usage

```bash
# Build the client
cd examples/oauth-client
go build -o oauth-client .

# Run with required parameters
./oauth-client \
  -client-id "your-client-id" \
  -client-secret "your-client-secret"

# Optional parameters
./oauth-client \
  -client-id "your-client-id" \
  -client-secret "your-client-secret" \
  -auth-server "http://localhost:3101" \
  -port 8085 \
  -scopes "openid profile email"
```

## Command-Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-client-id` | (required) | OAuth client ID |
| `-client-secret` | (required) | OAuth client secret |
| `-auth-server` | `http://localhost:3101` | OAuth authorization server URL |
| `-port` | `8085` | Local callback server port |
| `-scopes` | `openid profile email` | Space-separated OAuth scopes |

## Flow

1. Client generates PKCE code verifier and challenge
2. Opens browser to authorization server's `/oauth/authorize`
3. User authenticates with Google/GitHub
4. User grants consent to the application
5. Authorization server redirects to local callback with code
6. Client exchanges code for tokens at `/oauth/token`
7. Client validates JWT signature using JWKS
8. Client demonstrates token refresh

## Example Output

```
=== OAuth Client Test ===

Opening browser for authorization...

Authorization URL:
http://localhost:3101/oauth/authorize?response_type=code&client_id=...

Waiting for callback...

Authorization code received!

Exchanging authorization code for tokens...

=== Tokens Received ===

Access Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6...
Token Type: Bearer
Expires In: 3600 seconds
Refresh Token: 550e8400-e29b-41d4-a716-446655440000
Scope: openid profile email

=== JWT Validation ===

JWT signature verified successfully!

JWT Claims:
  iss: altalune-oauth
  sub: 123
  aud: your-client-id
  exp: 1234567890
  iat: 1234564290
  scope: openid profile email
  email: user@example.com
  name: John Doe

=== Token Refresh Test ===

Refreshing tokens...
Token refresh successful!
New Access Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIs...
New Refresh Token: 660e9411-f3ac-52e5-b827-557766551111

=== Test Complete ===
```

## Troubleshooting

### "No partition found" error
Make sure the OAuth client's project has partitions created for all required tables.

### "Invalid redirect_uri" error
Ensure the OAuth client has `http://localhost:8085/callback` registered as a valid redirect URI.

### "Invalid client" error
Check that client_id and client_secret are correct.

### Browser doesn't open
Manually copy the authorization URL and paste it in your browser.
```

## Files to Create

- `examples/oauth-client/main.go`
- `examples/oauth-client/go.mod`
- `examples/oauth-client/go.sum` (generated by go mod tidy)
- `examples/oauth-client/README.md`

## Testing Requirements

- Test complete authorization code flow
- Test PKCE validation
- Test JWT signature verification
- Test token refresh
- Test error scenarios (invalid client, expired code, etc.)
- Test with different scopes

## Commands to Run

```bash
# Initialize module
cd examples/oauth-client
go mod init github.com/hrz8/buf-starter/examples/oauth-client
go mod tidy

# Build
go build -o oauth-client .

# Run test (replace with actual credentials)
./oauth-client \
  -client-id "your-oauth-client-id" \
  -client-secret "your-oauth-client-secret"
```

## Validation Checklist

- [ ] Client builds without errors
- [ ] Starts callback server on specified port
- [ ] Opens browser or displays URL
- [ ] Authorization URL includes all required parameters
- [ ] PKCE code_challenge is correctly generated
- [ ] State parameter validated on callback
- [ ] Code exchanged for tokens successfully
- [ ] JWT signature validated against JWKS
- [ ] Token refresh works correctly
- [ ] Error messages are clear and helpful
- [ ] README documentation is complete

## Definition of Done

- [ ] Example client code implemented
- [ ] go.mod and go.sum created
- [ ] README with usage instructions
- [ ] All command-line flags working
- [ ] PKCE flow implemented (S256)
- [ ] JWT validation via JWKS
- [ ] Token refresh demonstrated
- [ ] Error handling for common scenarios
- [ ] Browser auto-open on macOS/Linux/Windows
- [ ] Clear console output for debugging

## Dependencies

- T29: Token Endpoint & JWKS (provides the endpoints to test against)
- `github.com/golang-jwt/jwt/v5` - JWT parsing
- `github.com/lestrrat-go/jwx/v2` - JWKS fetching

## Risk Factors

- **Low Risk**: Standalone test utility
- **Low Risk**: No production impact

## Notes

- This is a development/testing tool only
- Not meant for production use
- Credentials should not be committed
- Can be extended in US9 for more comprehensive testing
- Consider adding a simple web UI version for non-CLI users
- The 5-minute timeout should be sufficient for manual testing
