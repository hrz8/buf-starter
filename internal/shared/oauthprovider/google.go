package oauthprovider

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleClient struct {
	config *oauth2.Config
}

func NewGoogleClient(clientID, clientSecret, redirectURL string) *GoogleClient {
	return &GoogleClient{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (c *GoogleClient) GetAuthorizationURL(state string) string {
	return c.config.AuthCodeURL(state)
}

func (c *GoogleClient) ExchangeCodeForUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange: %w", err)
	}

	client := c.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("fetch user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		GivenName string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture   string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("decode user info: %w", err)
	}

	return &UserInfo{
		ID:        googleUser.ID,
		Email:     googleUser.Email,
		FirstName: googleUser.GivenName,
		LastName:  googleUser.FamilyName,
		AvatarURL: googleUser.Picture,
	}, nil
}
