package oauth_auth

import (
	"encoding/base64"
	"testing"
)

func TestGenerateSecureToken(t *testing.T) {
	token, err := generateSecureToken(32)
	if err != nil {
		t.Fatalf("generateSecureToken failed: %v", err)
	}

	// Token should be base64url encoded (32 bytes = 43 chars without padding)
	expectedLen := 43 // base64url of 32 bytes without padding
	if len(token) != expectedLen {
		t.Errorf("expected token length %d, got %d", expectedLen, len(token))
	}

	// Should be valid base64url
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		t.Errorf("token is not valid base64url: %v", err)
	}

	if len(decoded) != 32 {
		t.Errorf("expected decoded length 32, got %d", len(decoded))
	}
}

func TestGenerateSecureTokenUniqueness(t *testing.T) {
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, err := generateSecureToken(32)
		if err != nil {
			t.Fatalf("generateSecureToken failed: %v", err)
		}
		if tokens[token] {
			t.Errorf("duplicate token generated: %s", token)
		}
		tokens[token] = true
	}
}

func TestGenerateSecureTokenDifferentLengths(t *testing.T) {
	testCases := []struct {
		length      int
		expectedLen int // base64url encoded length (without padding)
	}{
		{16, 22}, // 16 bytes -> 22 chars
		{32, 43}, // 32 bytes -> 43 chars
		{64, 86}, // 64 bytes -> 86 chars
	}

	for _, tc := range testCases {
		token, err := generateSecureToken(tc.length)
		if err != nil {
			t.Errorf("generateSecureToken(%d) failed: %v", tc.length, err)
			continue
		}
		if len(token) != tc.expectedLen {
			t.Errorf("generateSecureToken(%d): expected length %d, got %d", tc.length, tc.expectedLen, len(token))
		}
	}
}
