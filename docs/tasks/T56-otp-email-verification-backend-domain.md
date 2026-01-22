# Task T56: Backend OTP & Email Verification Domain

**Story Reference:** US14-standalone-idp-application.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T54 (Database Schema), T55 (Notification Service)

## Objective

Implement OTP and email verification logic within the `oauth_auth` domain, including repository methods, service layer with rate limiting, and token generation/validation utilities.

## Acceptance Criteria

- [ ] OTP generation (6 digits) with SHA256 hashing before storage
- [ ] OTP validation with expiry check and one-time use
- [ ] OTP rate limiting (3 per email per 15 minutes)
- [ ] Email verification token generation (32 bytes, base64url)
- [ ] Email verification token validation with expiry and one-time use
- [ ] Integration with NotificationService for sending emails
- [ ] All tokens hashed before database storage

## Technical Requirements

### OTP Repository Methods

Add to `internal/domain/oauth_auth/interface.go`:

```go
// OTP Repository methods
type OTPRepositor interface {
    CreateOTP(ctx context.Context, email, otpHash string, expiresAt time.Time) error
    GetValidOTP(ctx context.Context, email, otpHash string) (*OTPToken, error)
    MarkOTPUsed(ctx context.Context, id int64) error
    CountRecentOTPs(ctx context.Context, email string, since time.Time) (int, error)
}

// Email Verification Repository methods
type EmailVerificationRepositor interface {
    CreateVerificationToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
    GetValidToken(ctx context.Context, tokenHash string) (*EmailVerificationToken, error)
    MarkTokenUsed(ctx context.Context, id int64) error
    InvalidateUserTokens(ctx context.Context, userID int64) error
}
```

### OTP Model

Add to `internal/domain/oauth_auth/model.go`:

```go
type OTPToken struct {
    ID        int64
    Email     string
    OTPHash   string
    ExpiresAt time.Time
    UsedAt    *time.Time
    CreatedAt time.Time
}

type EmailVerificationToken struct {
    ID        int64
    UserID    int64
    TokenHash string
    ExpiresAt time.Time
    UsedAt    *time.Time
    CreatedAt time.Time
}
```

### OTP Repository Implementation

Create `internal/domain/oauth_auth/otp_repo.go`:

```go
package oauth_auth

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type OTPRepo struct {
    db *pgxpool.Pool
}

func NewOTPRepo(db *pgxpool.Pool) *OTPRepo {
    return &OTPRepo{db: db}
}

func (r *OTPRepo) CreateOTP(ctx context.Context, email, otpHash string, expiresAt time.Time) error {
    query := `
        INSERT INTO altalune_otp_tokens (email, otp_hash, expires_at)
        VALUES ($1, $2, $3)
    `
    _, err := r.db.Exec(ctx, query, email, otpHash, expiresAt)
    return err
}

func (r *OTPRepo) GetValidOTP(ctx context.Context, email, otpHash string) (*OTPToken, error) {
    query := `
        SELECT id, email, otp_hash, expires_at, used_at, created_at
        FROM altalune_otp_tokens
        WHERE email = $1 AND otp_hash = $2 AND used_at IS NULL AND expires_at > NOW()
        ORDER BY created_at DESC
        LIMIT 1
    `
    var otp OTPToken
    err := r.db.QueryRow(ctx, query, email, otpHash).Scan(
        &otp.ID, &otp.Email, &otp.OTPHash, &otp.ExpiresAt, &otp.UsedAt, &otp.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &otp, nil
}

func (r *OTPRepo) MarkOTPUsed(ctx context.Context, id int64) error {
    query := `UPDATE altalune_otp_tokens SET used_at = NOW() WHERE id = $1`
    _, err := r.db.Exec(ctx, query, id)
    return err
}

func (r *OTPRepo) CountRecentOTPs(ctx context.Context, email string, since time.Time) (int, error) {
    query := `
        SELECT COUNT(*) FROM altalune_otp_tokens
        WHERE email = $1 AND created_at > $2
    `
    var count int
    err := r.db.QueryRow(ctx, query, email, since).Scan(&count)
    return count, err
}
```

### Email Verification Repository Implementation

Create `internal/domain/oauth_auth/verification_repo.go`:

