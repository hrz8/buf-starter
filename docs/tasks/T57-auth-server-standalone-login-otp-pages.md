# Task T57: Auth Server Standalone Login & OTP Pages

**Story Reference:** US14-standalone-idp-application.md
**Type:** Backend/Templates
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T56 (OTP & Email Verification Domain)

## Objective

Implement standalone login support (no client_id required), OTP login flow pages, email verification handling, and pending activation page using Go HTML templates on the auth server.

## Acceptance Criteria

- [ ] Root `/` redirects to `/login` or `/profile` based on auth state
- [ ] `/login` works without `client_id` parameter
- [ ] "Login with Email" option on login page
- [ ] `/login/email` shows email input form
- [ ] `/login/otp` shows OTP input form with countdown timer
- [ ] `/verify-email` handles verification token validation
- [ ] `/pending-activation` shows waiting page for admin approval
- [ ] All pages styled consistently with existing auth server templates

## Technical Requirements

### New Routes

Add to auth server routes:

```go
// Standalone IDP routes
mux.HandleFunc("GET /", handler.HandleRoot)
mux.HandleFunc("GET /login/email", handler.HandleEmailLoginPage)
mux.HandleFunc("POST /login/email", handler.HandleEmailLoginSubmit)
mux.HandleFunc("GET /login/otp", handler.HandleOTPPage)
mux.HandleFunc("POST /login/otp/verify", handler.HandleOTPVerify)
mux.HandleFunc("GET /verify-email", handler.HandleVerifyEmail)
mux.HandleFunc("POST /resend-verification", handler.HandleResendVerification)
mux.HandleFunc("GET /pending-activation", handler.HandlePendingActivation)
```

### Handler Implementations

Add to `internal/domain/oauth_auth/handler.go`:

