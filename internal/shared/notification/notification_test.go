package notification

import (
	"context"
	"strings"
	"testing"
)

// mockEmailSender is a mock implementation for testing.
type mockEmailSender struct {
	lastTo      string
	lastSubject string
	lastHTML    string
	lastText    string
}

func (m *mockEmailSender) SendEmail(_ context.Context, to, subject, htmlBody, textBody string) error {
	m.lastTo = to
	m.lastSubject = subject
	m.lastHTML = htmlBody
	m.lastText = textBody
	return nil
}

func TestNewNotificationService(t *testing.T) {
	sender := &mockEmailSender{}
	svc, err := NewNotificationService(sender, "http://localhost:3300")
	if err != nil {
		t.Fatalf("Failed to create notification service: %v", err)
	}
	if svc == nil {
		t.Fatal("Service should not be nil")
	}
}

func TestSendVerificationEmail(t *testing.T) {
	sender := &mockEmailSender{}
	svc, err := NewNotificationService(sender, "http://localhost:3300")
	if err != nil {
		t.Fatalf("Failed to create notification service: %v", err)
	}

	err = svc.SendVerificationEmail(context.Background(), "test@example.com", "abc123token", "John Doe")
	if err != nil {
		t.Fatalf("Failed to send verification email: %v", err)
	}

	if sender.lastTo != "test@example.com" {
		t.Errorf("Expected to=test@example.com, got %s", sender.lastTo)
	}
	if sender.lastSubject != "Verify your email address" {
		t.Errorf("Expected subject='Verify your email address', got %s", sender.lastSubject)
	}
	if !strings.Contains(sender.lastHTML, "John Doe") {
		t.Error("HTML body should contain user name")
	}
	if !strings.Contains(sender.lastHTML, "http://localhost:3300/verify-email?token=abc123token") {
		t.Error("HTML body should contain verification link")
	}
	if !strings.Contains(sender.lastText, "John Doe") {
		t.Error("Text body should contain user name")
	}
	if !strings.Contains(sender.lastText, "http://localhost:3300/verify-email?token=abc123token") {
		t.Error("Text body should contain verification link")
	}
}

func TestSendOTPEmail(t *testing.T) {
	sender := &mockEmailSender{}
	svc, err := NewNotificationService(sender, "http://localhost:3300")
	if err != nil {
		t.Fatalf("Failed to create notification service: %v", err)
	}

	err = svc.SendOTPEmail(context.Background(), "test@example.com", "123456", "Jane Doe")
	if err != nil {
		t.Fatalf("Failed to send OTP email: %v", err)
	}

	if sender.lastTo != "test@example.com" {
		t.Errorf("Expected to=test@example.com, got %s", sender.lastTo)
	}
	if sender.lastSubject != "Your login code" {
		t.Errorf("Expected subject='Your login code', got %s", sender.lastSubject)
	}
	if !strings.Contains(sender.lastHTML, "Jane Doe") {
		t.Error("HTML body should contain user name")
	}
	if !strings.Contains(sender.lastHTML, "123456") {
		t.Error("HTML body should contain OTP code")
	}
	if !strings.Contains(sender.lastText, "Jane Doe") {
		t.Error("Text body should contain user name")
	}
	if !strings.Contains(sender.lastText, "123456") {
		t.Error("Text body should contain OTP code")
	}
}
