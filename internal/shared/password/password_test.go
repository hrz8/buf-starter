package password

import (
	"fmt"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	plainText := "mysecretpassword"

	hash, err := HashPassword(plainText, HashOption{
		Iterations: 2, // Production parameter
		Memory:     64 * 1024,
		Threads:    4,
		Len:        32,
	})
	if err != nil {
		t.Fatalf("HashPassword returned an unexpected error: %v", err)
	}

	// Verify it starts with $argon2id
	expectedPrefix := "$argon2id"
	if !strings.HasPrefix(hash, expectedPrefix) {
		t.Errorf("expected hash prefixed with %s", expectedPrefix)
	}

	// Verify it contains the correct parameters
	if !strings.Contains(hash, "m=65536") {
		t.Errorf("expected hash to contain m=65536")
	}
	if !strings.Contains(hash, "t=2") {
		t.Errorf("expected hash to contain t=2")
	}
	if !strings.Contains(hash, "p=4") {
		t.Errorf("expected hash to contain p=4")
	}

	// Verify hash is not empty and has reasonable length
	if len(hash) < 50 {
		t.Errorf("expected hash length > 50, but got %d", len(hash))
	}
}

func TestVerifyPassword(t *testing.T) {
	// Pre-generated hash for "mysecretpassword" with production parameters
	// (Generated with: iterations=2, memory=64*1024, threads=4, len=32)
	plainText := "mysecretpassword"

	// Generate a fresh hash for testing
	hash, err := HashPassword(plainText, HashOption{
		Iterations: 2,
		Memory:     64 * 1024,
		Threads:    4,
		Len:        32,
	})
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	t.Run("Match = true", func(t *testing.T) {
		match, err := VerifyPassword(plainText, hash)
		if err != nil {
			t.Fatalf("VerifyPassword returned an unexpected error: %v", err)
		}

		if !match {
			t.Error("expected plain password and encoded hash to match")
		}
	})

	t.Run("Match = false", func(t *testing.T) {
		wrongPlainText := "wrongpassword"

		match, err := VerifyPassword(wrongPlainText, hash)
		if err != nil {
			t.Fatalf("VerifyPassword returned an unexpected error: %v", err)
		}

		if match {
			t.Error("expected wrong password and hash to not match, but they did")
		}
	})

	t.Run("Invalid hash format", func(t *testing.T) {
		invalidHash := "invalid-hash-format"

		_, err := VerifyPassword(plainText, invalidHash)
		if err != ErrInvalidHashedString {
			t.Errorf("expected ErrInvalidHashedString, got %v", err)
		}
	})

	t.Run("Incompatible version", func(t *testing.T) {
		// Hash with wrong version number
		invalidVersionHash := "$argon2id$v=18$m=65536,t=2,p=4$CO5hu/iRl5ey1rr8h4FbRQ$qgk/PEQzuAdh4b06CmxTS/djb7F7Fojdhubl0QEKWQw"

		_, err := VerifyPassword(plainText, invalidVersionHash)
		if err != ErrIncompatibleVersion {
			t.Errorf("expected ErrIncompatibleVersion, got %v", err)
		}
	})
}

func TestHashPasswordWithDifferentParameters(t *testing.T) {
	plainText := "testsecret123"

	tests := []struct {
		name       string
		iterations uint32
		memory     uint32
		threads    uint8
		length     uint32
	}{
		{"Low security", 1, 32 * 1024, 2, 16},
		{"Production", 2, 64 * 1024, 4, 32},
		{"High security", 3, 128 * 1024, 8, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(plainText, HashOption{
				Iterations: tt.iterations,
				Memory:     tt.memory,
				Threads:    tt.threads,
				Len:        tt.length,
			})
			if err != nil {
				t.Fatalf("HashPassword failed: %v", err)
			}

			// Verify password
			match, err := VerifyPassword(plainText, hash)
			if err != nil {
				t.Fatalf("VerifyPassword failed: %v", err)
			}
			if !match {
				t.Error("password should match")
			}

			// Verify hash contains correct parameters
			expectedParams := []string{
				"$argon2id$",
				"v=19",
				fmt.Sprintf("m=%d", tt.memory),
				fmt.Sprintf("t=%d", tt.iterations),
				fmt.Sprintf("p=%d", tt.threads),
			}
			for _, param := range expectedParams {
				if !strings.Contains(hash, param) {
					t.Errorf("hash missing expected parameter: %s", param)
				}
			}
		})
	}
}

func TestHashPasswordProducesUniqueHashes(t *testing.T) {
	plainText := "samepassword"

	hash1, err := HashPassword(plainText, HashOption{
		Iterations: 2,
		Memory:     64 * 1024,
		Threads:    4,
		Len:        32,
	})
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	hash2, err := HashPassword(plainText, HashOption{
		Iterations: 2,
		Memory:     64 * 1024,
		Threads:    4,
		Len:        32,
	})
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Hashes should be different due to random salt
	if hash1 == hash2 {
		t.Error("expected hashes to be different (random salt), but they were the same")
	}

	// But both should verify correctly
	match1, _ := VerifyPassword(plainText, hash1)
	match2, _ := VerifyPassword(plainText, hash2)

	if !match1 || !match2 {
		t.Error("both hashes should verify correctly")
	}
}
