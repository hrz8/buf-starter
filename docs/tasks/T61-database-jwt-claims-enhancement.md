# Task T61: Database Migration + JWT Claims Enhancement

**Story Reference:** US15-authorization-rbac.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** None

## Objective

Seed predefined permissions with entity:action format and enhance JWT claims to include `memberships` map (project_public_id -> role). The predefined permissions migration already exists, so this task focuses on JWT memberships claim enhancement and verification.

## Acceptance Criteria

- [ ] Predefined permissions migration applied (already exists: `20260202000000_seed_predefined_permissions.sql`)
- [ ] JWT `memberships` claim populated during token generation
- [ ] Memberships map format: `{ "project_public_id": "role" }`
- [ ] Token generation fetches memberships from `altalune_project_members` table
- [ ] Refresh token returns updated memberships

## Technical Requirements

### JWT Claims Structure (Already Exists)

The `AccessTokenClaims` in `internal/shared/jwt/claims.go` already includes the `Memberships` field:

```go
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    Scope         string            `json:"scope,omitempty"`
    Email         string            `json:"email,omitempty"`
    Name          string            `json:"name,omitempty"`
    Perms         []string          `json:"perms"`
    Memberships   map[string]string `json:"memberships,omitempty"` // project_public_id -> role
    EmailVerified bool              `json:"email_verified"`
}
```

### GenerateTokenParams Extension

Modify `internal/shared/jwt/jwt.go` to add `Memberships` to params:

```go
type GenerateTokenParams struct {
    UserPublicID  string
    ClientID      string
    Scope         string
    Email         string
    Name          string
    Perms         []string
    Memberships   map[string]string  // NEW: project_public_id -> role
    EmailVerified bool
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
        Memberships:   params.Memberships, // NEW
        EmailVerified: params.EmailVerified,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    token.Header["kid"] = s.kid
    return token.SignedString(s.privateKey)
}
```

### Membership Provider Interface

Create `internal/domain/project_member/provider.go`:

```go
package project_member

import (
    "context"
)

// MembershipProvider provides user's project memberships for JWT
type MembershipProvider struct {
    repo *Repo
}

func NewMembershipProvider(repo *Repo) *MembershipProvider {
    return &MembershipProvider{repo: repo}
}

// GetUserMemberships returns map of project_public_id -> role for a user
func (p *MembershipProvider) GetUserMemberships(ctx context.Context, userID int64) (map[string]string, error) {
    query := `
        SELECT p.public_id, pm.role
        FROM altalune_project_members pm
        JOIN altalune_projects p ON pm.project_id = p.id
        WHERE pm.user_id = $1
    `

    rows, err := p.repo.db.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    memberships := make(map[string]string)
    for rows.Next() {
        var projectPublicID, role string
        if err := rows.Scan(&projectPublicID, &role); err != nil {
            return nil, err
        }
        memberships[projectPublicID] = role
    }

    return memberships, rows.Err()
}
```

### Update Token Generation in OAuth Auth Service

Modify `internal/domain/oauth_auth/service.go` to include memberships:

