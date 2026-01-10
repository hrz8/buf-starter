// Package pkce provides PKCE (Proof Key for Code Exchange) utilities per RFC 7636.
package pkce

import (
	"crypto/sha256"
	"encoding/base64"
)

// PKCE code challenge methods.
const (
	MethodS256  = "S256"
	MethodPlain = "plain"
)

// VerifyCodeChallenge verifies a code verifier against a code challenge.
func VerifyCodeChallenge(verifier, challenge, method string) bool {
	switch method {
	case MethodS256:
		hash := sha256.Sum256([]byte(verifier))
		computed := base64.RawURLEncoding.EncodeToString(hash[:])
		return computed == challenge
	case MethodPlain:
		return verifier == challenge
	default:
		return false
	}
}

// GenerateCodeChallenge generates a code challenge from a code verifier.
func GenerateCodeChallenge(verifier, method string) string {
	switch method {
	case MethodS256:
		hash := sha256.Sum256([]byte(verifier))
		return base64.RawURLEncoding.EncodeToString(hash[:])
	case MethodPlain:
		return verifier
	default:
		return ""
	}
}
