package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenClaims represents the claims in an OAuth access token.
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Scope string   `json:"scope,omitempty"`
	Email string   `json:"email,omitempty"`
	Name  string   `json:"name,omitempty"`
	Perms []string `json:"perms"`
}
