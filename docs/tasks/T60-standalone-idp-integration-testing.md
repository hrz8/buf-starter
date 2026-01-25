# Task T60: Integration Testing & Configuration

**Story Reference:** US14-standalone-idp-application.md
**Type:** Integration
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** T54-T59 (All previous tasks)

## Objective

Verify all US14 features work together end-to-end, update configuration files with proper defaults, and document the complete standalone IDP functionality.

## Acceptance Criteria

- [ ] Complete OTP login flow works end-to-end
- [ ] Email verification flow works end-to-end
- [ ] Dashboard overlay appears and disappears correctly
- [ ] Role assignment works based on registration context
- [ ] JWT contains correct `email_verified` claim
- [ ] Configuration files updated with all new settings
- [ ] All error cases handled gracefully

## Technical Requirements

### Configuration Updates

**config.yaml additions:**

```yaml
auth:
  host: "localhost"
  port: 3300
  sessionSecret: "your-session-secret-min-32-chars"
  codeExpiry: 600
  accessTokenExpiry: 3600
  refreshTokenExpiry: 2592000
  # NEW settings for US14
  autoActivate: true                    # true = users active immediately, false = require admin approval
  emailVerificationExpiry: 86400        # 24 hours in seconds
  otpExpiry: 300                        # 5 minutes in seconds
  otpRateLimit: 3                       # max OTPs per email per window
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
```

**config.example.yaml:**

```yaml
auth:
  host: "localhost"
  port: 3300
  sessionSecret: "change-this-to-a-secure-secret-min-32-chars"
  codeExpiry: 600
  accessTokenExpiry: 3600
  refreshTokenExpiry: 2592000
  autoActivate: true
  emailVerificationExpiry: 86400
  otpExpiry: 300
  otpRateLimit: 3
  otpRateLimitWindow: 900

notification:
  email:
    provider: "resend"
    resend:
      apiKey: "YOUR_RESEND_API_KEY"
      fromEmail: "noreply@yourdomain.com"
      fromName: "Your App Name"
    ses:
      region: "us-east-1"
      fromEmail: "noreply@yourdomain.com"
```

### Integration Test Scenarios

#### Scenario 1: Complete OTP Login Flow

```bash
# 1. Navigate to auth server
curl -I http://localhost:3300/
# Expected: 302 redirect to /login

# 2. Access login page (no client_id)
curl http://localhost:3300/login
# Expected: 200 with login page including "Login with Email" option

# 3. Request OTP
curl -X POST http://localhost:3300/login/email \
  -d "email=test@example.com"
# Expected: 302 redirect to /login/otp
# Expected: Email sent to test@example.com

# 4. Verify OTP (from email)
curl -X POST http://localhost:3300/login/otp/verify \
  -d "otp=123456" \
  --cookie "session=..."
# Expected: 302 redirect to /profile (session created)
```

#### Scenario 2: Email Verification Flow

```bash
# 1. Register new user via OAuth
# (User auto-activated if autoActivate=true)
# Expected: Verification email sent

# 2. Click verification link
curl http://localhost:3300/verify-email?token=abc123...
# Expected: Success page, user.email_verified = true

# 3. Verify JWT contains email_verified=true
# Decode JWT and check claim
```

#### Scenario 3: Dashboard Overlay Flow

```bash
# 1. Log in with unverified user
# 2. Access dashboard
# Expected: Overlay appears

# 3. Click "I already verified"
# Expected: Token refresh, page reload

# 4. If verified, overlay disappears
# If not verified, overlay remains
```

#### Scenario 4: Role Assignment

```bash
# Test 1: Standalone registration (no client_id)
# Expected: user assigned to project 1 with role "user"

# Test 2: Dashboard OAuth registration
# Expected: user assigned to project 1 with role "member"

# Test 3: Custom client registration
# Expected: user assigned to client's project with role "user"
```

#### Scenario 5: Admin Activation (autoActivate=false)

```bash
# 1. Configure autoActivate: false
# 2. Register new user
# Expected: User created with is_active=false
# Expected: Redirect to /pending-activation
# Expected: No verification email sent yet

# 3. Admin activates user
# Expected: is_active=true, activated_at set
# Expected: Verification email sent
# Expected: User can now log in
```

