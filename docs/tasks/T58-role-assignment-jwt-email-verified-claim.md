# Task T58: Role Assignment & JWT email_verified Claim

**Story Reference:** US14-standalone-idp-application.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** T56 (Backend Domain), T57 (Auth Server Pages)

## Objective

Implement context-aware role assignment based on registration method and add `email_verified` claim to JWT access tokens. Also implement user activation logic that sends verification emails.

## Acceptance Criteria

- [ ] JWT access tokens include `email_verified` boolean claim
- [ ] Role assignment based on registration context (standalone/dashboard/custom client)
- [ ] User activation sets `activated_at` timestamp (first activation only)
- [ ] User activation triggers verification email automatically
- [ ] Refresh token returns updated `email_verified` status
- [ ] `autoActivate` config controls initial activation behavior

## Technical Requirements

### JWT Claims Extension

Modify `internal/shared/jwt/claims.go`:

```go
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    Scope         string   `json:"scope,omitempty"`
    Email         string   `json:"email,omitempty"`
    Name          string   `json:"name,omitempty"`
    Perms         []string `json:"perms,omitempty"`
    EmailVerified bool     `json:"email_verified"` // NEW
}
```

### GenerateTokenParams Extension

Modify `internal/shared/jwt/jwt.go`:

```go
type GenerateTokenParams struct {
    UserPublicID  string
    ClientID      string
    Scope         string
    Email         string
    Name          string
    Perms         []string
    EmailVerified bool   // NEW
    Expiry        time.Duration
}

func (s *Signer) GenerateAccessToken(params GenerateTokenParams) (string, error) {
    claims := AccessTokenClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    Issuer,
            Subject:   params.UserPublicID,
            Audience:  jwt.ClaimStrings{params.ClientID},
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
        Scope:         params.Scope,
        Email:         params.Email,
        Name:          params.Name,
        Perms:         params.Perms,
        EmailVerified: params.EmailVerified, // NEW
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    token.Header["kid"] = s.kid
    return token.SignedString(s.privateKey)
}
```

### Role Assignment Matrix

| Registration Method | Project Role | Global Role | Logic |
|---------------------|-------------|-------------|-------|
| Standalone IDP (no client_id) | user | user | Default project (id=1), basic access |
| Dashboard OAuth client | member | user | Default project (id=1), dashboard access |
| Custom OAuth client | user | user | Client's project, basic access |
| Admin creates user | member | user | Selected project, dashboard access |

### Registration Context Detection

Add to `internal/domain/oauth_auth/service.go`:

```go
type RegistrationContext string

const (
    RegistrationContextStandalone RegistrationContext = "standalone"
    RegistrationContextDashboard  RegistrationContext = "dashboard"
    RegistrationContextCustom     RegistrationContext = "custom"
    RegistrationContextAdmin      RegistrationContext = "admin"
)

// DetermineRegistrationContext determines how user registered
func (s *Service) DetermineRegistrationContext(clientID string, dashboardClientID string) RegistrationContext {
    if clientID == "" {
        return RegistrationContextStandalone
    }
    if clientID == dashboardClientID {
        return RegistrationContextDashboard
    }
    return RegistrationContextCustom
}

// GetRoleForContext returns the appropriate project role
func GetProjectRoleForContext(ctx RegistrationContext) string {
    switch ctx {
    case RegistrationContextDashboard, RegistrationContextAdmin:
        return "member"
    default:
        return "user"
    }
}
```

### User Creation with Role Assignment

Modify user creation in `internal/domain/oauth_auth/handler.go` (HandleOAuthCallback):

```go
func (h *Handler) createOrUpdateUser(ctx context.Context, oauthUser OAuthUserInfo, clientID string) (*User, error) {
    // Check if user exists
    existingUser, err := h.userRepo.GetByEmail(ctx, oauthUser.Email)
    if err == nil && existingUser != nil {
        // Update existing user (e.g., update name from OAuth)
        return existingUser, nil
    }

    // Determine registration context
    regCtx := h.svc.DetermineRegistrationContext(clientID, h.cfg.DashboardClientID)
    projectRole := GetProjectRoleForContext(regCtx)

    // Determine if auto-activate
    isActive := h.cfg.AutoActivate // From config

    // Create user
    user, err := h.userRepo.Create(ctx, CreateUserInput{
        Email:         oauthUser.Email,
        FirstName:     oauthUser.FirstName,
        LastName:      oauthUser.LastName,
        IsActive:      isActive,
        EmailVerified: false, // Always starts unverified
    })
    if err != nil {
        return nil, err
    }

    // Assign global 'user' role
    if err := h.roleService.AssignGlobalRole(ctx, user.ID, "user"); err != nil {
        h.log.Warn("failed to assign global role", "error", err)
    }

    // Assign to project with appropriate role
    projectID := h.getProjectIDForContext(ctx, regCtx, clientID)
    if err := h.projectMemberRepo.AddMember(ctx, projectID, user.ID, projectRole); err != nil {
        h.log.Warn("failed to add project member", "error", err)
    }

    // Send verification email if auto-activated
    if isActive {
        if err := h.verificationService.GenerateAndSendVerificationEmail(ctx, user.ID); err != nil {
            h.log.Warn("failed to send verification email", "error", err)
        }
    }

    return user, nil
}

func (h *Handler) getProjectIDForContext(ctx context.Context, regCtx RegistrationContext, clientID string) int64 {
    switch regCtx {
    case RegistrationContextStandalone, RegistrationContextDashboard:
        return 1 // Default project
    case RegistrationContextCustom:
        // Get project from OAuth client
        client, _ := h.oauthClientRepo.GetByClientID(ctx, clientID)
        if client != nil && client.ProjectID != 0 {
            return client.ProjectID
        }
        return 1
    default:
        return 1
    }
}
```

