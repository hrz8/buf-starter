package api_key

import "errors"

var (
	ErrApiKeyNotFound      = errors.New("api key not found")
	ErrApiKeyAlreadyExists = errors.New("api key with this name already exists")
	ErrApiKeyExpired       = errors.New("api key has expired")
)