```go
// HandleRoot redirects based on auth state
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
    session, _ := h.sessionStore.Get(r)
    if session.Data.UserID != 0 {
        http.Redirect(w, r, "/profile", http.StatusFound)
        return
    }
    http.Redirect(w, r, "/login", http.StatusFound)
}

// HandleLoginPage - modify existing to work without client_id
func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
    clientID := r.URL.Query().Get("client_id")

    // Get providers (existing logic)
    providers, _ := h.oauthProviderRepo.GetEnabledProviders(r.Context())

    data := LoginPageData{
        Providers:      providers,
        ClientID:       clientID, // May be empty for standalone
        ShowEmailLogin: true,     // Always show for standalone IDP
        IsStandalone:   clientID == "",
    }

    h.templates.Render(w, "login.html", data)
}

// HandleEmailLoginPage shows email input form
func (h *Handler) HandleEmailLoginPage(w http.ResponseWriter, r *http.Request) {
    data := EmailLoginPageData{
        Error: r.URL.Query().Get("error"),
    }
    h.templates.Render(w, "email_input.html", data)
}

// HandleEmailLoginSubmit processes email and sends OTP
func (h *Handler) HandleEmailLoginSubmit(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Redirect(w, r, "/login/email?error=invalid_request", http.StatusFound)
        return
    }

    email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
    if email == "" {
        http.Redirect(w, r, "/login/email?error=email_required", http.StatusFound)
        return
    }

    // Generate and send OTP
    err := h.otpService.GenerateAndSendOTP(r.Context(), email)
    if err != nil {
        switch {
        case errors.Is(err, ErrEmailNotRegistered):
            http.Redirect(w, r, "/login/email?error=email_not_registered", http.StatusFound)
        case errors.Is(err, ErrOTPRateLimited):
            http.Redirect(w, r, "/login/email?error=rate_limited", http.StatusFound)
        default:
            h.log.Error("failed to send OTP", "error", err)
            http.Redirect(w, r, "/login/email?error=server_error", http.StatusFound)
        }
        return
    }

    // Store email in session for OTP verification
    session, _ := h.sessionStore.Get(r)
    session.Data.PendingOTPEmail = email
    h.sessionStore.Save(r, w, session)

    http.Redirect(w, r, "/login/otp", http.StatusFound)
}

// HandleOTPPage shows OTP input form
func (h *Handler) HandleOTPPage(w http.ResponseWriter, r *http.Request) {
    session, _ := h.sessionStore.Get(r)
    email := session.Data.PendingOTPEmail

    if email == "" {
        http.Redirect(w, r, "/login/email", http.StatusFound)
        return
    }

    // Mask email for display (j***@example.com)
    maskedEmail := maskEmail(email)

    data := OTPPageData{
        Email:       maskedEmail,
        Error:       r.URL.Query().Get("error"),
        ExpiryMins:  5,
    }
    h.templates.Render(w, "otp_input.html", data)
}

// HandleOTPVerify validates OTP and creates session
func (h *Handler) HandleOTPVerify(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Redirect(w, r, "/login/otp?error=invalid_request", http.StatusFound)
        return
    }

    session, _ := h.sessionStore.Get(r)
    email := session.Data.PendingOTPEmail
    otp := r.FormValue("otp")

    if email == "" {
        http.Redirect(w, r, "/login/email", http.StatusFound)
        return
    }

    // Validate OTP
    user, err := h.otpService.ValidateOTP(r.Context(), email, otp)
    if err != nil {
        http.Redirect(w, r, "/login/otp?error=invalid_otp", http.StatusFound)
        return
    }

    // Check if user is active
    if !user.IsActive {
        http.Redirect(w, r, "/pending-activation", http.StatusFound)
        return
    }

    // Create session
    session.Data.UserID = user.ID
    session.Data.AuthenticatedAt = time.Now()
    session.Data.PendingOTPEmail = "" // Clear
    h.sessionStore.Save(r, w, session)

    // Redirect to original URL or profile
    redirectURL := session.Data.OriginalURL
    if redirectURL == "" {
        redirectURL = "/profile"
    }
    session.Data.OriginalURL = ""
    h.sessionStore.Save(r, w, session)

    http.Redirect(w, r, redirectURL, http.StatusFound)
}

// HandleVerifyEmail handles verification link clicks
func (h *Handler) HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    if token == "" {
        h.templates.Render(w, "verify_email_result.html", VerifyEmailResultData{
            Success: false,
            Error:   "missing_token",
        })
        return
    }

    err := h.verificationService.VerifyEmail(r.Context(), token)
    if err != nil {
        errorCode := "invalid_token"
        if errors.Is(err, ErrInvalidVerificationToken) {
            errorCode = "expired_or_used"
        }
        h.templates.Render(w, "verify_email_result.html", VerifyEmailResultData{
            Success: false,
            Error:   errorCode,
        })
        return
    }

    h.templates.Render(w, "verify_email_result.html", VerifyEmailResultData{
        Success: true,
    })
}

// HandleResendVerification resends verification email
func (h *Handler) HandleResendVerification(w http.ResponseWriter, r *http.Request) {
    session, _ := h.sessionStore.Get(r)
    if session.Data.UserID == 0 {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    err := h.verificationService.ResendVerificationEmail(r.Context(), session.Data.UserID)
    if err != nil {
        http.Error(w, "Failed to resend", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"success": true}`))
}

// HandlePendingActivation shows waiting page
func (h *Handler) HandlePendingActivation(w http.ResponseWriter, r *http.Request) {
    h.templates.Render(w, "pending_activation.html", nil)
}