```go
package oauth_auth

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type EmailVerificationRepo struct {
    db *pgxpool.Pool
}

func NewEmailVerificationRepo(db *pgxpool.Pool) *EmailVerificationRepo {
    return &EmailVerificationRepo{db: db}
}

func (r *EmailVerificationRepo) CreateVerificationToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
    query := `
        INSERT INTO altalune_email_verification_tokens (user_id, token_hash, expires_at)
        VALUES ($1, $2, $3)
    `
    _, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
    return err
}

func (r *EmailVerificationRepo) GetValidToken(ctx context.Context, tokenHash string) (*EmailVerificationToken, error) {
    query := `
        SELECT id, user_id, token_hash, expires_at, used_at, created_at
        FROM altalune_email_verification_tokens
        WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()
    `
    var token EmailVerificationToken
    err := r.db.QueryRow(ctx, query, tokenHash).Scan(
        &token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &token, nil
}

func (r *EmailVerificationRepo) MarkTokenUsed(ctx context.Context, id int64) error {
    query := `UPDATE altalune_email_verification_tokens SET used_at = NOW() WHERE id = $1`
    _, err := r.db.Exec(ctx, query, id)
    return err
}

func (r *EmailVerificationRepo) InvalidateUserTokens(ctx context.Context, userID int64) error {
    query := `UPDATE altalune_email_verification_tokens SET used_at = NOW() WHERE user_id = $1 AND used_at IS NULL`
    _, err := r.db.Exec(ctx, query, userID)
    return err
}
```

### OTP Service

Create `internal/domain/oauth_auth/otp_service.go`:

```go
package oauth_auth

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "time"

    "altalune/internal/shared/notification"
)

var (
    ErrEmailNotRegistered = errors.New("email not registered")
    ErrOTPRateLimited     = errors.New("too many OTP requests, please try again later")
    ErrInvalidOTP         = errors.New("invalid or expired OTP")
    ErrOTPAlreadyUsed     = errors.New("OTP has already been used")
)

type OTPConfig struct {
    Expiry          time.Duration // Default: 5 minutes
    RateLimit       int           // Default: 3
    RateLimitWindow time.Duration // Default: 15 minutes
}

type OTPService struct {
    repo         OTPRepositor
    userRepo     UserLookupRepositor // Interface to check if email exists
    notification *notification.NotificationService
    cfg          OTPConfig
}

func NewOTPService(repo OTPRepositor, userRepo UserLookupRepositor, notification *notification.NotificationService, cfg OTPConfig) *OTPService {
    if cfg.Expiry == 0 {
        cfg.Expiry = 5 * time.Minute
    }
    if cfg.RateLimit == 0 {
        cfg.RateLimit = 3
    }
    if cfg.RateLimitWindow == 0 {
        cfg.RateLimitWindow = 15 * time.Minute
    }
    return &OTPService{repo: repo, userRepo: userRepo, notification: notification, cfg: cfg}
}

// GenerateAndSendOTP creates OTP, stores hash, sends email
func (s *OTPService) GenerateAndSendOTP(ctx context.Context, email string) error {
    // 1. Check if email exists
    user, err := s.userRepo.GetUserByEmail(ctx, email)
    if err != nil || user == nil {
        return ErrEmailNotRegistered
    }

    // 2. Check rate limit
    since := time.Now().Add(-s.cfg.RateLimitWindow)
    count, err := s.repo.CountRecentOTPs(ctx, email, since)
    if err != nil {
        return fmt.Errorf("failed to check rate limit: %w", err)
    }
    if count >= s.cfg.RateLimit {
        return ErrOTPRateLimited
    }

    // 3. Generate 6-digit OTP
    otp, err := generateOTP(6)
    if err != nil {
        return fmt.Errorf("failed to generate OTP: %w", err)
    }

    // 4. Hash and store
    otpHash := hashToken(otp)
    expiresAt := time.Now().Add(s.cfg.Expiry)
    if err := s.repo.CreateOTP(ctx, email, otpHash, expiresAt); err != nil {
        return fmt.Errorf("failed to store OTP: %w", err)
    }

    // 5. Send email
    userName := user.FirstName
    if userName == "" {
        userName = email
    }
    if err := s.notification.SendOTPEmail(ctx, email, otp, userName); err != nil {
        return fmt.Errorf("failed to send OTP email: %w", err)
    }

    return nil
}

// ValidateOTP checks OTP and marks as used
func (s *OTPService) ValidateOTP(ctx context.Context, email, otp string) (*UserInfo, error) {
    otpHash := hashToken(otp)

    otpToken, err := s.repo.GetValidOTP(ctx, email, otpHash)
    if err != nil {
        return nil, ErrInvalidOTP
    }

    // Mark as used
    if err := s.repo.MarkOTPUsed(ctx, otpToken.ID); err != nil {
        return nil, fmt.Errorf("failed to mark OTP as used: %w", err)
    }

    // Get user info for session creation
    user, err := s.userRepo.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func generateOTP(length int) (string, error) {
    const digits = "0123456789"
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    for i := range b {
        b[i] = digits[int(b[i])%len(digits)]
    }
    return string(b), nil
}

func hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}
```

