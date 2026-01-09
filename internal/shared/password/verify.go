package password

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// decodeHash decodes a PHC format hash string and returns the hash options, salt, and hash
func decodeHash(encodedHash string) (*HashOption, []byte, []byte, error) {
	encodedSplit := strings.Split(encodedHash, "$")
	if len(encodedSplit) != splitN {
		return nil, nil, nil, ErrInvalidHashedString
	}

	// Parse version
	var version int
	_, err := fmt.Sscanf(encodedSplit[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Parse parameters (m, t, p)
	o := &HashOption{}
	_, err = fmt.Sscanf(encodedSplit[3], "m=%d,t=%d,p=%d", &o.Memory, &o.Iterations, &o.Threads)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.Strict().DecodeString(encodedSplit[4])
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode hash
	hash, err := base64.RawStdEncoding.Strict().DecodeString(encodedSplit[5])
	if err != nil {
		return nil, nil, nil, err
	}
	o.Len = uint32(len(hash))

	return o, salt, hash, nil
}

// VerifyPassword verifies a plaintext password against an encoded hash
// Uses constant-time comparison to prevent timing attacks
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Decode the hash to get parameters, salt, and hash
	o, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Re-hash the password with the same parameters and salt
	otherHash := argon2.IDKey([]byte(password), salt, o.Iterations, o.Memory, o.Threads, o.Len)

	// Constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
