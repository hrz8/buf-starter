package password

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// generateRandomBytes generates n random bytes using crypto/rand
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// HashPassword hashes a plaintext password using Argon2id with the given options
// Returns a PHC string format: $argon2id$v=19$m=65536,t=2,p=4$<base64-salt>$<base64-hash>
func HashPassword(plainText string, opt HashOption) (string, error) {
	// Generate cryptographically secure salt
	salt, err := generateRandomBytes(saltLen)
	if err != nil {
		return "", err
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(plainText), salt, opt.Iterations, opt.Memory, opt.Threads, opt.Len)

	// Encode salt and hash to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return PHC string format
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		opt.Memory,
		opt.Iterations,
		opt.Threads,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}