### Email Verification Service

Create `internal/domain/oauth_auth/verification_service.go`:

```go
package oauth_auth

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "time"

    "altalune/internal/shared/notification"
)

var (
    ErrInvalidVerificationToken = errors.New("invalid or expired verification token")
    ErrTokenAlreadyUsed         = errors.New("verification token has already been used")
    ErrUserNotFound             = errors.New("user not found")
)

type VerificationConfig struct {
    TokenExpiry time.Duration // Default: 24 hours
}

type EmailVerificationService struct {
    repo         EmailVerificationRepositor
    userRepo     UserEmailVerificationRepositor
    notification *notification.NotificationService
    cfg          VerificationConfig
}

func NewEmailVerificationService(
    repo EmailVerificationRepositor,
    userRepo UserEmailVerificationRepositor,
    notification *notification.NotificationService,
    cfg VerificationConfig,
) *EmailVerificationService {
    if cfg.TokenExpiry == 0 {
        cfg.TokenExpiry = 24 * time.Hour
    }
    return &EmailVerificationService{
        repo: repo, userRepo: userRepo, notification: notification, cfg: cfg,
    }
}

// GenerateAndSendVerificationEmail creates token and sends email
func (s *EmailVerificationService) GenerateAndSendVerificationEmail(ctx context.Context, userID int64) error {
    // Get user info
    user, err := s.userRepo.GetUserByID(ctx, userID)
    if err != nil {
        return ErrUserNotFound
    }

    // Invalidate any existing tokens
    _ = s.repo.InvalidateUserTokens(ctx, userID)

    // Generate token (32 bytes = 256 bits)
    token, err := generateSecureToken(32)
    if err != nil {
        return fmt.Errorf("failed to generate token: %w", err)
    }

    // Hash and store
    tokenHash := hashToken(token)
    expiresAt := time.Now().Add(s.cfg.TokenExpiry)
    if err := s.repo.CreateVerificationToken(ctx, userID, tokenHash, expiresAt); err != nil {
        return fmt.Errorf("failed to store token: %w", err)
    }

    // Send email
    userName := user.FirstName
    if userName == "" {
        userName = user.Email
    }
    if err := s.notification.SendVerificationEmail(ctx, user.Email, token, userName); err != nil {
        return fmt.Errorf("failed to send verification email: %w", err)
    }

    return nil
}

// VerifyEmail validates token and marks user as verified
func (s *EmailVerificationService) VerifyEmail(ctx context.Context, token string) error {
    tokenHash := hashToken(token)

    // Get and validate token
    verificationToken, err := s.repo.GetValidToken(ctx, tokenHash)
    if err != nil {
        return ErrInvalidVerificationToken
    }

    // Mark token as used
    if err := s.repo.MarkTokenUsed(ctx, verificationToken.ID); err != nil {
        return fmt.Errorf("failed to mark token as used: %w", err)
    }

    // Update user's email_verified status
    if err := s.userRepo.SetEmailVerified(ctx, verificationToken.UserID, true); err != nil {
        return fmt.Errorf("failed to update email verification status: %w", err)
    }

    return nil
}

// ResendVerificationEmail for authenticated users
func (s *EmailVerificationService) ResendVerificationEmail(ctx context.Context, userID int64) error {
    return s.GenerateAndSendVerificationEmail(ctx, userID)
}

func generateSecureToken(length int) (string, error) {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}
```

