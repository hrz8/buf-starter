package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

// JWK represents a JSON Web Key for RSA.
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS represents a JSON Web Key Set.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// GenerateJWKS creates a JWKS from an RSA public key.
func GenerateJWKS(publicKey *rsa.PublicKey, kid string) *JWKS {
	return &JWKS{
		Keys: []JWK{
			{
				Kty: "RSA",
				Use: "sig",
				Kid: kid,
				Alg: "RS256",
				N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
				E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
			},
		},
	}
}

// GenerateJWKS creates a JWKS from the signer's public key.
func (s *Signer) GenerateJWKS() *JWKS {
	return GenerateJWKS(s.publicKey, s.kid)
}
