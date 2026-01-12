# Task T28: Authorization & Consent Endpoints

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Backend Handlers
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T27 (OAuth Provider Login Flow)

## Objective

Implement the /oauth/authorize endpoint for handling OAuth 2.0 authorization requests, including parameter validation, PKCE support, consent screen rendering, and authorization code generation.

## Acceptance Criteria

- [ ] GET /oauth/authorize validates all OAuth parameters
- [ ] Validates response_type is "code"
- [ ] Validates client_id exists and retrieves OAuth client
- [ ] Validates redirect_uri matches registered URIs (exact match)
- [ ] Validates PKCE parameters if client requires PKCE
- [ ] Checks for existing user consent
- [ ] Renders consent screen if no prior consent
- [ ] Skips consent if user already consented to same scopes
- [ ] POST /oauth/authorize processes consent (Allow/Deny)
- [ ] Generates authorization code and redirects on Allow
- [ ] Redirects with error on Deny

## Technical Requirements

### Authorization GET Handler (`handlers/authorize.go`)

```go
func Authorize(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 1. Require authenticated session
        session := session.GetFromContext(ctx)
        if !session.IsAuthenticated() {
            // Store original URL and redirect to login
            session.OriginalURL = r.URL.String()
            session.Save(r, w)
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        // 2. Parse and validate parameters
        params, err := parseAuthorizationParams(r)
        if err != nil {
            renderAuthError(w, params.RedirectURI, params.State, err)
            return
        }

        // 3. Validate response_type
        if params.ResponseType != "code" {
            renderAuthError(w, params.RedirectURI, params.State, ErrUnsupportedResponseType)
            return
        }

        // 4. Get OAuth client by client_id
        client, err := c.OAuthClientService().GetByClientID(ctx, params.ClientID)
        if err != nil {
            // Cannot redirect - render error page
            templates.Render(w, "error.html", ErrorPageData{
                Error:       "invalid_client",
                Description: "Unknown client_id",
            })
            return
        }

        // 5. Validate redirect_uri (exact match)
        if !client.ContainsRedirectURI(params.RedirectURI) {
            templates.Render(w, "error.html", ErrorPageData{
                Error:       "invalid_redirect_uri",
                Description: "Redirect URI does not match registered URIs",
            })
            return
        }

        // 6. Validate PKCE if required
        if client.PKCERequired {
            if params.CodeChallenge == "" {
                renderAuthError(w, params.RedirectURI, params.State, ErrMissingCodeChallenge)
                return
            }
            if params.CodeChallengeMethod != "S256" && params.CodeChallengeMethod != "plain" {
                renderAuthError(w, params.RedirectURI, params.State, ErrInvalidCodeChallengeMethod)
                return
            }
        }

        // 7. Validate and parse scopes
        requestedScopes, err := validateScopes(ctx, c.ScopeService(), params.Scope, client.ID)
        if err != nil {
            renderAuthError(w, params.RedirectURI, params.State, ErrInvalidScope)
            return
        }

        // 8. Check for existing consent
        hasConsent, err := c.OAuthAuthService().CheckUserConsent(ctx, session.UserID, params.ClientID, params.Scope)
        if err != nil {
            renderAuthError(w, params.RedirectURI, params.State, ErrServerError)
            return
        }

        if hasConsent {
            // Skip consent - generate code and redirect
            code, err := c.OAuthAuthService().GenerateAuthorizationCode(ctx, GenerateAuthCodeInput{
                ClientID:            params.ClientID,
                UserID:              session.UserID,
                RedirectURI:         params.RedirectURI,
                Scope:               params.Scope,
                Nonce:               params.Nonce,
                CodeChallenge:       params.CodeChallenge,
                CodeChallengeMethod: params.CodeChallengeMethod,
            })
            if err != nil {
                renderAuthError(w, params.RedirectURI, params.State, ErrServerError)
                return
            }

            redirectWithCode(w, r, params.RedirectURI, code.Code.String(), params.State)
            return
        }

        // 9. Render consent screen
        data := templates.ConsentPageData{
            ClientName:   client.Name,
            Scopes:       requestedScopes,
            CSRFToken:    generateCSRFToken(session),
            ClientID:     params.ClientID.String(),
            RedirectURI:  params.RedirectURI,
            Scope:        params.Scope,
            State:        params.State,
            Nonce:        params.Nonce,
            CodeChallenge:       params.CodeChallenge,
            CodeChallengeMethod: params.CodeChallengeMethod,
        }
        templates.Render(w, "consent.html", data)
    }
}
```