// Helper to mask email
func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return email
    }
    local := parts[0]
    if len(local) <= 2 {
        return local[0:1] + "***@" + parts[1]
    }
    return local[0:1] + "***" + local[len(local)-1:] + "@" + parts[1]
}
```

### Session Data Extension

Add to `internal/session/model.go`:

```go
type Data struct {
    UserID          int64
    AuthenticatedAt time.Time
    OAuthState      string
    OriginalURL     string
    CSRFToken       string
    PendingOTPEmail string // NEW: Email waiting for OTP verification
}
```

### Template Data Structures

Add to `internal/authserver/templates/data.go`:

```go
type LoginPageData struct {
    Providers      []OAuthProvider
    ClientID       string
    Error          string
    ShowEmailLogin bool
    IsStandalone   bool
}

type EmailLoginPageData struct {
    Error string
}

type OTPPageData struct {
    Email      string // Masked email
    Error      string
    ExpiryMins int
}

type VerifyEmailResultData struct {
    Success bool
    Error   string
}
```

### HTML Templates

**email_input.html:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login with Email - Altalune</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body>
    <div class="auth-container">
        <div class="auth-card">
            <h1>Login with Email</h1>
            <p class="subtitle">Enter your email to receive a login code</p>

            {{if .Error}}
            <div class="error-message">
                {{if eq .Error "email_required"}}Please enter your email address{{end}}
                {{if eq .Error "email_not_registered"}}This email is not registered{{end}}
                {{if eq .Error "rate_limited"}}Too many requests. Please try again in a few minutes{{end}}
                {{if eq .Error "server_error"}}Something went wrong. Please try again{{end}}
            </div>
            {{end}}

            <form method="POST" action="/login/email" class="login-form">
                <div class="form-group">
                    <label for="email">Email address</label>
                    <input type="email" id="email" name="email" required
                           placeholder="you@example.com" autocomplete="email">
                </div>
                <button type="submit" class="btn btn-primary">Send Login Code</button>
            </form>

            <div class="auth-links">
                <a href="/login">← Back to login options</a>
            </div>
        </div>
    </div>
</body>
</html>
```

**otp_input.html:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Enter Code - Altalune</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body>
    <div class="auth-container">
        <div class="auth-card">
            <h1>Check your email</h1>
            <p class="subtitle">We sent a 6-digit code to <strong>{{.Email}}</strong></p>

            {{if .Error}}
            <div class="error-message">
                {{if eq .Error "invalid_otp"}}Invalid or expired code. Please try again{{end}}
                {{if eq .Error "invalid_request"}}Something went wrong. Please try again{{end}}
            </div>
            {{end}}

            <form method="POST" action="/login/otp/verify" class="login-form">
                <div class="form-group">
                    <label for="otp">6-digit code</label>
                    <input type="text" id="otp" name="otp" required
                           pattern="[0-9]{6}" maxlength="6" inputmode="numeric"
                           placeholder="000000" autocomplete="one-time-code"
                           class="otp-input">
                </div>
                <p class="expiry-notice">Code expires in <span id="countdown">{{.ExpiryMins}}:00</span></p>
                <button type="submit" class="btn btn-primary">Verify Code</button>
            </form>

            <div class="auth-links">
                <a href="/login/email">Didn't receive the code? Send again</a>
            </div>
        </div>
    </div>

    <script>
    // Countdown timer
    let seconds = {{.ExpiryMins}} * 60;
    const countdown = document.getElementById('countdown');
    const timer = setInterval(() => {
        seconds--;
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        countdown.textContent = mins + ':' + (secs < 10 ? '0' : '') + secs;
        if (seconds <= 0) {
            clearInterval(timer);
            countdown.textContent = 'Expired';
        }
    }, 1000);

    // Auto-submit when 6 digits entered
    document.getElementById('otp').addEventListener('input', function(e) {
        if (e.target.value.length === 6) {
            e.target.form.submit();
        }
    });
    </script>
