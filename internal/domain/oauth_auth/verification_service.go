package oauth_auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/shared/notification"
)

// EmailVerificationService handles email verification token generation and validation.
type EmailVerificationService struct {
	repo         EmailVerificationRepositor
	userRepo     UserEmailVerificationRepositor
	notification *notification.NotificationService
	log          altalune.Logger
	cfg          altalune.Config
}

// NewEmailVerificationService creates a new email verification service.
func NewEmailVerificationService(
	repo EmailVerificationRepositor,
	userRepo UserEmailVerificationRepositor,
	notificationSvc *notification.NotificationService,
	log altalune.Logger,
	cfg altalune.Config,
) *EmailVerificationService {
	return &EmailVerificationService{
		repo:         repo,
		userRepo:     userRepo,
		notification: notificationSvc,
		log:          log,
		cfg:          cfg,
	}
}

// GenerateAndSendVerificationEmail creates a verification token and sends it via email.
// Invalidates any existing tokens for the user before creating a new one.
func (s *EmailVerificationService) GenerateAndSendVerificationEmail(ctx context.Context, userID int64) error {
	// Get user info
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Error("user not found for verification email", "error", err, "userID", userID)
		return ErrUserNotFound
	}

	// Check if already verified
	if user.EmailVerified {
		s.log.Debug("user already verified, skipping email", "userID", userID)
		return nil // Silently succeed - user is already verified
	}

	// Invalidate any existing tokens (ignore errors)
	if err := s.repo.InvalidateUserTokens(ctx, userID); err != nil {
		s.log.Warn("failed to invalidate existing tokens", "error", err, "userID", userID)
	}

	// Generate secure token (32 bytes = 256 bits of entropy)
	token, err := generateSecureToken(32)
	if err != nil {
		s.log.Error("failed to generate verification token", "error", err)
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Hash and store (never log the actual token)
	tokenHash := hashToken(token)
	tokenExpiry := time.Duration(s.cfg.GetVerificationTokenExpiryHours()) * time.Hour
	expiresAt := time.Now().Add(tokenExpiry)
	if err := s.repo.CreateVerificationToken(ctx, userID, tokenHash, expiresAt); err != nil {
		s.log.Error("failed to store verification token", "error", err, "userID", userID)
		return fmt.Errorf("failed to store token: %w", err)
	}

	// Send email
	userName := user.FirstName
	if userName == "" {
		userName = user.Email
	}
	if err := s.notification.SendVerificationEmail(ctx, user.Email, token, userName); err != nil {
		s.log.Error("failed to send verification email", "error", err, "userID", userID)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	s.log.Info("verification email sent", "userID", userID, "email", user.Email)
	return nil
}

// VerifyEmail validates the token and marks the user's email as verified.
func (s *EmailVerificationService) VerifyEmail(ctx context.Context, token string) error {
	tokenHash := hashToken(token)

	// Get and validate token
	verificationToken, err := s.repo.GetValidToken(ctx, tokenHash)
	if err != nil {
		s.log.Debug("invalid verification token attempt")
		return ErrInvalidVerificationToken
	}

	// Mark token as used first (atomic operation)
	if err := s.repo.MarkTokenUsed(ctx, verificationToken.ID); err != nil {
		s.log.Error("failed to mark verification token as used", "error", err, "tokenID", verificationToken.ID)
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Update user's email_verified status
	if err := s.userRepo.SetEmailVerified(ctx, verificationToken.UserID, true); err != nil {
		s.log.Error("failed to set email verified", "error", err, "userID", verificationToken.UserID)
		return fmt.Errorf("failed to update verification status: %w", err)
	}

	s.log.Info("email verified successfully", "userID", verificationToken.UserID)
	return nil
}

// ResendVerificationEmail sends a new verification email to the user.
// Alias for GenerateAndSendVerificationEmail for clarity in the API.
func (s *EmailVerificationService) ResendVerificationEmail(ctx context.Context, userID int64) error {
	return s.GenerateAndSendVerificationEmail(ctx, userID)
}

// generateSecureToken generates a cryptographically secure random token.
// Returns base64url-encoded string (URL-safe without padding).
func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
