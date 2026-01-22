# User Story US14: OAuth Server as Standalone IDP Application

## Story Overview

**As a** system administrator and end user
**I want** the OAuth server to function as a standalone Identity Provider (IDP) with email verification, OTP login, and proper user activation flows
**So that** users can authenticate directly to the auth server without requiring a client application, and administrators can manage user activation and verification

## Acceptance Criteria

### Core Functionality

#### Standalone Login Support (No client_id Required)

- **Given** a user navigates directly to `http://localhost:3300/login`
- **When** they are not authenticated
- **Then** the login page should display without requiring `client_id` parameter
- **And** OAuth provider buttons (Google, GitHub) should work for direct login
- **And** an "Login with Email" option should be available for OTP-based login

#### Homepage Redirect

- **Given** a user navigates to `http://localhost:3300/` (root)
- **When** they are authenticated
- **Then** they should be redirected to `/profile`
- **When** they are not authenticated
- **Then** they should be redirected to `/login`

#### Email Verification System

- **Given** a new user registers via OAuth provider
- **When** `autoActivate` config is `true` (default)
- **Then** user is created with `is_active=true`, `email_verified=false`
- **And** user is assigned to default project (id=1) with appropriate role
- **And** user is assigned global role 'user' from `altalune_roles`
- **And** verification email is sent with one-time use token (24h expiry)
- **And** user can access the dashboard but sees blocking verification overlay

- **Given** a new user registers via OAuth provider
- **When** `autoActivate` config is `false`
- **Then** user is created with `is_active=false`, `email_verified=false`
- **And** user is shown "Account pending admin approval" page
- **And** NO verification email is sent yet

#### Admin User Activation

- **Given** a user exists with `is_active=false` (pending admin approval)
- **When** admin activates the user from the dashboard
- **Then** user's `is_active` is set to `true`
- **And** `activated_at` timestamp is set (first activation only)
- **And** verification email is automatically sent
- **And** user is assigned to default project and global 'user' role if not already

#### Email Verification Flow

- **Given** a user receives a verification email
- **When** they click the verification link
- **Then** the token is validated (not expired, not used)
- **And** user's `email_verified` is set to `true`
- **And** token is marked as used (one-time)
- **And** user is redirected to login page with success message

- **Given** a user clicks an invalid/expired verification link
- **Then** appropriate error message is displayed
- **And** option to request new verification email is shown

#### OTP Login Flow

- **Given** a user is on the login page
- **When** they click "Login with Email"
- **Then** an email input form is displayed

- **Given** a user enters their email for OTP login
- **When** the email is registered in the system
- **Then** OTP (6 digits, 5 min expiry) is generated and sent to email
- **And** user is redirected to OTP input page
- **And** rate limit: max 3 OTPs per email per 15 minutes

- **Given** a user enters their email for OTP login
- **When** the email is NOT registered in the system
- **Then** error message "Email not registered" is displayed
- **And** no OTP is sent

- **Given** a user is on OTP input page
- **When** they enter valid OTP within 5 minutes
- **Then** OTP is validated and marked as used
- **And** session is created
- **And** user is redirected based on context (profile or original URL)

#### Dashboard Email Verification Overlay

- **Given** a user is logged in with `email_verified=false` in JWT claims
- **When** they access any dashboard page
- **Then** a blocking full-screen overlay is displayed
- **And** overlay shows message about email verification required
- **And** overlay has "I already verified my email" button
- **And** clicking button refreshes JWT (via token refresh) and page
- **And** overlay reappears on any navigation until verified

#### Project Membership Role Assignment

- **Given** a new user registers via OAuth
- **When** registration is via standalone IDP (no client_id)
- **Then** user is assigned to project_id=1 with role `user`

- **Given** a new user registers via OAuth
- **When** registration is via dashboard OAuth client (from config.yaml `dashboardOauth`)
- **Then** user is assigned to project_id=1 with role `member`

- **Given** a new user registers via OAuth
- **When** registration is via custom OAuth client (not dashboard client)
- **Then** user is assigned to client's project with role `user`

- **Given** admin creates a user from dashboard
- **Then** user is assigned to selected project with role `member`

### Security Requirements

#### Email Verification Tokens

- Tokens must be cryptographically random (32 bytes, base64url encoded)
- Tokens must be one-time use only
- Tokens must expire after 24 hours
- Tokens must be stored hashed in database
- Expired/used tokens must return clear error messages

