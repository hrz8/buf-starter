package notification

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/hrz8/altalune/internal/shared/notification/email"
)

// NotificationService handles sending notifications (emails, etc.)
type NotificationService struct {
	emailSender email.EmailSender
	templates   *template.Template
	baseURL     string // Auth server URL for verification links
}

// VerificationEmailData contains data for verification email templates.
type VerificationEmailData struct {
	UserName         string
	VerificationLink string
	ExpiryHours      int
}

// OTPEmailData contains data for OTP email templates.
type OTPEmailData struct {
	UserName      string
	OTPCode       string
	ExpiryMinutes int
}

// NewNotificationService creates a new notification service with embedded templates.
func NewNotificationService(sender email.EmailSender, baseURL string) (*NotificationService, error) {
	// Parse HTML templates
	htmlTmpl, err := template.ParseFS(email.Templates, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML templates: %w", err)
	}

	// Parse text templates
	txtTmpl, err := template.ParseFS(email.Templates, "templates/*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to parse text templates: %w", err)
	}

	// Merge text templates into the main template set
	for _, t := range txtTmpl.Templates() {
		if _, err := htmlTmpl.AddParseTree(t.Name(), t.Tree); err != nil {
			return nil, fmt.Errorf("failed to merge template %s: %w", t.Name(), err)
		}
	}

	return &NotificationService{
		emailSender: sender,
		templates:   htmlTmpl,
		baseURL:     baseURL,
	}, nil
}

// SendVerificationEmail sends an email verification link to the user.
func (n *NotificationService) SendVerificationEmail(ctx context.Context, toEmail, token, userName string) error {
	data := VerificationEmailData{
		UserName:         userName,
		VerificationLink: fmt.Sprintf("%s/verify-email?token=%s", n.baseURL, token),
		ExpiryHours:      24,
	}

	htmlBody, textBody, err := n.renderTemplates("verification", data)
	if err != nil {
		return fmt.Errorf("failed to render verification templates: %w", err)
	}

	if err := n.emailSender.SendEmail(ctx, toEmail, "Verify your email address", htmlBody, textBody); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// SendOTPEmail sends a one-time password code to the user.
func (n *NotificationService) SendOTPEmail(ctx context.Context, toEmail, otp, userName string) error {
	data := OTPEmailData{
		UserName:      userName,
		OTPCode:       otp,
		ExpiryMinutes: 5,
	}

	htmlBody, textBody, err := n.renderTemplates("otp", data)
	if err != nil {
		return fmt.Errorf("failed to render OTP templates: %w", err)
	}

	if err := n.emailSender.SendEmail(ctx, toEmail, "Your login code", htmlBody, textBody); err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}

// renderTemplates renders both HTML and text versions of a template.
func (n *NotificationService) renderTemplates(name string, data any) (string, string, error) {
	var htmlBuf, textBuf bytes.Buffer

	// Template names from embed.FS use just the filename
	htmlTemplateName := name + ".html"
	if err := n.templates.ExecuteTemplate(&htmlBuf, htmlTemplateName, data); err != nil {
		return "", "", fmt.Errorf("failed to render HTML template %s: %w", htmlTemplateName, err)
	}

	textTemplateName := name + ".txt"
	if err := n.templates.ExecuteTemplate(&textBuf, textTemplateName, data); err != nil {
		return "", "", fmt.Errorf("failed to render text template %s: %w", textTemplateName, err)
	}

	return htmlBuf.String(), textBuf.String(), nil
}