### Authorization POST Handler (same file)

```go
func AuthorizeProcess(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 1. Require authenticated session
        session := session.GetFromContext(ctx)
        if !session.IsAuthenticated() {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        // 2. Parse form
        if err := r.ParseForm(); err != nil {
            http.Error(w, "Invalid form data", http.StatusBadRequest)
            return
        }

        // 3. Validate CSRF token
        csrfToken := r.FormValue("csrf_token")
        if !validateCSRFToken(session, csrfToken) {
            http.Error(w, "Invalid CSRF token", http.StatusForbidden)
            return
        }

        // 4. Get authorization parameters from form
        params := AuthorizationParams{
            ClientID:            uuid.MustParse(r.FormValue("client_id")),
            RedirectURI:         r.FormValue("redirect_uri"),
            Scope:               r.FormValue("scope"),
            State:               r.FormValue("state"),
            Nonce:               r.FormValue("nonce"),
            CodeChallenge:       r.FormValue("code_challenge"),
            CodeChallengeMethod: r.FormValue("code_challenge_method"),
        }

        // 5. Check user's decision
        decision := r.FormValue("decision")

        if decision == "deny" {
            // Redirect with access_denied error
            redirectWithError(w, r, params.RedirectURI, "access_denied", "User denied the request", params.State)
            return
        }

        // 6. User allowed - generate authorization code
        code, err := c.OAuthAuthService().GenerateAuthorizationCode(ctx, GenerateAuthCodeInput{
            ClientID:            params.ClientID,
            UserID:              session.UserID,
            RedirectURI:         params.RedirectURI,
            Scope:               params.Scope,
            Nonce:               params.Nonce,
            CodeChallenge:       params.CodeChallenge,
            CodeChallengeMethod: params.CodeChallengeMethod,
        })
        if err != nil {
            renderAuthError(w, params.RedirectURI, params.State, ErrServerError)
            return
        }

        // 7. Save user consent for future requests
        err = c.OAuthAuthService().SaveUserConsent(ctx, session.UserID, params.ClientID, params.Scope)
        if err != nil {
            // Log but don't fail - consent is optional
            c.Logger().Warn("failed to save user consent", "error", err)
        }

        // 8. Redirect to client with authorization code
        redirectWithCode(w, r, params.RedirectURI, code.Code.String(), params.State)
    }
}
```

### Helper Functions

```go
type AuthorizationParams struct {
    ResponseType        string
    ClientID            uuid.UUID
    RedirectURI         string
    Scope               string
    State               string
    Nonce               string
    CodeChallenge       string
    CodeChallengeMethod string
}

func parseAuthorizationParams(r *http.Request) (*AuthorizationParams, error) {
    params := &AuthorizationParams{
        ResponseType:        r.URL.Query().Get("response_type"),
        RedirectURI:         r.URL.Query().Get("redirect_uri"),
        Scope:               r.URL.Query().Get("scope"),
        State:               r.URL.Query().Get("state"),
        Nonce:               r.URL.Query().Get("nonce"),
        CodeChallenge:       r.URL.Query().Get("code_challenge"),
        CodeChallengeMethod: r.URL.Query().Get("code_challenge_method"),
    }

    // Parse client_id as UUID
    clientIDStr := r.URL.Query().Get("client_id")
    if clientIDStr == "" {
        return nil, ErrMissingClientID
    }
    clientID, err := uuid.Parse(clientIDStr)
    if err != nil {
        return nil, ErrInvalidClientID
    }
    params.ClientID = clientID

    // Validate required parameters
    if params.ResponseType == "" {
        return nil, ErrMissingResponseType
    }
    if params.RedirectURI == "" {
        return nil, ErrMissingRedirectURI
    }

    // Default code_challenge_method to S256 if code_challenge is provided
    if params.CodeChallenge != "" && params.CodeChallengeMethod == "" {
        params.CodeChallengeMethod = "S256"
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
```