#### OTP Security

- OTP must be 6 digits numeric
- OTP must expire after 5 minutes
- OTP must be one-time use only
- Rate limit: 3 OTP requests per email per 15 minutes
- OTP must be stored hashed in database
- Failed attempts should be tracked (future: lockout after N failures)

#### JWT Claims Extension

- `email_verified` boolean claim must be added to access tokens
- Claim must be refreshed when user verifies email
- Dashboard must decode JWT client-side (no verification, just decode)

### Data Validation

#### Users Table Extension

- `email_verified` - Required boolean, default `false`
- `activated_at` - Optional timestamp, set on first admin activation

#### Email Verification Tokens Table

- `id` - Required, primary key
- `user_id` - Required, references altalune_users
- `token_hash` - Required, SHA256 hash of token
- `expires_at` - Required timestamp
- `used_at` - Optional timestamp (soft-delete pattern)
- `created_at` - Required timestamp

#### OTP Tokens Table

- `id` - Required, primary key
- `email` - Required, indexed for lookup
- `otp_hash` - Required, SHA256 hash of OTP
- `expires_at` - Required timestamp
- `used_at` - Optional timestamp (soft-delete pattern)
- `created_at` - Required timestamp

#### Predefined Roles Migration

- Role `user` with description "Default role for authenticated users"
- Permission `dashboard:read` with description "Basic dashboard read access"
- Role-permission mapping: `user` has `dashboard:read`

### User Experience

#### Login Page Enhancement

- Add "Login with Email" button/link below OAuth provider buttons
- Email input form with clear placeholder
- Loading states during OTP send
- Success message when OTP is sent
- Error messages for unregistered email or rate limit

#### OTP Input Page

- 6-digit input field (or 6 separate boxes)
- Auto-submit on complete entry
- Countdown timer showing OTP expiry
- "Resend OTP" link (respects rate limit)
- Clear error message on invalid OTP

#### Pending Activation Page

- Clear message: "Your account has been registered but requires admin approval"
- Contact information or instructions
- Professional design consistent with login page

#### Email Verification Overlay (Dashboard)

- Full-screen modal blocking all interaction
- Clear message explaining email verification required
- "Resend verification email" button
- "I already verified my email" button (refreshes token + page)
- Cannot be dismissed by clicking outside or pressing escape

#### Email Templates

- Verification email: Subject, greeting, verification link, expiry info
- OTP email: Subject, OTP code prominently displayed, expiry info
- Professional HTML templates with Altalune branding
- Plain text fallback versions

## Technical Requirements

### Backend Architecture

#### Database Migration

File: `database/migrations/YYYYMMDDHHMMSS_add_email_verification.sql`

```sql
-- Add email_verified and activated_at to users table
ALTER TABLE altalune_users
ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN activated_at TIMESTAMPTZ;

-- Create email verification tokens table
CREATE TABLE altalune_email_verification_tokens (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES altalune_users(id) ON DELETE CASCADE,
  token_hash VARCHAR(64) NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ux_email_verification_token UNIQUE (token_hash)
);

CREATE INDEX ix_email_verification_user_id ON altalune_email_verification_tokens(user_id);
CREATE INDEX ix_email_verification_expires ON altalune_email_verification_tokens(expires_at)
  WHERE used_at IS NULL;

-- Create OTP tokens table
CREATE TABLE altalune_otp_tokens (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  otp_hash VARCHAR(64) NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_otp_email ON altalune_otp_tokens(email);
CREATE INDEX ix_otp_expires ON altalune_otp_tokens(expires_at) WHERE used_at IS NULL;

-- Create predefined role and permission
INSERT INTO altalune_roles (public_id, name, description, created_at, updated_at)
VALUES (
  'usr' || substring(md5(random()::text) from 1 for 11),
  'user',
  'Default role for authenticated users',
  NOW(), NOW()
) ON CONFLICT (name) DO NOTHING;

INSERT INTO altalune_permissions (public_id, name, description, created_at, updated_at)
VALUES (
  'perm' || substring(md5(random()::text) from 1 for 10),
  'dashboard:read',
  'Basic dashboard read access',
  NOW(), NOW()
) ON CONFLICT (name) DO NOTHING;

-- Link role to permission
INSERT INTO altalune_roles_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM altalune_roles r, altalune_permissions p
WHERE r.name = 'user' AND p.name = 'dashboard:read'
ON CONFLICT DO NOTHING;
```

