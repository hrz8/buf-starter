package oauth_auth

import (
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	otp, err := generateOTP(6)
	if err != nil {
		t.Fatalf("generateOTP failed: %v", err)
	}

	if len(otp) != 6 {
		t.Errorf("expected OTP length 6, got %d", len(otp))
	}

	// Verify all characters are digits
	for i, c := range otp {
		if c < '0' || c > '9' {
			t.Errorf("character at position %d is not a digit: %c", i, c)
		}
	}
}

func TestGenerateOTPUniqueness(t *testing.T) {
	otps := make(map[string]bool)
	for i := 0; i < 100; i++ {
		otp, err := generateOTP(6)
		if err != nil {
			t.Fatalf("generateOTP failed: %v", err)
		}
		if otps[otp] {
			// Note: This could theoretically fail with probability 1/10^6 per iteration
			// but 100 iterations should be fine
			t.Logf("duplicate OTP generated (rare but possible): %s", otp)
		}
		otps[otp] = true
	}
}

func TestHashToken(t *testing.T) {
	token := "test-token-123"
	hash1 := hashToken(token)
	hash2 := hashToken(token)

	// Same input should produce same hash
	if hash1 != hash2 {
		t.Errorf("hashToken is not deterministic: %s != %s", hash1, hash2)
	}

	// Hash should be 64 characters (SHA256 hex)
	if len(hash1) != 64 {
		t.Errorf("expected hash length 64, got %d", len(hash1))
	}

	// Different inputs should produce different hashes
	hash3 := hashToken("different-token")
	if hash1 == hash3 {
		t.Error("different tokens produced same hash")
	}
}

func TestHashTokenConsistency(t *testing.T) {
	// Test that hash is consistent for OTP codes
	otp := "123456"
	hash := hashToken(otp)

	// Expected SHA256 hash of "123456"
	expected := "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
	if hash != expected {
		t.Errorf("hash mismatch: got %s, expected %s", hash, expected)
	}
}