```go
func (s *Service) generateTokenPair(ctx context.Context, user *User, client *OAuthClient, scope string) (*TokenPair, error) {
    // Get user permissions
    perms, _ := s.permissionProvider.GetUserPermissions(ctx, user.ID)

    // Get user memberships (NEW)
    memberships, _ := s.membershipProvider.GetUserMemberships(ctx, user.ID)

    // Get email_verified status
    emailVerified, _ := s.userRepo.GetEmailVerified(ctx, user.ID)

    // Generate access token with memberships
    accessToken, err := s.jwtSigner.GenerateAccessToken(jwt.GenerateTokenParams{
        UserPublicID:  user.PublicID,
        ClientID:      client.ClientID,
        Scope:         scope,
        Email:         user.Email,
        Name:          fmt.Sprintf("%s %s", user.FirstName, user.LastName),
        Perms:         perms,
        Memberships:   memberships, // NEW
        EmailVerified: emailVerified,
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

### Predefined Permissions Verification

Verify migration `database/migrations/20260202000000_seed_predefined_permissions.sql` includes all required permissions:

| Entity | Read | Write | Delete |
|--------|------|-------|--------|
| employee | employee:read | employee:write | employee:delete |
| user | user:read | user:write | user:delete |
| role | role:read | role:write | role:delete |
| permission | permission:read | permission:write | permission:delete |
| project | project:read | project:write | project:delete |
| apikey | apikey:read | apikey:write | apikey:delete |
| chatbot | chatbot:read | chatbot:write | chatbot:delete |
| client | client:read | client:write | client:delete |
| member | member:read | member:write | member:delete |
| iam | iam:read | iam:write | - |

## Files to Create

- `internal/domain/project_member/provider.go` - Membership provider for JWT

## Files to Modify

- `internal/shared/jwt/jwt.go` - Add Memberships to GenerateTokenParams
- `internal/domain/oauth_auth/service.go` - Include memberships in token generation
- `internal/container/container.go` - Wire MembershipProvider

## Testing Requirements

```go
func TestMembershipProvider(t *testing.T) {
    // Setup test user with project memberships
    // ...

    memberships, err := provider.GetUserMemberships(ctx, userID)
    assert.NoError(t, err)
    assert.Equal(t, "admin", memberships["proj_abc123"])
    assert.Equal(t, "member", memberships["proj_xyz789"])
}

func TestJWTMembershipsClaimIncluded(t *testing.T) {
    token, _ := signer.GenerateAccessToken(jwt.GenerateTokenParams{
        UserPublicID: "user123",
        ClientID:     "client-uuid",
        Memberships:  map[string]string{"proj_abc": "admin"},
        // ...
    })

    claims, _ := signer.ValidateAccessToken(token)
    assert.Equal(t, "admin", claims.Memberships["proj_abc"])
}
```

## Commands to Run

```bash
# Apply migrations
./bin/app migrate -c config.yaml

# Verify permissions seeded
psql -c "SELECT name FROM altalune_permissions ORDER BY name;"

# Build to verify compilation
make build

# Run tests
go test ./internal/shared/jwt/...
go test ./internal/domain/project_member/...
go test ./internal/domain/oauth_auth/...
```

## Validation Checklist

- [ ] Predefined permissions exist in database
- [ ] JWT GenerateTokenParams includes Memberships field
- [ ] Token generation fetches memberships from database
- [ ] Generated JWT contains memberships claim
- [ ] Memberships format is `{ "project_public_id": "role" }`
- [ ] Refresh token returns updated memberships

## Definition of Done

- [ ] Migration verified (already applied or apply now)
- [ ] MembershipProvider implemented
- [ ] JWT params extended with Memberships
- [ ] Token generation includes memberships
- [ ] All tests pass
- [ ] Build succeeds

## Dependencies

- None (foundation task)

## Risk Factors

- **Low Risk**: Simple data fetching and JWT claim addition
- **Low Risk**: Migration already exists and is well-structured

## Notes

- The `memberships` claim allows frontend to filter visible projects without API calls
- Superadmins (with `root` permission) bypass membership checks but still get memberships for display
- Memberships are refreshed on token refresh to reflect role changes
- Consider caching memberships if performance becomes an issue

### JWT Example with Memberships

```json
{
  "iss": "http://localhost:3300",
  "sub": "usr_abc123xyz",
  "aud": ["client_dashboard"],
  "exp": 1704067200,
  "iat": 1704063600,
  "scope": "openid profile email",
  "email": "user@example.com",
  "name": "John Doe",
  "email_verified": true,
  "perms": ["employee:read", "employee:write", "dashboard:read"],
  "memberships": {
    "proj_abc123": "admin",
    "proj_xyz789": "member"
  }
}
```