#### Notification Package

Directory: `internal/shared/notification/`

```
internal/shared/notification/
├── email/
│   ├── interface.go       # EmailSender interface
│   ├── resend.go          # Resend implementation
│   ├── ses.go             # AWS SES stub (not implemented)
│   └── templates/
│       ├── verification.html
│       ├── verification.txt
│       ├── otp.html
│       └── otp.txt
├── notification.go        # Main notification service
└── config.go              # Email config types
```

Interface:
```go
type EmailSender interface {
    SendEmail(ctx context.Context, to, subject string, htmlBody, textBody string) error
}

type NotificationService struct {
    emailSender EmailSender
    templates   *template.Template
}

func (n *NotificationService) SendVerificationEmail(ctx context.Context, email, token string) error
func (n *NotificationService) SendOTPEmail(ctx context.Context, email, otp string) error
```

#### Configuration Extension

```yaml
# config.yaml additions
auth:
  autoActivate: true                    # false = require admin approval
  emailVerificationExpiry: 86400        # 24 hours in seconds
  otpExpiry: 300                        # 5 minutes in seconds
  otpRateLimit: 3                       # max OTPs per 15 minutes
  otpRateLimitWindow: 900               # 15 minutes in seconds

notification:
  email:
    provider: "resend"                  # "resend" or "ses"
    resend:
      apiKey: "re_xxxxxxxxxxxx"
      fromEmail: "noreply@altalune.com"
      fromName: "Altalune"
    ses:
      region: "ap-southeast-1"
      fromEmail: "noreply@altalune.com"
      # Not implemented - stub only
```

#### JWT Claims Extension

```go
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    Scope         string   `json:"scope"`
    Email         string   `json:"email,omitempty"`
    Name          string   `json:"name,omitempty"`
    Perms         []string `json:"perms,omitempty"`
    EmailVerified bool     `json:"email_verified"` // NEW
}
```

#### New Endpoints (Auth Server)

```
GET  /                           - Redirect to /login or /profile
GET  /login                      - Login page (no client_id required)
POST /login/email                - Initiate OTP login (email input)
GET  /login/otp                  - OTP input page
POST /login/otp/verify           - Verify OTP
GET  /verify-email               - Email verification link handler
POST /resend-verification        - Resend verification email
GET  /pending-activation         - Pending admin approval page
```

### Frontend (Auth Server)

#### New/Updated Pages

- `/login` - Add "Login with Email" option
- `/login/otp` - OTP input page (new)
- `/pending-activation` - Pending admin approval page (new)
- `/verify-email` - Email verification result page (new)

#### Dashboard Plugin/Middleware

File: `frontend/app/plugins/email-verification.client.ts`

```typescript
// Plugin that checks email_verified claim from JWT on route change
// If not verified, shows blocking overlay with verification prompt
// Overlay has "I already verified" button that refreshes token + page
```

### API Design

#### Auth Server Endpoints

```
POST /login/email
Request: { email: string }
Response:
  Success: { success: true, message: "OTP sent" }
  Error: { error: "email_not_registered" } or { error: "rate_limit_exceeded" }

POST /login/otp/verify
Request: { email: string, otp: string }
Response:
  Success: Redirect to /profile or original URL (session created)
  Error: { error: "invalid_otp" } or { error: "otp_expired" }

GET /verify-email?token=xxx
Response:
  Success: Redirect to /login?verified=true
  Error: Render error page with reason

POST /resend-verification
Request: (session required)
Response: { success: true } or { error: "rate_limit" }
```

## Out of Scope

- Authorization validation/enforcement (will be separate user story)
- Password-based login (only OAuth and OTP)
- MFA/2FA beyond OTP (e.g., TOTP apps)
- Email change workflow
- Password reset (no passwords in system)
- Account deletion self-service
- Social account linking UI
- AWS SES implementation (stub only, mark as "not implemented")

## Dependencies

- US5: OAuth Server Foundation (existing tables)
- US7: OAuth Authorization Server (existing auth flows)
- Resend API account and API key
- Existing session management infrastructure
- Existing JWT generation utilities

## Definition of Done

### Database

- [ ] Migration adds `email_verified`, `activated_at` to users table
- [ ] Migration creates `altalune_email_verification_tokens` table
- [ ] Migration creates `altalune_otp_tokens` table
- [ ] Migration creates predefined 'user' role
- [ ] Migration creates predefined 'dashboard:read' permission
- [ ] Migration links role to permission
- [ ] Migration is reversible

