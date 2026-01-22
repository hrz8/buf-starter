package oauth_auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/shared/notification"
)

// OTPService handles OTP generation, validation, and email sending.
type OTPService struct {
	repo         OTPRepositor
	userRepo     UserLookupRepositor
	notification *notification.NotificationService
	log          altalune.Logger
	cfg          altalune.Config
}

// NewOTPService creates a new OTP service with the given dependencies.
func NewOTPService(
	repo OTPRepositor,
	userRepo UserLookupRepositor,
	notificationSvc *notification.NotificationService,
	log altalune.Logger,
	cfg altalune.Config,
) *OTPService {
	return &OTPService{
		repo:         repo,
		userRepo:     userRepo,
		notification: notificationSvc,
		log:          log,
		cfg:          cfg,
	}
}

// GenerateAndSendOTP creates an OTP, stores its hash, and sends it via email.
// Returns ErrEmailNotRegistered if the email is not in the system.
// Returns ErrOTPRateLimited if too many OTPs have been requested recently.
func (s *OTPService) GenerateAndSendOTP(ctx context.Context, email string) error {
	// 1. Check if email exists and get user info
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Debug("OTP request for unknown email", "email", email)
		return ErrEmailNotRegistered
	}

	// 2. Check rate limit
	rateLimitWindow := time.Duration(s.cfg.GetOTPRateLimitWindowMins()) * time.Minute
	since := time.Now().Add(-rateLimitWindow)
	count, err := s.repo.CountRecentOTPs(ctx, email, since)
	if err != nil {
		s.log.Error("failed to check OTP rate limit", "error", err, "email", email)
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
	rateLimit := s.cfg.GetOTPRateLimit()
	if count >= rateLimit {
		s.log.Warn("OTP rate limit exceeded", "email", email, "count", count, "limit", rateLimit)
		return ErrOTPRateLimited
	}

	// 3. Generate 6-digit OTP
	otp, err := generateOTP(6)
	if err != nil {
		s.log.Error("failed to generate OTP", "error", err)
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// 4. Hash and store (never log the actual OTP)
	otpHash := hashToken(otp)
	expiry := time.Duration(s.cfg.GetOTPExpirySeconds()) * time.Second
	expiresAt := time.Now().Add(expiry)
	if err := s.repo.CreateOTP(ctx, email, otpHash, expiresAt); err != nil {
		s.log.Error("failed to store OTP", "error", err, "email", email)
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	// 5. Send email
	userName := user.FirstName
	if userName == "" {
		userName = user.Email
	}
	if err := s.notification.SendOTPEmail(ctx, email, otp, userName); err != nil {
		s.log.Error("failed to send OTP email", "error", err, "email", email)
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	s.log.Info("OTP sent successfully", "email", email)
	return nil
}

// ValidateOTP checks if the provided OTP is valid for the email and marks it as used.
// Returns the user info on success for session creation.
func (s *OTPService) ValidateOTP(ctx context.Context, email, otp string) (*UserInfo, error) {
	otpHash := hashToken(otp)

	// Get and validate OTP
	otpToken, err := s.repo.GetValidOTP(ctx, email, otpHash)
	if err != nil {
		s.log.Debug("invalid OTP attempt", "email", email)
		return nil, ErrInvalidOTP
	}

	// Mark as used (atomic operation)
	if err := s.repo.MarkOTPUsed(ctx, otpToken.ID); err != nil {
		s.log.Error("failed to mark OTP as used", "error", err, "otpID", otpToken.ID)
		return nil, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	// Get user info for session creation
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Error("user not found after OTP validation", "error", err, "email", email)
		return nil, err
	}

	s.log.Info("OTP validated successfully", "email", email, "userID", user.ID)
	return user, nil
}

// generateOTP generates a cryptographically secure random numeric OTP.
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

// hashToken creates a SHA256 hash of a token and returns it as a hex string.
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
