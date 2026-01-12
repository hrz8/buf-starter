package oauthprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubClient struct {
	config *oauth2.Config
}

func NewGitHubClient(clientID, clientSecret, redirectURL string) *GitHubClient {
	return &GitHubClient{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (c *GitHubClient) GetAuthorizationURL(state string) string {
	return c.config.AuthCodeURL(state)
}

func (c *GitHubClient) ExchangeCodeForUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange: %w", err)
	}

	client := c.config.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("fetch user info: %w", err)
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("decode user info: %w", err)
	}

	email := githubUser.Email
	if email == "" {
		email, err = c.fetchPrimaryEmail(client)
		if err != nil {
			return nil, fmt.Errorf("fetch email: %w", err)
		}
	}

	firstName, lastName := parseName(githubUser.Name)

	return &UserInfo{
		ID:        fmt.Sprintf("%d", githubUser.ID),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		AvatarURL: githubUser.AvatarURL,
	}, nil
}

func (c *GitHubClient) fetchPrimaryEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	if len(emails) > 0 && emails[0].Verified {
		return emails[0].Email, nil
	}

	return "", fmt.Errorf("no verified email found")
}

func parseName(fullName string) (string, string) {
	if fullName == "" {
		return "", ""
	}

	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], strings.Join(parts[1:], " ")
}
