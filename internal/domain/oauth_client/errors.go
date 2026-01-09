package oauth_client

import "errors"

var (
	ErrOAuthClientNotFound          = errors.New("oauth client not found")
	ErrOAuthClientAlreadyExists     = errors.New("oauth client already exists")
	ErrInvalidRedirectURI           = errors.New("invalid redirect URI")
	ErrClientSecretMismatch         = errors.New("client secret does not match")
	ErrDefaultClientCannotBeDeleted = errors.New("default dashboard client cannot be deleted")
	ErrClientBelongsToOtherProject  = errors.New("oauth client belongs to another project")
	ErrPKCECannotBeDisabled         = errors.New("PKCE cannot be disabled for default client")
)