### User Activation with Email

Add to `internal/domain/user/service.go`:

```go
// ActivateUser activates a user and sends verification email
func (s *Service) ActivateUser(ctx context.Context, userID int64) error {
    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        return err
    }

    if user.IsActive {
        return nil // Already active
    }

    // Update user
    now := time.Now()
    updates := UpdateUserInput{
        IsActive:    true,
        ActivatedAt: &now, // Only set on first activation
    }

    if err := s.repo.Update(ctx, userID, updates); err != nil {
        return err
    }

    // Assign global 'user' role if not already assigned
    if err := s.roleService.EnsureGlobalRole(ctx, userID, "user"); err != nil {
        s.log.Warn("failed to ensure global role", "error", err)
    }

    // Assign to default project if not already member
    if err := s.projectService.EnsureProjectMembership(ctx, userID, 1, "user"); err != nil {
        s.log.Warn("failed to ensure project membership", "error", err)
    }

    // Send verification email
    if err := s.verificationService.GenerateAndSendVerificationEmail(ctx, userID); err != nil {
        s.log.Warn("failed to send verification email", "error", err)
    }

    return nil
}
```

### User Repository Extension

Add to `internal/domain/user/repo.go`:

```go
func (r *Repo) SetEmailVerified(ctx context.Context, userID int64, verified bool) error {
    query := `UPDATE altalune_users SET email_verified = $1, updated_at = NOW() WHERE id = $2`
    _, err := r.db.Exec(ctx, query, verified, userID)
    return err
}

func (r *Repo) GetEmailVerified(ctx context.Context, userID int64) (bool, error) {
    query := `SELECT email_verified FROM altalune_users WHERE id = $1`
    var verified bool
    err := r.db.QueryRow(ctx, query, userID).Scan(&verified)
    return verified, err
}

func (r *Repo) SetActivatedAt(ctx context.Context, userID int64, activatedAt time.Time) error {
    // Only set if not already set
    query := `UPDATE altalune_users SET activated_at = $1, updated_at = NOW() WHERE id = $2 AND activated_at IS NULL`
    _, err := r.db.Exec(ctx, query, activatedAt, userID)
    return err
}
```

### Token Generation with email_verified

Update token generation in `internal/domain/oauth_auth/service.go`:

```go
func (s *Service) generateTokenPair(ctx context.Context, user *User, client *OAuthClient, scope string) (*TokenPair, error) {
    // Get user permissions
    perms, _ := s.permissionProvider.GetUserPermissions(ctx, user.ID)

    // Get email_verified status
    emailVerified, _ := s.userRepo.GetEmailVerified(ctx, user.ID)

    // Generate access token with email_verified
    accessToken, err := s.jwtSigner.GenerateAccessToken(jwt.GenerateTokenParams{
        UserPublicID:  user.PublicID,
        ClientID:      client.ClientID,
        Scope:         scope,
        Email:         user.Email,
        Name:          fmt.Sprintf("%s %s", user.FirstName, user.LastName),
        Perms:         perms,
        EmailVerified: emailVerified, // NEW
        Expiry:        time.Duration(s.cfg.AccessTokenExpiry) * time.Second,
    })
    if err != nil {
        return nil, err
    }

    // Generate refresh token (existing logic)
    // ...

    return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
```

### Configuration Extension

Add to `internal/config/app.go`:

