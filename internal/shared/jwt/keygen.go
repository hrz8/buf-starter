package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// GenerateRSAKeyPair generates an RSA key pair with the specified bit size.
// Recommended sizes: 2048 (faster) or 4096 (more secure).
//
// Example:
//
//	privateKey, err := GenerateRSAKeyPair(2048)
//	if err != nil {
//	    log.Fatal(err)
//	}
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	if bits != 2048 && bits != 4096 {
		return nil, fmt.Errorf("key size must be 2048 or 4096, got %d", bits)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("generate RSA key: %w", err)
	}

	// Validate the key
	if err := privateKey.Validate(); err != nil {
		return nil, fmt.Errorf("validate key: %w", err)
	}

	return privateKey, nil
}

// SavePrivateKeyPEM saves an RSA private key to a PEM file with secure permissions (0600).
// The file will be encoded in PKCS#1 format.
//
// Security: The file is created with 0600 permissions (read/write for owner only).
//
// Example:
//
//	err := SavePrivateKeyPEM(privateKey, "keys/jwt-private.pem")
//	if err != nil {
//	    log.Fatal(err)
//	}
func SavePrivateKeyPEM(key *rsa.PrivateKey, filepath string) error {
	if key == nil {
		return errors.New("private key cannot be nil")
	}

	// Encode key to PKCS#1 ASN.1 DER format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)

	// Create PEM block
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Create file with secure permissions (0600 = rw-------)
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create private key file: %w", err)
	}
	defer file.Close()

	// Write PEM to file
	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("encode private key: %w", err)
	}

	return nil
}

// SavePublicKeyPEM saves an RSA public key to a PEM file with 0644 permissions.
// The file will be encoded in PKIX format (standard for public keys).
//
// Example:
//
//	err := SavePublicKeyPEM(&privateKey.PublicKey, "keys/jwt-public.pem")
//	if err != nil {
//	    log.Fatal(err)
//	}
func SavePublicKeyPEM(key *rsa.PublicKey, filepath string) error {
	if key == nil {
		return errors.New("public key cannot be nil")
	}

	// Encode key to PKIX ASN.1 DER format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}

	// Create PEM block
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// Create file with public permissions (0644 = rw-r--r--)
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("create public key file: %w", err)
	}
	defer file.Close()

	// Write PEM to file
	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("encode public key: %w", err)
	}

	return nil
}

// LoadPrivateKeyPEM loads an RSA private key from a PEM file.
// The file must be in PKCS#1 format (RSA PRIVATE KEY).
//
// Example:
//
//	privateKey, err := LoadPrivateKeyPEM("keys/jwt-private.pem")
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadPrivateKeyPEM(filepath string) (*rsa.PrivateKey, error) {
	// Read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read private key file: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block from private key file")
	}

	// Verify block type
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected key type: %s (expected: RSA PRIVATE KEY)", block.Type)
	}

	// Parse PKCS#1 private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	// Validate key
	if err := privateKey.Validate(); err != nil {
		return nil, fmt.Errorf("validate private key: %w", err)
	}

	return privateKey, nil
}

// LoadPublicKeyPEM loads an RSA public key from a PEM file.
// The file must be in PKIX format (RSA PUBLIC KEY).
//
// Example:
//
//	publicKey, err := LoadPublicKeyPEM("keys/jwt-public.pem")
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadPublicKeyPEM(filepath string) (*rsa.PublicKey, error) {
	// Read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read public key file: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block from public key file")
	}

	// Verify block type
	if block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected key type: %s (expected: RSA PUBLIC KEY)", block.Type)
	}

	// Parse PKIX public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	// Type assert to RSA public key
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not an RSA public key")
	}

	return publicKey, nil
}
