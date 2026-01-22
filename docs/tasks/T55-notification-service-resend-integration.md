# Task T55: Notification Service - Resend Integration

**Story Reference:** US14-standalone-idp-application.md
**Type:** Backend Foundation
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** None (can run parallel with T54)

## Objective

Create a notification service with email sending capabilities using Resend as the primary provider, with an AWS SES stub for future expansion. Include HTML and text email templates for verification and OTP emails.

## Acceptance Criteria

- [ ] `EmailSender` interface defined for provider abstraction
- [ ] Resend implementation working with API key
- [ ] AWS SES stub created (returns "not implemented" error)
- [ ] Email templates created (verification + OTP, HTML + text versions)
- [ ] `NotificationService` created with `SendVerificationEmail` and `SendOTPEmail` methods
- [ ] Configuration parsing for notification settings
- [ ] Templates render correctly with provided data

## Technical Requirements

### Package Structure

```
internal/shared/notification/
├── email/
│   ├── interface.go       # EmailSender interface
│   ├── resend.go          # Resend implementation
│   ├── ses.go             # AWS SES stub
│   └── templates/
│       ├── verification.html
│       ├── verification.txt
│       ├── otp.html
│       └── otp.txt
├── notification.go        # Main NotificationService
└── config.go              # Notification config types
```

### EmailSender Interface

```go
// internal/shared/notification/email/interface.go
package email

import "context"

type EmailSender interface {
    SendEmail(ctx context.Context, to, subject, htmlBody, textBody string) error
}
```

### Resend Implementation

```go
// internal/shared/notification/email/resend.go
package email

import (
    "context"
    "github.com/resend/resend-go/v2"
)

type ResendEmailSender struct {
    client    *resend.Client
    fromEmail string
    fromName  string
}

func NewResendEmailSender(apiKey, fromEmail, fromName string) *ResendEmailSender {
    return &ResendEmailSender{
        client:    resend.NewClient(apiKey),
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (r *ResendEmailSender) SendEmail(ctx context.Context, to, subject, htmlBody, textBody string) error {
    params := &resend.SendEmailRequest{
        From:    fmt.Sprintf("%s <%s>", r.fromName, r.fromEmail),
        To:      []string{to},
        Subject: subject,
        Html:    htmlBody,
        Text:    textBody,
    }
    _, err := r.client.Emails.SendWithContext(ctx, params)
    return err
}
```

### AWS SES Stub

```go
// internal/shared/notification/email/ses.go
package email

import (
    "context"
    "errors"
)

var ErrSESNotImplemented = errors.New("AWS SES email sender not implemented")

type SESEmailSender struct {
    region    string
    fromEmail string
}

func NewSESEmailSender(region, fromEmail string) *SESEmailSender {
    return &SESEmailSender{region: region, fromEmail: fromEmail}
}

func (s *SESEmailSender) SendEmail(ctx context.Context, to, subject, htmlBody, textBody string) error {
    return ErrSESNotImplemented
}
```

### NotificationService

```go
// internal/shared/notification/notification.go
package notification

import (
    "bytes"
    "context"
    "html/template"
    "path/filepath"

    "altalune/internal/shared/notification/email"
)

type NotificationService struct {
    emailSender email.EmailSender
    templates   *template.Template
    baseURL     string // Auth server URL for verification links
}

type VerificationEmailData struct {
    UserName         string
    VerificationLink string
    ExpiryHours      int
}

type OTPEmailData struct {
    UserName    string
    OTPCode     string
    ExpiryMinutes int
}

func NewNotificationService(sender email.EmailSender, templatesDir, baseURL string) (*NotificationService, error) {
    tmpl, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
    if err != nil {
        return nil, err
    }
    txtTmpl, err := template.ParseGlob(filepath.Join(templatesDir, "*.txt"))
    if err != nil {
        return nil, err
    }
    // Merge templates
    for _, t := range txtTmpl.Templates() {
        tmpl, _ = tmpl.AddParseTree(t.Name(), t.Tree)
    }

    return &NotificationService{
        emailSender: sender,
        templates:   tmpl,
        baseURL:     baseURL,
    }, nil
}

func (n *NotificationService) SendVerificationEmail(ctx context.Context, toEmail, token, userName string) error {
    data := VerificationEmailData{
        UserName:         userName,
        VerificationLink: fmt.Sprintf("%s/verify-email?token=%s", n.baseURL, token),
        ExpiryHours:      24,
    }

    htmlBody, textBody, err := n.renderTemplates("verification", data)
    if err != nil {
        return err
    }

    return n.emailSender.SendEmail(ctx, toEmail, "Verify your email address", htmlBody, textBody)
}

func (n *NotificationService) SendOTPEmail(ctx context.Context, toEmail, otp, userName string) error {
    data := OTPEmailData{
        UserName:      userName,
        OTPCode:       otp,
        ExpiryMinutes: 5,
    }

    htmlBody, textBody, err := n.renderTemplates("otp", data)
    if err != nil {
        return err
    }

    return n.emailSender.SendEmail(ctx, toEmail, "Your login code", htmlBody, textBody)
}

func (n *NotificationService) renderTemplates(name string, data interface{}) (string, string, error) {
    var htmlBuf, textBuf bytes.Buffer

    if err := n.templates.ExecuteTemplate(&htmlBuf, name+".html", data); err != nil {
        return "", "", err
    }
    if err := n.templates.ExecuteTemplate(&textBuf, name+".txt", data); err != nil {
        return "", "", err
    }

    return htmlBuf.String(), textBuf.String(), nil
}
```

### Email Templates

