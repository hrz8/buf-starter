package email

import (
	"context"
	"errors"
	"testing"
)

func TestSESEmailSender_NotImplemented(t *testing.T) {
	sender := NewSESEmailSender("ap-southeast-1", "test@example.com")

	err := sender.SendEmail(context.Background(), "to@example.com", "Test", "<p>Test</p>", "Test")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, ErrSESNotImplemented) {
		t.Errorf("Expected ErrSESNotImplemented, got %v", err)
	}
}
