package email

import (
	"context"
	"errors"
)

// ErrSESNotImplemented is returned when SES is selected but not yet implemented.
var ErrSESNotImplemented = errors.New("AWS SES email sender not implemented")

// SESEmailSender is a stub for AWS SES email sending.
// This allows configuration of SES as a provider for future implementation.
type SESEmailSender struct {
	region    string
	fromEmail string
}

// NewSESEmailSender creates a new SES email sender stub.
func NewSESEmailSender(region, fromEmail string) *SESEmailSender {
	return &SESEmailSender{
		region:    region,
		fromEmail: fromEmail,
	}
}

// SendEmail returns an error indicating SES is not yet implemented.
func (s *SESEmailSender) SendEmail(_ context.Context, _, _, _, _ string) error {
	return ErrSESNotImplemented
}