### Scope Validation

```go
type ScopeInfo struct {
    Name        string
    Description string
}

func validateScopes(ctx context.Context, scopeService *ScopeService, scopeString string, clientID uuid.UUID) ([]ScopeInfo, error) {
    if scopeString == "" {
        // Return default scopes (openid)
        return []ScopeInfo{{Name: "openid", Description: "Verify your identity"}}, nil
    }

    requestedScopes := strings.Split(scopeString, " ")
    validScopes := []ScopeInfo{}

    for _, scope := range requestedScopes {
        scopeInfo, err := scopeService.GetScopeByName(ctx, scope)
        if err != nil {
            return nil, fmt.Errorf("invalid scope: %s", scope)
        }
        validScopes = append(validScopes, ScopeInfo{
            Name:        scopeInfo.Name,
            Description: scopeInfo.Description,
        })
    }

    return validScopes, nil
}

// Human-readable scope descriptions
var scopeDescriptions = map[string]string{
    "openid":         "Verify your identity",
    "profile":        "Access your profile information (name)",
    "email":          "Access your email address",
    "offline_access": "Access your data while you're offline",
}
```

### CSRF Token Handling

```go
func generateCSRFToken(session *SessionData) string {
    token := generateSecureRandomString(32)
    session.CSRFToken = token
    return token
}

func validateCSRFToken(session *SessionData, token string) bool {
    return session.CSRFToken != "" && session.CSRFToken == token
}
```

## Files to Create

- `internal/authserver/handlers/authorize.go` - Authorization endpoint handlers

## Files to Modify

- `internal/authserver/routes.go` - Update route handlers
- `internal/domain/oauth_auth/service.go` - Ensure consent methods are implemented

## Testing Requirements

- Test with all required parameters
- Test missing parameters (error responses)
- Test invalid client_id
- Test redirect_uri mismatch
- Test PKCE required but not provided
- Test consent screen rendering
- Test allow and deny decisions
- Test existing consent skip

## Commands to Run

```bash
# Build application
make build

# Start auth server
./bin/app serve-auth -c config.yaml

# Test authorization endpoint with minimal params
# (requires prior authentication)
curl "http://localhost:3101/oauth/authorize?\
response_type=code&\
client_id=e730207a-0fce-495d-bac3-6211963ac423&\
redirect_uri=http://localhost:8089/callback&\
scope=openid%20profile%20email&\
state=xyz123"
```

## Validation Checklist

- [ ] GET /oauth/authorize requires authentication
- [ ] Missing response_type returns error
- [ ] Invalid client_id returns error
- [ ] Mismatched redirect_uri returns error
- [ ] Missing PKCE when required returns error
- [ ] Consent screen shows client name and scopes
- [ ] Allow redirects with authorization code
- [ ] Deny redirects with access_denied error
- [ ] State parameter is echoed back
- [ ] CSRF token validates on POST

## Definition of Done

- [ ] Authorization parameter parsing and validation
- [ ] OAuth client lookup by client_id
- [ ] Redirect URI exact match validation
- [ ] PKCE parameter validation when required
- [ ] Scope parsing and validation
- [ ] User consent checking
- [ ] Consent screen rendering with scope descriptions
- [ ] Authorization code generation on allow
- [ ] Error redirect on deny
- [ ] CSRF protection on consent form
- [ ] Consent saved for future requests

## Dependencies

- T27: Login flow (user must be authenticated)
- T25: OAuth auth domain (code generation, consent)
- T30: HTML templates (consent.html)
- Existing: oauth_client domain for client lookup

## Risk Factors

- **Low Risk**: Standard OAuth 2.0 authorization endpoint
- **Medium Risk**: Redirect URI validation must be exact match (security critical)

## Notes

- Redirect URI validation is security-critical - no wildcards, no partial matches
- State parameter is client-side CSRF protection (echoed back)
- CSRF token is server-side form protection
- PKCE is required for public clients (pkce_required=true on client)
- Consent can be remembered per user-client pair
- Authorization code expires in 10 minutes (configurable)