### Backend - Notification Service

- [ ] Email interface created in `internal/shared/notification/email/`
- [ ] Resend implementation working
- [ ] AWS SES stub created (returns "not implemented" error)
- [ ] Email templates created (HTML + text versions)
- [ ] NotificationService created with verification and OTP methods
- [ ] Config parsing for notification settings

### Backend - Auth Server

- [ ] Root `/` redirects to `/login` or `/profile` based on auth
- [ ] Login page works without client_id
- [ ] OTP login flow implemented (email input → send OTP → verify)
- [ ] OTP rate limiting (3 per 15 min per email)
- [ ] Email verification token generation
- [ ] Email verification link handler
- [ ] Pending activation page for `autoActivate=false`
- [ ] JWT includes `email_verified` claim
- [ ] User creation respects `autoActivate` config
- [ ] Role assignment logic based on registration context

### Backend - User Domain

- [ ] ActivateUser also sends verification email
- [ ] `activated_at` set only on first activation
- [ ] Global 'user' role assignment on registration

### Frontend - Auth Server

- [ ] Login page has "Login with Email" option
- [ ] OTP input page implemented
- [ ] Pending activation page implemented
- [ ] Email verification result pages

### Frontend - Dashboard

- [ ] Email verification overlay component
- [ ] Plugin/middleware to check `email_verified` claim
- [ ] Overlay shows on every page when not verified
- [ ] "I already verified" button refreshes token and page
- [ ] Overlay cannot be dismissed

### Configuration

- [ ] config.yaml updated with new auth settings
- [ ] config.yaml updated with notification settings
- [ ] config.example.yaml updated with examples

### Testing

- [ ] OTP generation and verification tested
- [ ] Email verification flow tested
- [ ] Rate limiting tested
- [ ] Role assignment logic tested
- [ ] JWT claims verified
- [ ] Dashboard overlay behavior verified

### Documentation

- [ ] API endpoints documented
- [ ] Configuration options documented
- [ ] Email templates documented

## Notes

### Role Assignment Matrix

| Registration Method | Project Role | Global Role | Notes |
|---------------------|-------------|-------------|-------|
| Standalone IDP (no client_id) | user | user | Requires admin upgrade for dashboard |
| Dashboard OAuth client | member | user | Can access dashboard |
| Custom OAuth client | user | user | Assigned to client's project |
| Admin creates user | member | user | Full dashboard access |

### Email Verification States

| is_active | email_verified | State | Access |
|-----------|----------------|-------|--------|
| false | false | Pending admin approval | None |
| true | false | Active, unverified | Dashboard with overlay |
| true | true | Fully active | Full access |

### Config Option: `autoActivate`

- `true` (default): Users auto-activated, need email verification
- `false`: Users need admin approval first, then email verification

### OTP Flow Security

1. User enters email
2. Server validates email exists
3. Server checks rate limit (3 per 15 min)
4. Server generates 6-digit OTP
5. Server stores SHA256(OTP) with expiry
6. Server sends OTP via email
7. User enters OTP
8. Server validates: not expired, not used, hash matches
9. Server marks OTP as used
10. Server creates session

### Email Service Architecture

```
NotificationService
    ├── EmailSender (interface)
    │   ├── ResendEmailSender (implemented)
    │   └── SESEmailSender (stub - not implemented)
    └── Templates
        ├── verification.html/.txt
        └── otp.html/.txt
```

### JWT Refresh for Email Verification

When user verifies email:
1. User clicks verification link
2. Server marks email_verified=true in database
3. User's existing JWT still has email_verified=false
4. Dashboard overlay shows "I already verified" button
5. Button triggers token refresh (existing refresh token flow)
6. New JWT has email_verified=true
7. Page reloads, overlay disappears

### Related Stories

- US5: OAuth Server Foundation (tables this extends)
- US7: OAuth Authorization Server (auth flows this extends)
- US8: Dashboard OAuth Integration (dashboard this integrates with)
- Future: US15+ - Authorization enforcement

### Future Enhancements

- AWS SES full implementation
- Email templates customization UI
- OTP via SMS
- Remember device for OTP
- Account lockout after failed OTP attempts
- Audit logging for all verification events
- Admin notification when new user pending approval
