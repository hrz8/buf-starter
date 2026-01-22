package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

// ResendEmailSender implements EmailSender using the Resend API.
type ResendEmailSender struct {
	client    *resend.Client
	fromEmail string
	fromName  string
}

// NewResendEmailSender creates a new Resend email sender.
func NewResendEmailSender(apiKey, fromEmail, fromName string) *ResendEmailSender {
	return &ResendEmailSender{
		client:    resend.NewClient(apiKey),
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// SendEmail sends an email using the Resend API.
func (r *ResendEmailSender) SendEmail(ctx context.Context, to, subject, htmlBody, textBody string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", r.fromName, r.fromEmail),
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
		Text:    textBody,
	}

	_, err := r.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}

	return nil
}
