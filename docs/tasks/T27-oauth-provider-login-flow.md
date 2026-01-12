# Task T27: OAuth Provider Login Flow

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Backend Handlers
**Priority:** High
**Estimated Effort:** 5-6 hours
**Prerequisites:** T26 (serve-auth CLI Command & Server)

## Objective

Implement the complete OAuth provider login flow including login page, provider-specific OAuth initiation (Google/GitHub), callback handling with user creation/update, and logout functionality.

## Acceptance Criteria

- [ ] GET /login renders login page with provider buttons
- [ ] GET /login/google initiates Google OAuth flow
- [ ] GET /login/github initiates GitHub OAuth flow
- [ ] GET /auth/callback handles OAuth provider callback
- [ ] New users: creates user + user_identity + project_member (role=user)
- [ ] Existing users: updates user_identity.last_login_at
- [ ] POST /logout destroys session and clears cookies
- [ ] State parameter provides CSRF protection
- [ ] Original URL is preserved for post-login redirect

## Technical Requirements

### Login Page Handler (`handlers/login.go`)

```go
func LoginPage(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check if already authenticated
        session := session.GetFromContext(r.Context())
        if session.IsAuthenticated() {
            // Redirect to authorize or home
            http.Redirect(w, r, "/oauth/authorize", http.StatusFound)
            return
        }

        // Get available providers from database
        providers, err := c.OAuthProviderService().GetEnabledProviders(r.Context())

        // Render login template
        data := templates.LoginPageData{
            Providers:    providers,
            ErrorMessage: r.URL.Query().Get("error"),
        }
        templates.Render(w, "login.html", data)
    }
}
```

### Provider Login Handler (`handlers/login_provider.go`)

```go
func LoginProvider(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        providerName := chi.URLParam(r, "provider")

        // Get provider config from database
        provider, err := c.OAuthProviderService().GetByType(r.Context(), providerName)
        if err != nil {
            // Redirect to login with error
            http.Redirect(w, r, "/login?error=invalid_provider", http.StatusFound)
            return
        }

        // Generate state parameter (CSRF protection)
        state := generateSecureRandomString(32)

        // Store state and original URL in session
        session := session.GetFromContext(r.Context())
        session.OAuthState = state
        session.OriginalURL = r.URL.Query().Get("next")
        session.Save(r, w)

        // Build authorization URL
        authURL := buildAuthorizationURL(provider, state)

        // Redirect to OAuth provider
        http.Redirect(w, r, authURL, http.StatusFound)
    }
}
```

### OAuth Callback Handler (`handlers/callback.go`)

```go
func OAuthCallback(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Validate state parameter
        state := r.URL.Query().Get("state")
        session := session.GetFromContext(r.Context())
        if state != session.OAuthState {
            http.Redirect(w, r, "/login?error=invalid_state", http.StatusFound)
            return
        }

        // 2. Get authorization code
        code := r.URL.Query().Get("code")
        if code == "" {
            error := r.URL.Query().Get("error")
            http.Redirect(w, r, "/login?error="+error, http.StatusFound)
            return
        }

        // 3. Determine provider from session or cookie
        provider := session.Provider // Set during login initiation

        // 4. Exchange code for user info
        userInfo, err := c.OAuthProviderClient().ExchangeCodeForUserInfo(r.Context(), provider, code)
        if err != nil {
            http.Redirect(w, r, "/login?error=token_exchange_failed", http.StatusFound)
            return
        }

        // 5. Find or create user
        user, identity, isNew, err := c.UserRegistrationService().FindOrCreateUser(r.Context(), FindOrCreateUserInput{
            Email:          userInfo.Email,
            FirstName:      userInfo.FirstName,
            LastName:       userInfo.LastName,
            Provider:       provider,
            ProviderUserID: userInfo.ID,
            OAuthClientID:  getOAuthClientIDFromSession(session), // Client that initiated the flow
        })

        // 6. Update session with authenticated user
        session.UserID = user.ID
        session.AuthenticatedAt = time.Now()
        session.OAuthState = "" // Clear state
        session.Save(r, w)

        // 7. Redirect to original URL or authorize endpoint
        redirectURL := session.OriginalURL
        if redirectURL == "" {
            redirectURL = "/oauth/authorize"
        }
        http.Redirect(w, r, redirectURL, http.StatusFound)
    }
}
```