### Manual Testing Checklist

**Auth Server Testing:**
- [ ] `http://localhost:3300/` redirects to `/login` when not authenticated
- [ ] `http://localhost:3300/` redirects to `/profile` when authenticated
- [ ] Login page displays without `client_id` parameter
- [ ] "Login with Email" button visible on login page
- [ ] Email input page (`/login/email`) displays correctly
- [ ] Invalid email shows "Email not registered" error
- [ ] Rate limit triggers after 3 OTP requests in 15 minutes
- [ ] OTP page shows masked email
- [ ] OTP countdown timer works
- [ ] Invalid OTP shows error
- [ ] Valid OTP creates session and redirects
- [ ] Inactive user redirected to pending activation page
- [ ] Verification link with valid token shows success
- [ ] Verification link with invalid/expired token shows error

**Dashboard Testing:**
- [ ] Login with unverified user shows overlay
- [ ] Overlay cannot be closed by clicking outside
- [ ] Overlay cannot be closed by pressing Escape
- [ ] "Resend verification email" button works
- [ ] "I already verified" refreshes token and page
- [ ] After verification, overlay disappears
- [ ] Overlay appears on navigation to other pages

**JWT Testing:**
- [ ] Access token contains `email_verified` claim
- [ ] Claim is boolean type (not string)
- [ ] Refresh token returns updated `email_verified` status

**Role Testing:**
- [ ] Standalone user gets project role "user"
- [ ] Dashboard OAuth user gets project role "member"
- [ ] All users get global role "user"

### Error Handling Verification

| Error Case | Expected Behavior |
|------------|------------------|
| Unregistered email for OTP | Show "email not registered" error |
| Rate limit exceeded | Show "try again later" error |
| Invalid OTP | Show "invalid or expired code" error |
| Expired verification token | Show "link expired" error |
| Used verification token | Show "link already used" error |
| Network error on resend | Show generic error message |

## Files to Create

- None (integration testing task)

## Files to Modify

- `config.yaml` - Add new auth and notification settings
- `config.example.yaml` - Add new settings with placeholder values

## Commands to Run

```bash
# Build all components
make build

# Run database migrations
./bin/app migrate -c config.yaml

# Start auth server
./bin/app serve-auth -c config.yaml &

# Start main server
./bin/app serve -c config.yaml &

# Start frontend
cd frontend && pnpm dev &

# Run integration tests (if automated)
go test ./tests/integration/...
```

## Validation Checklist

- [ ] All configuration options documented
- [ ] OTP login flow works end-to-end
- [ ] Email verification flow works end-to-end
- [ ] Dashboard overlay works correctly
- [ ] Role assignment works for all contexts
- [ ] Error messages are user-friendly
- [ ] Rate limiting works correctly
- [ ] Token refresh updates verification status

## Definition of Done

- [ ] All integration scenarios pass
- [ ] Configuration files updated
- [ ] Manual testing completed
- [ ] Error cases verified
- [ ] Documentation updated if needed
- [ ] No console errors in frontend
- [ ] No error logs in backend (except expected validation errors)

## Dependencies

- T54: Database schema must be applied
- T55: Notification service must be configured
- T56: Backend services must be implemented
- T57: Auth server pages must be working
- T58: JWT claims must be extended
- T59: Dashboard overlay must be implemented

## Risk Factors

- **Low Risk**: Testing existing implementations
- **Medium Risk**: Integration issues between components

## Notes

- Test with both `autoActivate: true` and `autoActivate: false`
- Test with Resend in sandbox mode first
- Ensure email templates render correctly
- Check for race conditions in token refresh
- Monitor rate limiting behavior
- Verify session handling across auth server and dashboard
- Test with different browsers for frontend overlay

### Email States Reference

| is_active | email_verified | State | Access |
|-----------|----------------|-------|--------|
| false | false | Pending admin approval | None |
| true | false | Active, unverified | Dashboard with overlay |
| true | true | Fully active | Full access |

### Role Assignment Reference

| Registration Method | Project | Project Role | Global Role |
|---------------------|---------|--------------|-------------|
| Standalone IDP | 1 | user | user |
| Dashboard OAuth | 1 | member | user |
| Custom OAuth | Client's project | user | user |
| Admin creates | Selected | member | user |
