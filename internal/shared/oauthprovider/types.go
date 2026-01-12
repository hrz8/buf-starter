package oauthprovider

import "context"

type UserInfo struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	AvatarURL string
}

type Client interface {
	ExchangeCodeForUserInfo(ctx context.Context, code string) (*UserInfo, error)
	GetAuthorizationURL(state string) string
}
