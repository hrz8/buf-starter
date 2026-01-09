package api_key

import (
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApiKeyQueryResult struct {
	ID         int64
	PublicID   string
	Name       string
	Expiration time.Time
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (r *ApiKeyQueryResult) ToApiKey() *ApiKey {
	return &ApiKey{
		ID:         r.PublicID,
		Name:       r.Name,
		Expiration: r.Expiration,
		Active:     r.Active,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

type ApiKey struct {
	ID         string
	Name       string
	Expiration time.Time
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *ApiKey) ToApiKeyProto() *altalunev1.ApiKey {
	return &altalunev1.ApiKey{
		Id:         m.ID,
		Name:       m.Name,
		Expiration: timestamppb.New(m.Expiration),
		Active:     m.Active,
		CreatedAt:  timestamppb.New(m.CreatedAt),
		UpdatedAt:  timestamppb.New(m.UpdatedAt),
	}
}

type CreateApiKeyInput struct {
	ProjectID  int64
	Name       string
	Expiration time.Time
}

type CreateApiKeyResult struct {
	ID         int64
	PublicID   string
	Name       string
	Key        string // Generated API key (only in result)
	Expiration time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UpdateApiKeyInput struct {
	ProjectID  int64
	PublicID   string // API Key's public ID
	Name       string
	Expiration time.Time
}

type UpdateApiKeyResult struct {
	ID         int64
	PublicID   string
	Name       string
	Expiration time.Time
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DeleteApiKeyInput struct {
	ProjectID int64
	PublicID  string // API Key's public ID
}

type ActivateApiKeyInput struct {
	ProjectID int64
	PublicID  string // API Key's public ID
}

type ActivateApiKeyResult struct {
	ID         int64
	PublicID   string
	Name       string
	Expiration time.Time
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DeactivateApiKeyInput struct {
	ProjectID int64
	PublicID  string // API Key's public ID
}

type DeactivateApiKeyResult struct {
	ID         int64
	PublicID   string
	Name       string
	Expiration time.Time
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
