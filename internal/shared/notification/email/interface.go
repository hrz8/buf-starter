package email

import "context"

// EmailSender defines the interface for sending emails.
// This abstraction allows for different email providers (Resend, SES, etc.)
type EmailSender interface {
	// SendEmail sends an email with both HTML and plain text bodies.
	// The implementation should handle provider-specific details.
	SendEmail(ctx context.Context, to, subject, htmlBody, textBody string) error
}