### Logout Handler (`handlers/logout.go`)

```go
func Logout(c *Container) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Destroy session
        session := session.GetFromContext(r.Context())
        session.Destroy(r, w)

        // Render logout success page
        templates.Render(w, "logout_success.html", nil)
    }
}
```

### OAuth Provider Client (`internal/shared/oauthprovider/`)

#### Client Interface (`client.go`)

```go
type ProviderClient interface {
    ExchangeCodeForUserInfo(ctx context.Context, code string) (*UserInfo, error)
}

type UserInfo struct {
    ID        string // Provider's unique user ID
    Email     string
    FirstName string
    LastName  string
    AvatarURL string
}
```

#### Google Client (`google.go`)

```go
type GoogleClient struct {
    config *oauth2.Config
}

func NewGoogleClient(clientID, clientSecret, redirectURL string) *GoogleClient {
    return &GoogleClient{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Scopes:       []string{"openid", "profile", "email"},
            Endpoint:     google.Endpoint,
        },
    }
}

func (c *GoogleClient) ExchangeCodeForUserInfo(ctx context.Context, code string) (*UserInfo, error) {
    // Exchange code for token
    token, err := c.config.Exchange(ctx, code)
    if err != nil {
        return nil, fmt.Errorf("token exchange: %w", err)
    }

    // Fetch user info from Google
    client := c.config.Client(ctx, token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    // Parse response...

    return &UserInfo{
        ID:        googleUserInfo.ID,
        Email:     googleUserInfo.Email,
        FirstName: googleUserInfo.GivenName,
        LastName:  googleUserInfo.FamilyName,
        AvatarURL: googleUserInfo.Picture,
    }, nil
}
```

#### GitHub Client (`github.go`)

```go
type GitHubClient struct {
    config *oauth2.Config
}

func NewGitHubClient(clientID, clientSecret, redirectURL string) *GitHubClient {
    return &GitHubClient{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Scopes:       []string{"read:user", "user:email"},
            Endpoint:     github.Endpoint,
        },
    }
}

func (c *GitHubClient) ExchangeCodeForUserInfo(ctx context.Context, code string) (*UserInfo, error) {
    // Exchange code for token
    token, err := c.config.Exchange(ctx, code)

    // Fetch user info from GitHub API
    client := c.config.Client(ctx, token)
    resp, err := client.Get("https://api.github.com/user")
    // Parse response...

    // GitHub might not return email in user endpoint, need to fetch separately
    emailResp, err := client.Get("https://api.github.com/user/emails")
    // Find primary email...

    return &UserInfo{
        ID:        strconv.Itoa(githubUser.ID),
        Email:     primaryEmail,
        FirstName: parseName(githubUser.Name).First,
        LastName:  parseName(githubUser.Name).Last,
        AvatarURL: githubUser.AvatarURL,
    }, nil
}
```

### User Registration Logic

