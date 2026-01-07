package oauth_provider

import "errors"

var (
	// ErrOAuthProviderNotFound is returned when an OAuth provider is not found
	ErrOAuthProviderNotFound = errors.New("oauth provider not found")

	// ErrDuplicateProviderType is returned when trying to create a provider with an existing type
	ErrDuplicateProviderType = errors.New("oauth provider with this provider type already exists")

	// ErrEncryptionFailed is returned when client secret encryption fails
	ErrEncryptionFailed = errors.New("failed to encrypt client secret")

	// ErrDecryptionFailed is returned when client secret decryption fails
	ErrDecryptionFailed = errors.New("failed to decrypt client secret")

	// ErrProviderTypeImmutable is returned when trying to change provider type
	ErrProviderTypeImmutable = errors.New("provider type cannot be changed after creation")
)