**verification.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button { display: inline-block; padding: 12px 24px; background: #2563eb; color: white; text-decoration: none; border-radius: 6px; }
        .footer { margin-top: 32px; color: #6b7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Verify your email address</h1>
        <p>Hi {{.UserName}},</p>
        <p>Please click the button below to verify your email address:</p>
        <p><a href="{{.VerificationLink}}" class="button">Verify Email</a></p>
        <p>Or copy and paste this link: {{.VerificationLink}}</p>
        <p>This link will expire in {{.ExpiryHours}} hours.</p>
        <div class="footer">
            <p>If you didn't create an account, you can safely ignore this email.</p>
            <p>— The Altalune Team</p>
        </div>
    </div>
</body>
</html>
```

**verification.txt:**
```
Verify your email address

Hi {{.UserName}},

Please click the link below to verify your email address:

{{.VerificationLink}}

This link will expire in {{.ExpiryHours}} hours.

If you didn't create an account, you can safely ignore this email.

— The Altalune Team
```

**otp.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .otp-code { font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #2563eb; padding: 16px; background: #f3f4f6; border-radius: 8px; text-align: center; }
        .footer { margin-top: 32px; color: #6b7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Your login code</h1>
        <p>Hi {{.UserName}},</p>
        <p>Use the following code to complete your login:</p>
        <div class="otp-code">{{.OTPCode}}</div>
        <p>This code will expire in {{.ExpiryMinutes}} minutes.</p>
        <div class="footer">
            <p>If you didn't request this code, you can safely ignore this email.</p>
            <p>— The Altalune Team</p>
        </div>
    </div>
</body>
</html>
```

**otp.txt:**
```
Your login code

Hi {{.UserName}},

Use the following code to complete your login:

{{.OTPCode}}

This code will expire in {{.ExpiryMinutes}} minutes.

If you didn't request this code, you can safely ignore this email.

— The Altalune Team
```

### Configuration Extension

Add to `internal/config/app.go`:

```go
type NotificationConfig struct {
    Email EmailNotificationConfig `yaml:"email" validate:"required"`
}

type EmailNotificationConfig struct {
    Provider string       `yaml:"provider" validate:"required,oneof=resend ses"`
    Resend   ResendConfig `yaml:"resend"`
    SES      SESConfig    `yaml:"ses"`
}

type ResendConfig struct {
    APIKey    string `yaml:"apiKey" validate:"required_if=Provider resend"`
    FromEmail string `yaml:"fromEmail" validate:"required_if=Provider resend,email"`
    FromName  string `yaml:"fromName" validate:"required_if=Provider resend"`
}

type SESConfig struct {
    Region    string `yaml:"region" validate:"required_if=Provider ses"`
    FromEmail string `yaml:"fromEmail" validate:"required_if=Provider ses,email"`
}
```

Add `Notification *NotificationConfig` to `AppConfig` struct.

### Config.yaml Example

```yaml
notification:
  email:
    provider: "resend"
    resend:
      apiKey: "re_xxxxxxxxxxxx"
      fromEmail: "noreply@altalune.com"
      fromName: "Altalune"
    ses:
      region: "ap-southeast-1"
      fromEmail: "noreply@altalune.com"
```

## Files to Create

- `internal/shared/notification/email/interface.go`
- `internal/shared/notification/email/resend.go`
- `internal/shared/notification/email/ses.go`
- `internal/shared/notification/email/templates/verification.html`
- `internal/shared/notification/email/templates/verification.txt`
- `internal/shared/notification/email/templates/otp.html`
- `internal/shared/notification/email/templates/otp.txt`
- `internal/shared/notification/notification.go`
- `internal/shared/notification/config.go`

## Files to Modify

- `internal/config/app.go` - Add NotificationConfig struct and field
- `internal/config/config.go` - Add getter methods for notification config
- `config.yaml` - Add notification section
- `config.example.yaml` - Add notification section with placeholders

## Testing Requirements

```go
// Manual test in main or test file
func TestSendVerificationEmail() {
    sender := email.NewResendEmailSender(apiKey, "noreply@test.com", "Test")
    svc, _ := notification.NewNotificationService(sender, "./templates", "http://localhost:3300")
    err := svc.SendVerificationEmail(context.Background(), "test@example.com", "abc123token", "John")
    // Check email received
}
```

## Commands to Run

```bash
# Add Resend SDK dependency
go get github.com/resend/resend-go/v2

# Build to verify compilation
make build
```

## Validation Checklist

- [ ] ResendEmailSender compiles and initializes
- [ ] SESEmailSender returns "not implemented" error
- [ ] Templates render with correct data substitution
- [ ] NotificationService loads templates successfully
- [ ] SendVerificationEmail generates correct email content
- [ ] SendOTPEmail generates correct email content
- [ ] Config validation works for notification settings
- [ ] Build succeeds without errors

## Definition of Done

- [ ] EmailSender interface defined
- [ ] Resend implementation complete and tested
- [ ] SES stub implemented with clear error message
- [ ] All 4 email templates created (HTML + text × 2)
- [ ] NotificationService created with both methods
- [ ] Configuration parsing implemented
- [ ] Templates embedded or loaded correctly
- [ ] Code follows established patterns

## Dependencies

- `github.com/resend/resend-go/v2` - Resend API client
- Go standard library `html/template` for templating

## Risk Factors

- **Low Risk**: Well-documented Resend API
- **Medium Risk**: Template path configuration needs to work in different environments

## Notes

- Resend free tier allows 100 emails/day, 3000/month
- Templates use Go's html/template for safety (auto-escaping)
- baseURL configuration allows different environments (dev/staging/prod)
- SES stub allows future implementation without code changes
- Consider using embed.FS for templates in production