```go
// internal/authserver/services/user_registration.go
type UserRegistrationService struct {
    userRepo            *user.Repo
    userIdentityRepo    *useridentity.Repo
    projectMemberRepo   *projectmember.Repo
    oauthClientRepo     *oauth_client.Repo
}

func (s *UserRegistrationService) FindOrCreateUser(ctx context.Context, input FindOrCreateUserInput) (*User, *UserIdentity, bool, error) {
    // 1. Check if user exists by email
    existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)

    if err == user.ErrUserNotFound {
        // NEW USER
        // 1a. Create user
        newUser, err := s.userRepo.Create(ctx, CreateUserInput{
            Email:     input.Email,
            FirstName: input.FirstName,
            LastName:  input.LastName,
        })

        // 1b. Create user identity
        identity, err := s.userIdentityRepo.Create(ctx, CreateIdentityInput{
            UserID:         newUser.ID,
            Provider:       input.Provider,
            ProviderUserID: input.ProviderUserID,
            Email:          input.Email,
            FirstName:      input.FirstName,
            LastName:       input.LastName,
            OAuthClientID:  input.OAuthClientID, // Links to OAuth client
            LastLoginAt:    time.Now(),
        })

        // 1c. Get OAuth client's project_id
        client, err := s.oauthClientRepo.GetByClientID(ctx, input.OAuthClientID)

        // 1d. Create project membership with role='user'
        _, err = s.projectMemberRepo.Create(ctx, CreateMemberInput{
            ProjectID: client.ProjectID,
            UserID:    newUser.ID,
            Role:      "user",
        })

        return newUser, identity, true, nil
    }

    if err != nil {
        return nil, nil, false, err
    }

    // EXISTING USER
    // 2a. Check if identity for this provider exists
    identity, err := s.userIdentityRepo.GetByUserAndProvider(ctx, existingUser.ID, input.Provider)

    if err == useridentity.ErrNotFound {
        // Create new identity for existing user (different provider)
        identity, err = s.userIdentityRepo.Create(ctx, CreateIdentityInput{
            UserID:         existingUser.ID,
            Provider:       input.Provider,
            ProviderUserID: input.ProviderUserID,
            Email:          input.Email,
            OAuthClientID:  input.OAuthClientID,
            LastLoginAt:    time.Now(),
        })
    } else {
        // Update last_login_at
        err = s.userIdentityRepo.UpdateLastLogin(ctx, identity.ID, time.Now())
    }

    return existingUser, identity, false, nil
}
```

## Files to Create

- `internal/authserver/handlers/login.go` - Login page handler
- `internal/authserver/handlers/login_provider.go` - Provider OAuth initiation
- `internal/authserver/handlers/callback.go` - OAuth callback handler
- `internal/authserver/handlers/logout.go` - Logout handler
- `internal/shared/oauthprovider/client.go` - Provider client interface
- `internal/shared/oauthprovider/google.go` - Google OAuth client
- `internal/shared/oauthprovider/github.go` - GitHub OAuth client
- `internal/authserver/services/user_registration.go` - User registration logic

## Files to Modify

- `internal/authserver/routes.go` - Connect handlers to routes
- `internal/authserver/container.go` - Add provider clients and registration service

## Testing Requirements

- Manual testing with real Google/GitHub OAuth credentials
- Test state parameter validation (CSRF protection)
- Test new user creation flow
- Test existing user login flow
- Test logout clears session

## Commands to Run

```bash
# Build application
make build

# Start auth server
./bin/app serve-auth -c config.yaml

# Open browser to test login
open http://localhost:3101/login

# Test login flow:
# 1. Click "Continue with Google"
# 2. Authenticate with Google
# 3. Should redirect back and create session
```

## Validation Checklist

- [ ] Login page shows available providers
- [ ] State parameter is cryptographically random
- [ ] State is validated on callback
- [ ] Google OAuth flow completes successfully
- [ ] GitHub OAuth flow completes successfully
- [ ] New user creates all required records
- [ ] Existing user updates last_login_at
- [ ] Session contains user_id after login
- [ ] Logout destroys session

## Definition of Done

- [ ] Login page renders with provider buttons
- [ ] Google OAuth initiation works
- [ ] GitHub OAuth initiation works
- [ ] Callback validates state parameter
- [ ] Callback exchanges code for user info
- [ ] New users are created with identity and project membership
- [ ] Existing users have identity updated
- [ ] Session is created with authenticated user
- [ ] Logout destroys session and renders success page
- [ ] Error handling redirects to /login with error message

## Dependencies

- T26: serve-auth server infrastructure
- Existing: oauth_provider domain for provider configs
- Existing: user domain for user operations
- Existing: oauth_client domain for client lookup
- `golang.org/x/oauth2` - OAuth2 client library

## Risk Factors

- **Medium Risk**: Provider API changes - use well-maintained oauth2 library
- **Medium Risk**: GitHub email might be private - handle email fetch separately
- **Low Risk**: Session management - using battle-tested Gorilla sessions

## Notes

- Google requires scopes: openid, profile, email
- GitHub requires scopes: read:user, user:email
- GitHub doesn't return email directly, need separate API call
- Provider credentials come from database (oauth_providers table)
- OAuth client ID in session determines which app user is signing into
- Project membership is created for the OAuth client's project
