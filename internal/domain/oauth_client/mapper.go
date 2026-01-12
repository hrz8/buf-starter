package oauth_client

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToOAuthClientProto converts domain model to protobuf message
// OAuth clients are GLOBAL entities (not project-scoped)
func (c *OAuthClient) ToOAuthClientProto() *altalunev1.OAuthClient {
	return &altalunev1.OAuthClient{
		Id:              c.ID,
		Name:            c.Name,
		ClientId:        c.ClientID.String(),
		RedirectUris:    c.RedirectURIs,
		PkceRequired:    c.PKCERequired,
		IsDefault:       c.IsDefault,
		ClientSecretSet: true,       // Always true (never expose actual secret)
		AllowedScopes:   []string{}, // TODO: Implement scope assignment
		CreatedAt:       timestamppb.New(c.CreatedAt),
		UpdatedAt:       timestamppb.New(c.UpdatedAt),
	}
}
