# OAuth Authorization Server

## Overview

Altalune implements a dual-server architecture:
- **Resources API** (`serve` command): Main API for projects, users, chatbots
- **Auth Server** (`serve-auth` command): OAuth/OIDC authorization server

## Three OAuth Flow Types

### 1. Standalone IDP (Email-based Authentication)

Direct email/OTP login without external providers.

**Routes:**
```
GET  /login/email              Email login form
POST /login/email              Process email login
GET  /login/otp                OTP verification form
POST /login/otp/verify         Verify OTP code
GET  /verify-email             Email verification page
POST /resend-verification      Resend verification email
GET  /pending-activation       Pending account activation
```

### 2. OAuth Client (External Provider Authentication)

Login via Google, GitHub, Microsoft, Apple.

**Routes:**
```
GET  /login                    Provider selection page
GET  /login/{provider}         Initiate OAuth flow
GET  /auth/callback            Receive OAuth callback
GET  /profile                  View linked identities
GET  /edit-profile             Edit profile form
POST /edit-profile             Update profile
POST /profile/consents/revoke  Revoke provider consent
POST /logout                   Logout
```

### 3. OAuth Provider (This App as OAuth Server)

Allow external apps to authenticate users via this system.

**Routes:**
```
GET  /oauth/authorize          Authorization endpoint
POST /oauth/authorize          Grant authorization
POST /oauth/token              Token endpoint
GET  /oauth/userinfo           UserInfo endpoint
POST /oauth/revoke             Token revocation
POST /oauth/introspect         Token introspection
GET  /.well-known/jwks.json    Public key set
GET  /.well-known/openid-configuration  OIDC discovery
```

## BFF (Backend-for-Frontend) Pattern

Frontend uses BFF endpoints for secure token management. Tokens stored in httpOnly cookies.

**BFF Endpoints (Resources API):**
```
POST /oauth/exchange     Exchange auth code for tokens (sets cookies)
POST /oauth/logout       Clear httpOnly cookies
POST /oauth/refresh      Refresh access token
GET  /oauth/me           Get current user from cookie token
```

**Security Benefits:**
- Frontend cannot access tokens directly (XSS protection)
- Tokens stored in secure httpOnly cookies
- CSRF protection via SameSite cookie attribute

## Project Membership & Auto-Registration

When users authenticate via an OAuth client:

1. User record created in `altalune_users`
2. Identity created in `altalune_user_identities` (linked to OAuth client)
3. **Project membership auto-created** with `role='user'`

**Role Hierarchy:**
| Role | Dashboard Access | Data Access | Admin Actions |
|------|-----------------|-------------|---------------|
| owner | Yes | All projects | Full (superadmin only) |
| admin | Yes | Assigned projects | Settings, members |
| member | Yes | Assigned projects | View/edit data |
| user | No | None | API only (via tokens) |

## OAuth Client Management

**Client Types:**
- **Confidential**: Server-side apps with client_secret
- **Public**: SPA/mobile apps using PKCE

**Secret Handling:**
- Secrets encrypted with AES-256-GCM in database
- Never exposed in list/get responses (`client_secret_set` boolean instead)
- Separate `RevealClientSecret` RPC for authorized access with audit logging

## Configuration

**config.yaml sections:**

```yaml
auth:
  sessionSecret: "session-encryption-key"
  jwtPrivateKeyPath: "keys/private.pem"
  jwtPublicKeyPath: "keys/public.pem"
  accessTokenExpiry: 15m
  refreshTokenExpiry: 7d

security:
  iamEncryptionKey: "32-byte-key-for-aes-256"

dashboardOauth:
  externalServer: false  # true if auth server runs separately
  baseURL: "http://localhost:8081"

notification:
  provider: "resend"  # or "ses"
  otp:
    expiry: 5m
    length: 6
  verification:
    expiry: 24h
```

## Implementing OAuth-Protected Endpoints

**Middleware Pattern:**

```go
func (s *Server) authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header or cookie
        token := extractToken(r)
        if token == "" {
            http.Error(w, "Unauthorized", 401)
            return
        }

        // Verify JWT
        claims, err := s.jwtVerifier.Verify(token)
        if err != nil {
            http.Error(w, "Invalid token", 401)
            return
        }

        // Add user context
        ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Checking Project Access:**

```go
func (s *Service) checkProjectAccess(ctx context.Context, projectID int64, minRole string) error {
    userID := getUserIDFromContext(ctx)

    member, err := s.projectMemberRepo.GetByProjectAndUser(ctx, projectID, userID)
    if err != nil {
        return ErrProjectAccessDenied
    }

    if !hasMinimumRole(member.Role, minRole) {
        return ErrInsufficientPermissions
    }

    return nil
}

func hasMinimumRole(userRole, minRole string) bool {
    hierarchy := map[string]int{
        "user":   1,
        "member": 2,
        "admin":  3,
        "owner":  4,
    }
    return hierarchy[userRole] >= hierarchy[minRole]
}
```

## User Identity Linking

Users can link multiple OAuth providers:

```go
type UserIdentity struct {
    ID                    int64
    UserID                int64
    Provider              string    // "google", "github", "system"
    ProviderUserID        string
    OAuthClientID         *uuid.UUID
    Email                 string
    FirstName             string
    LastName              string
    AvatarURL             string
    OriginOAuthClientName string    // App user registered through
    LastLoginAt           time.Time
}
```

## OAuth Consent Management

Track user consent per OAuth client:

```go
type UserConsent struct {
    ID            int64
    UserID        int64
    OAuthClientID uuid.UUID
    Scopes        []string
    GrantedAt     time.Time
    RevokedAt     *time.Time
}
```

**Revocation endpoint allows users to disconnect apps.**