</body>
</html>
```

**verify_email_result.html:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification - Altalune</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body>
    <div class="auth-container">
        <div class="auth-card">
            {{if .Success}}
            <div class="success-icon">✓</div>
            <h1>Email Verified!</h1>
            <p class="subtitle">Your email has been verified successfully.</p>
            <a href="/login" class="btn btn-primary">Continue to Login</a>
            {{else}}
            <div class="error-icon">✕</div>
            <h1>Verification Failed</h1>
            <p class="subtitle">
                {{if eq .Error "missing_token"}}The verification link is invalid.{{end}}
                {{if eq .Error "expired_or_used"}}This verification link has expired or has already been used.{{end}}
                {{if eq .Error "invalid_token"}}This verification link is invalid.{{end}}
            </p>
            <p>Please request a new verification email from your dashboard.</p>
            <a href="/login" class="btn btn-secondary">Go to Login</a>
            {{end}}
        </div>
    </div>
</body>
</html>
```

**pending_activation.html:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Account Pending - Altalune</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body>
    <div class="auth-container">
        <div class="auth-card">
            <div class="pending-icon">⏳</div>
            <h1>Account Pending Approval</h1>
            <p class="subtitle">Your account has been registered but requires administrator approval before you can sign in.</p>
            <p>You will receive an email once your account has been activated.</p>
            <div class="auth-links">
                <a href="/login">← Back to login</a>
            </div>
        </div>
    </div>
</body>
</html>
```

### Update Login Template

Modify existing `login.html` to add email login option:

```html
<!-- Add after OAuth provider buttons -->
{{if .ShowEmailLogin}}
<div class="divider">
    <span>or</span>
</div>
<a href="/login/email" class="btn btn-secondary btn-email">
    <svg class="icon" viewBox="0 0 24 24">
        <path d="M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"/>
    </svg>
    Login with Email
</a>
{{end}}
```

## Files to Create

- `internal/authserver/templates/email_input.html`
- `internal/authserver/templates/otp_input.html`
- `internal/authserver/templates/verify_email_result.html`
- `internal/authserver/templates/pending_activation.html`

## Files to Modify

- `internal/domain/oauth_auth/handler.go` - Add new handlers
- `internal/authserver/routes.go` - Add new routes
- `internal/authserver/templates/login.html` - Add email login option
- `internal/authserver/templates/data.go` - Add new data structures
- `internal/session/model.go` - Add PendingOTPEmail field

## Testing Requirements

**Manual Testing:**
1. Navigate to `http://localhost:3300/` - should redirect
2. Navigate to `http://localhost:3300/login` - should work without client_id
3. Click "Login with Email" - should show email form
4. Enter valid email - should redirect to OTP page
5. Enter valid OTP - should create session and redirect
6. Test email verification link - should show success/error
7. Test with inactive user - should show pending page

## Commands to Run

```bash
# Build to verify compilation
make build

# Start auth server
./bin/app serve-auth -c config.yaml

# Test in browser
open http://localhost:3300/login
```

## Validation Checklist

- [ ] Root `/` redirects correctly
- [ ] Login page works without client_id
- [ ] Email login form displays correctly
- [ ] OTP page shows masked email
- [ ] OTP countdown timer works
- [ ] OTP auto-submits on 6 digits
- [ ] Invalid OTP shows error
- [ ] Verification success page displays
- [ ] Verification error page displays
- [ ] Pending activation page displays
- [ ] All pages match existing design

## Definition of Done

- [ ] All new routes implemented
- [ ] All new templates created
- [ ] Session extended for OTP flow
- [ ] Error handling covers all cases
- [ ] Templates styled consistently
- [ ] Code follows established patterns
- [ ] Build succeeds without errors

## Dependencies

- T56: OTPService and EmailVerificationService must be available
- Existing session infrastructure
- Existing template rendering system

## Risk Factors

- **Low Risk**: Standard HTTP handler patterns
- **Medium Risk**: Session state management for OTP flow

## Notes

- Email is masked in OTP page for privacy
- Countdown timer is client-side JavaScript
- Auto-submit improves UX but form still works without JS
- All errors redirect with query params (no direct rendering on POST)
- Session stores pending OTP email to prevent tampering