```go
type AuthConfig struct {
    Host               string `yaml:"host" validate:"required,hostname|ip"`
    Port               int    `yaml:"port" validate:"required,gte=1,lte=65535"`
    SessionSecret      string `yaml:"sessionSecret" validate:"required,min=32"`
    CodeExpiry         int    `yaml:"codeExpiry" validate:"gte=1"`
    AccessTokenExpiry  int    `yaml:"accessTokenExpiry" validate:"gte=1"`
    RefreshTokenExpiry int    `yaml:"refreshTokenExpiry" validate:"gte=1"`
    AutoActivate       bool   `yaml:"autoActivate"` // NEW: default true
    EmailVerificationExpiry int `yaml:"emailVerificationExpiry"` // NEW: seconds, default 86400
    OTPExpiry          int    `yaml:"otpExpiry"` // NEW: seconds, default 300
    OTPRateLimit       int    `yaml:"otpRateLimit"` // NEW: default 3
    OTPRateLimitWindow int    `yaml:"otpRateLimitWindow"` // NEW: seconds, default 900
}

func (c *AuthConfig) setDefaults() {
    // ... existing defaults ...

    // Default to auto-activate
    // Note: bool defaults to false, so we need explicit handling
    // Consider using *bool or separate defaultAutoActivate field

    if c.EmailVerificationExpiry == 0 {
        c.EmailVerificationExpiry = 86400 // 24 hours
    }
    if c.OTPExpiry == 0 {
        c.OTPExpiry = 300 // 5 minutes
    }
    if c.OTPRateLimit == 0 {
        c.OTPRateLimit = 3
    }
    if c.OTPRateLimitWindow == 0 {
        c.OTPRateLimitWindow = 900 // 15 minutes
    }
}
```

## Files to Create

- None (all modifications to existing files)

## Files to Modify

- `internal/shared/jwt/claims.go` - Add `EmailVerified` field
- `internal/shared/jwt/jwt.go` - Add `EmailVerified` to `GenerateTokenParams`
- `internal/domain/oauth_auth/service.go` - Add registration context, update token generation
- `internal/domain/oauth_auth/handler.go` - Add role assignment on user creation
- `internal/domain/user/repo.go` - Add `SetEmailVerified`, `GetEmailVerified`, `SetActivatedAt`
- `internal/domain/user/service.go` - Add `ActivateUser` method
- `internal/config/app.go` - Add new auth config fields
- `config.yaml` - Add new auth settings
- `config.example.yaml` - Add new auth settings with examples

## Testing Requirements

```go
func TestJWTEmailVerifiedClaim() {
    signer := jwt.NewSigner(privateKey, publicKey, "test-kid")

    // Test with verified email
    token, _ := signer.GenerateAccessToken(jwt.GenerateTokenParams{
        UserPublicID:  "user123",
        ClientID:      "client-uuid",
        EmailVerified: true,
        // ...
    })

    claims, _ := signer.ValidateAccessToken(token)
    assert.True(t, claims.EmailVerified)
}

func TestRoleAssignmentContext() {
    tests := []struct {
        clientID          string
        dashboardClientID string
        expectedContext   RegistrationContext
        expectedRole      string
    }{
        {"", "", RegistrationContextStandalone, "user"},
        {"dashboard-id", "dashboard-id", RegistrationContextDashboard, "member"},
        {"other-id", "dashboard-id", RegistrationContextCustom, "user"},
    }
    // Run tests
}
```

## Commands to Run

```bash
# Build to verify compilation
make build

# Run tests
go test ./internal/shared/jwt/...
go test ./internal/domain/oauth_auth/...
go test ./internal/domain/user/...
```

## Validation Checklist

- [ ] JWT includes `email_verified` claim
- [ ] Claim is boolean (not string)
- [ ] Standalone registration assigns 'user' project role
- [ ] Dashboard registration assigns 'member' project role
- [ ] Custom client registration assigns 'user' role to client's project
- [ ] User activation sets `activated_at` only once
- [ ] User activation sends verification email
- [ ] Refresh token returns updated `email_verified` status
- [ ] Config defaults work correctly

## Definition of Done

- [ ] JWT claims extended with `email_verified`
- [ ] Registration context detection implemented
- [ ] Role assignment based on context works
- [ ] User activation with email trigger works
- [ ] Configuration extended and documented
- [ ] Code follows established patterns
- [ ] Build succeeds without errors

## Dependencies

- T56: Verification service must be available
- T57: Auth server handlers must be implemented
- Existing user, role, and project member repositories

## Risk Factors

- **Low Risk**: Simple claim addition to existing JWT structure
- **Medium Risk**: Role assignment logic must handle edge cases

## Notes

- `email_verified` in JWT allows frontend to check without API call
- Refresh token flow must update `email_verified` from database
- `autoActivate: true` (default) means users can access dashboard immediately
- `autoActivate: false` means users wait for admin approval
- Dashboard client ID comes from config for comparison
- Global 'user' role provides `dashboard:read` permission