### User Repository Interface Extensions

Add to `internal/domain/oauth_auth/interface.go`:

```go
// UserLookupRepositor for OTP service
type UserLookupRepositor interface {
    GetUserByEmail(ctx context.Context, email string) (*UserInfo, error)
}

// UserEmailVerificationRepositor for verification service
type UserEmailVerificationRepositor interface {
    GetUserByID(ctx context.Context, userID int64) (*UserInfo, error)
    SetEmailVerified(ctx context.Context, userID int64, verified bool) error
}

// UserInfo minimal user info for services
type UserInfo struct {
    ID        int64
    PublicID  string
    Email     string
    FirstName string
    LastName  string
    IsActive  bool
}
```

## Implementation Details

### Error Handling

All custom errors should be defined in `internal/domain/oauth_auth/errors.go`:

```go
var (
    ErrEmailNotRegistered       = errors.New("email not registered")
    ErrOTPRateLimited           = errors.New("too many OTP requests")
    ErrInvalidOTP               = errors.New("invalid or expired OTP")
    ErrInvalidVerificationToken = errors.New("invalid or expired verification token")
)
```

### Security Considerations

1. **Never log tokens or OTPs** - only log hashes or IDs
2. **Use constant-time comparison** for hash validation (already handled by DB query)
3. **Rate limiting** prevents brute force on OTPs
4. **Token invalidation** prevents reuse and allows new token generation

## Files to Create

- `internal/domain/oauth_auth/otp_repo.go`
- `internal/domain/oauth_auth/verification_repo.go`
- `internal/domain/oauth_auth/otp_service.go`
- `internal/domain/oauth_auth/verification_service.go`

## Files to Modify

- `internal/domain/oauth_auth/interface.go` - Add new interfaces
- `internal/domain/oauth_auth/model.go` - Add OTP and verification token models
- `internal/domain/oauth_auth/errors.go` - Add new error types
- `internal/domain/user/repo.go` - Add `GetUserByEmail`, `SetEmailVerified` methods

## Testing Requirements

```go
func TestOTPGeneration() {
    otp, err := generateOTP(6)
    assert.NoError(t, err)
    assert.Len(t, otp, 6)
    // Verify all characters are digits
}

func TestOTPRateLimiting() {
    // Create 3 OTPs, 4th should fail
}

func TestTokenHashing() {
    token := "test-token"
    hash1 := hashToken(token)
    hash2 := hashToken(token)
    assert.Equal(t, hash1, hash2)
    assert.Len(t, hash1, 64) // SHA256 hex
}
```

## Commands to Run

```bash
# Build to verify compilation
make build

# Run tests
go test ./internal/domain/oauth_auth/...
```

## Validation Checklist

- [ ] OTP generation produces 6 numeric digits
- [ ] OTP hash is 64 characters (SHA256 hex)
- [ ] Rate limiting blocks after 3 OTPs in 15 minutes
- [ ] Verification token is base64url encoded
- [ ] Token hashing is consistent
- [ ] Used tokens cannot be reused
- [ ] Expired tokens are rejected
- [ ] NotificationService integration works

## Definition of Done

- [ ] All repository methods implemented
- [ ] OTPService with rate limiting complete
- [ ] EmailVerificationService complete
- [ ] Error types defined and used consistently
- [ ] Token/OTP hashing implemented securely
- [ ] User repository extended with needed methods
- [ ] Code follows established patterns
- [ ] Build succeeds without errors

## Dependencies

- T54: Database tables must exist
- T55: NotificationService must be available
- Existing user domain repository

## Risk Factors

- **Low Risk**: Standard cryptographic operations using Go stdlib
- **Medium Risk**: Rate limiting logic must be tested thoroughly

## Notes

- OTP is 6 digits for usability (easy to type)
- Verification token is 32 bytes (256 bits) for security
- Both use SHA256 hashing before storage
- Rate limiting window is configurable
- Services are designed for dependency injection
- Consider adding audit logging in future
